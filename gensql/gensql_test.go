package gensql

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestGenSQLInt(t *testing.T) {
	sql := "select * from test where id = ?, id32 = ?, id64 = ?"
	var args1 []interface{}
	args1 = append(args1, 16, int32(32), int64(64))
	ret, err := GenSQL(sql, args1, true, time.Local)
	if err != nil {
		fmt.Printf("err: %s\n", err.Error())
	}
	fmt.Printf("sql:%s\n", ret)
	if !strings.EqualFold(ret, "select * from test where id = 16, id32 = 32, id64 = 64") {
		t.Error("not match")
	}
}

func TestGenSQLFloat(t *testing.T) {
	sql := "select * from test where id = ?"
	var args1 []interface{}
	args1 = append(args1, float64(666))
	ret, err := GenSQL(sql, args1, true, time.Local)
	if err != nil {
		fmt.Printf("err: %s\n", err.Error())
	}
	fmt.Printf("sql:%s\n", ret)
	if !strings.EqualFold(ret, "select * from test where id = 666") {
		t.Error("not match")
	}
}

func TestGenSQLBool(t *testing.T) {
	sql := "select * from test where id = ?"
	var args1 []interface{}
	args1 = append(args1, true)
	ret, err := GenSQL(sql, args1, true, time.Local)
	if err != nil {
		fmt.Printf("err: %s\n", err.Error())
	}
	fmt.Printf("sql:%s\n", ret)
	if !strings.EqualFold(ret, "select * from test where id = 1") {
		t.Error("not match")
	}
}

func TestGenSQLTime(t *testing.T) {
	sql := "select * from test where id = ?"
	var args1 []interface{}
	args1 = append(args1, time.Now())
	ret, err := GenSQL(sql, args1, true, time.Local)
	if err != nil {
		fmt.Printf("err: %s\n", err.Error())
	}
	fmt.Printf("sql:%s\n", ret)
	if !strings.EqualFold(ret, "select * from test where id = '2020-05-06 20:07:21.492144'") {
		t.Error("not match")
	}
}

func TestGenSQLByte(t *testing.T) {
	sql := "select * from test where id = ?"
	var args1 []interface{}
	args1 = append(args1, []byte("hel'abc'lo"))
	ret, err := GenSQL(sql, args1, false, time.Local)
	if err != nil {
		fmt.Printf("err: %s\n", err.Error())
	}
	fmt.Printf("sql:%s\n", ret)
	if !strings.EqualFold(ret, "select * from test where id = '2020-05-06 20:07:21.492144'") {
		t.Error("not match")
	}
}

func TestGenSQLString(t *testing.T) {
	sql := "select * from test where id = ?"
	var args1 []interface{}
	args1 = append(args1, string("hel'abc'lo"))
	ret, err := GenSQL(sql, args1, true, time.Local)
	if err != nil {
		fmt.Printf("err: %s\n", err.Error())
	}
	fmt.Printf("sql:%s\n", ret)
	if !strings.EqualFold(ret, "select * from test where id = '2020-05-06 20:07:21.492144'") {
		t.Error("not match")
	}
}