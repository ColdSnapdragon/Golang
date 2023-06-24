package main

import (
	"io"
	"log"
	"net"
	"time"
)

func handleConn(c net.Conn) {
	defer c.Close()
	for {
		// 这里使用io.WriteString方法，可以检测写入失败
		_, err := io.WriteString(c, time.Now().Format("03:04:05PM 2006\n")) // Format方法的参数是一个格式化模板，标识如何来格式化时间。
		// 模板格式限定在Mon Jan 2 03:04:05PM 2006 UTC-0700

		if err!=nil {
			return
		}
		time.Sleep(time.Second)
	}
}

func main() {
	listener, err := net.Listen("tcp",":9090"); // 返回一个监听器
	if err != nil {
		log.Fatal(err) // 打印错误并执行os.Exit(1)
	}
	for {
		conn, err := listener.Accept() // 进行箭头，会阻塞直到一个tcp连接建立
		if err != nil {
			log.Fatal(err)
		}
		go handleConn(conn)
	}
}