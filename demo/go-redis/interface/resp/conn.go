package resp

// Connection 该接口代表一个redis的连接。写成接口是为了未来Connection可能有不同的实现
type Connection interface {
	Write([]byte) error
	GetDBIndex() int
	SelectDB(int)
}
