package info

import (
	"fmt"
	"sync"
	"strings"

	gosql "database/sql"
	"github.com/pingcap/errors"
)

const (
	colsSQL = `
SELECT column_name, extra FROM information_schema.columns
WHERE table_schema = ? AND table_name = ?;`
	uniqKeysSQL = `
SELECT non_unique, index_name, seq_in_index, column_name 
FROM information_schema.statistics
WHERE table_schema = ? AND table_name = ?
ORDER BY seq_in_index ASC;`
)

type tableInfo struct {
	columns    []string
	primaryKey *indexInfo
	// include primary key if have
	uniqueKeys []indexInfo
}

type indexInfo struct {
	name    string
	columns []string
}

type TableInfo struct {
	tableInfos sync.Map
	db *gosql.DB
}

func (ti *TableInfo) GetTableInfo(schema, table string) (info *tableInfo, err error){
	v, ok := ti.tableInfos.Load(quoteSchema(schema, table))
	if ok {
		info = v.(*tableInfo)
		return
	}
	return ti.RefreshTableInfo(schema, table)
}

func (ti *TableInfo) RefreshTableInfo(schema string, table string) (info *tableInfo, err error) {

	if len(schema) == 0 {
		return nil, errors.New("schema is empty")
	}

	if len(table) == 0 {
		return nil, nil
	}

	info, err = ti.UtilGetTableInfo(ti.db, schema, table)
	if err != nil {
		return info, errors.Trace(err)
	}

	if len(info.uniqueKeys) == 0 {
	}

	ti.tableInfos.Store(quoteSchema(schema, table), info)

	return
}

func (ti *TableInfo) UtilGetTableInfo(db *gosql.DB, schema string, table string) (info *tableInfo, err error) {
	info = new(tableInfo)

	if info.columns, err = getColsOfTbl(db, schema, table); err != nil {
		return nil, errors.Annotatef(err, "table `%s`.`%s`", schema, table)
	}

	if info.uniqueKeys, err = getUniqKeys(db, schema, table); err != nil {
		return nil, errors.Trace(err)
	}

	// put primary key at first place
	// and set primaryKey
	for i := 0; i < len(info.uniqueKeys); i++ {
		if info.uniqueKeys[i].name == "PRIMARY" {
			info.uniqueKeys[i], info.uniqueKeys[0] = info.uniqueKeys[0], info.uniqueKeys[i]
			info.primaryKey = &info.uniqueKeys[0]
			break
		}
	}

	return
}

func getColsOfTbl(db *gosql.DB, schema, table string) ([]string, error) {
	rows, err := db.Query(colsSQL, schema, table)
	if err != nil {
		return nil, errors.Trace(err)
	}
	defer rows.Close()

	cols := make([]string, 0, 1)
	for rows.Next() {
		var name, extra string
		err = rows.Scan(&name, &extra)
		if err != nil {
			return nil, errors.Trace(err)
		}
		isGenerated := strings.Contains(extra, "VIRTUAL GENERATED") || strings.Contains(extra, "STORED GENERATED")
		if isGenerated {
			continue
		}
		cols = append(cols, name)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Trace(err)
	}

	// if no any columns returns, means the table not exist.
	if len(cols) == 0 {
		return nil, errors.New("table not exist")
	}

	return cols, nil
}

func getUniqKeys(db *gosql.DB, schema, table string) (uniqueKeys []indexInfo, err error) {
	rows, err := db.Query(uniqKeysSQL, schema, table)
	if err != nil {
		err = errors.Trace(err)
		return
	}
	defer rows.Close()

	var nonUnique int
	var keyName string
	var columnName string
	var seqInIndex int // start at 1

	// get pk and uk
	// key for PRIMARY or other index name
	for rows.Next() {
		err = rows.Scan(&nonUnique, &keyName, &seqInIndex, &columnName)
		if err != nil {
			err = errors.Trace(err)
			return
		}

		if nonUnique == 1 {
			continue
		}

		var i int
		// Search for indexInfo with the current keyName
		for i = 0; i < len(uniqueKeys); i++ {
			if uniqueKeys[i].name == keyName {
				uniqueKeys[i].columns = append(uniqueKeys[i].columns, columnName)
				break
			}
		}
		// If we don't find the indexInfo with the loop above, create a new one
		if i == len(uniqueKeys) {
			uniqueKeys = append(uniqueKeys, indexInfo{keyName, []string{columnName}})
		}
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Trace(err)
	}

	return
}

func CreateDB(user string, password string, host string, port int) (db *gosql.DB, err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=utf8mb4,utf8&interpolateParams=true&readTimeout=1m&multiStatements=true", user, password, host, port)
	db, err = gosql.Open("mysql", dsn)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return db, nil
}

func quoteSchema(schema string, table string) string {
	return fmt.Sprintf("`%s`.`%s`", escapeName(schema), escapeName(table))
}

func escapeName(name string) string {
	return strings.Replace(name, "`", "``", -1)
}