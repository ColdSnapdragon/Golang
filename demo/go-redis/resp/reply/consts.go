package reply

// 实现一些固定的回复

// PongReply 回复PONG
type PongReply struct{}

var pongReply = []byte("+PONG\r\n")

func (r PongReply) ToBytes() []byte {
	return pongReply
}

// 很多包喜欢暴露一个make方法，以便获取一个保内对象

func MakePongReply() *PongReply {
	return new(PongReply)
}

// OkReply 回复OK
type OkReply struct{}

var okReply = []byte("+OK\r\n")

func (o OkReply) ToBytes() []byte {
	return okReply
}

// 很多包喜欢暴露一个make方法，以便获取一个保内对象

func MakeOkReply() *OkReply {
	return new(OkReply)
}

// NullBulkReply is empty string
type NullBulkReply struct{}

var nullBulkBytes = []byte("$-1\r\n") // 表示“空回复”

func (o NullBulkReply) ToBytes() []byte {
	return nullBulkBytes
}

func MakeNullBulkReply() *NullBulkReply {
	return new(NullBulkReply)
}

// EmptyMultiBulkReply is a empty list
type EmptyMultiBulkReply struct{}

var emptyMultiBulkBytes = []byte("*0\r\n")

// ToBytes marshal redis.Reply
func (r *EmptyMultiBulkReply) ToBytes() []byte {
	return emptyMultiBulkBytes
}

// NoReply respond nothing, for commands like subscribe
type NoReply struct{}

var noBytes = []byte("")

// ToBytes marshal redis.Reply
func (r *NoReply) ToBytes() []byte {
	return noBytes
}
