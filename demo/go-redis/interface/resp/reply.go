package resp

// Reply 该接口表示各种服务端对客户端的回复
type Reply interface {
	ToBytes() []byte
}

type ErrorReply interface {
	Error() string
	ToBytes() []byte
}
