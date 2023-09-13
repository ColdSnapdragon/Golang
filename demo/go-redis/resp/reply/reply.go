package reply

// 实现RESP的五种通信格式

import (
	"go-redis/interface/resp"
)

var (
	nullBulkReplyBytes = []byte("$-1")

	// CRLF 常用的结尾
	CRLF = "\r\n"
)

/* ---- Bulk Reply ---- */
// 多行字符串，以$开头

// BulkReply stores a binary-safe string
type BulkReply struct {
	Arg []byte
}

// MakeBulkReply creates  BulkReply
func MakeBulkReply(arg []byte) *BulkReply {
	return nil // todo
}

// ToBytes marshal redis.Reply
func (r *BulkReply) ToBytes() []byte {
	return nil // todo
}

/* ---- Multi Bulk Reply ---- */
// 数组，以*开头

// MultiBulkReply stores a list of string
type MultiBulkReply struct {
	Args [][]byte
}

// MakeMultiBulkReply creates MultiBulkReply
func MakeMultiBulkReply(args [][]byte) *MultiBulkReply {
	return nil // todo
}

// ToBytes marshal redis.Reply
func (r *MultiBulkReply) ToBytes() []byte {
	return nil // todo
}

/* ---- Status Reply ---- */

// StatusReply stores a simple status string
type StatusReply struct {
	Status string
}

// MakeStatusReply creates StatusReply
func MakeStatusReply(status string) *StatusReply {
	return nil // todo
}

// ToBytes marshal redis.Reply
func (r *StatusReply) ToBytes() []byte {
	return nil // todo
}

/* ---- Int Reply ---- */

// IntReply stores an int64 number
type IntReply struct {
	Code int64
}

// MakeIntReply creates int reply
func MakeIntReply(code int64) *IntReply {
	return nil // todo
}

// ToBytes marshal redis.Reply
func (r *IntReply) ToBytes() []byte {
	return nil // todo
}

/* ---- Error Reply ---- */

// ErrorReply is an error and redis.Reply
type ErrorReply interface {
	Error() string
	ToBytes() []byte
}

// StandardErrReply represents handler error
type StandardErrReply struct {
	Status string
}

// ToBytes marshal redis.Reply
func (r *StandardErrReply) ToBytes() []byte {
	return []byte("-" + r.Status + CRLF)
}

func (r *StandardErrReply) Error() string {
	return r.Status
}

// MakeErrReply creates StandardErrReply
func MakeErrReply(status string) *StandardErrReply {
	return &StandardErrReply{
		Status: status,
	}
}

// IsErrorReply returns true if the given reply is error
func IsErrorReply(reply resp.Reply) bool {
	return reply.ToBytes()[0] == '-'
}
