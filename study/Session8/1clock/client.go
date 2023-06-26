package main

import (
	"io"
	"log"
	"net"
	"os"
)

func keepCopy(dst io.Writer, src io.Reader) {
	if _,err := io.Copy(dst, src); err!=nil { // io.Copy方法只有当src到达EOF或者发生其他错误时才退出。如果因为EOF结束，则err==nil而非nil==EOF
		log.Fatal(err)
	} 
}

func main() {
	conn, err := net.Dial("tcp",":9090") // 连接服务端
	if err!=nil {
		log.Fatal(err)
	}
	defer conn.Close()
	keepCopy(os.Stdout,conn)
}