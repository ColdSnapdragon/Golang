package parse

import (
	"bufio"
	"go-redis/interface/resp"
	"io"
)

// 解析客户端发送的内容

type Payload struct {
	Data resp.Reply // 客户端发给服务器的格式也是一样的，这样也用Reply
	Err  error
}

// 维护解析过程的状态
type readState struct {
	readingMultiLine  bool     // 解析的是单行还是多行
	expectedArgsCount int      // 正在解析的指令有几条数据
	msgType           byte     // 数据类型
	args              [][]byte // 真正数据
	bulklen           int64    // 数组的串长度(若数据是数组的话)
}

func (s *readState) finished() bool {
	return true // todo
}

func ParseStream(reader io.Reader) <-chan *Payload {
	return nil // todo
}

func parse0(reader io.Reader, ch chan<- *Payload) {
	//todo
}

func readLine(bufReader *bufio.Reader, state *readState) ([]byte, bool, error) {
	var msg []byte
	var err error
	_ = err // todo
	// 正常串

	// 批量串)

	return msg, false, nil
}
