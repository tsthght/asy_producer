package message
type Message struct {
	Db     string `json:"database-name"`
	Tb     string `json:"table-name"`
	Sql    string `json:"statement"`
	Cts    int64  `json:"committed-timestamp"`
	Ats    int64  `json:"applied-timestamp"`
	seq    int64  `json:"uuid"`
}
