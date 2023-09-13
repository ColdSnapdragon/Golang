package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	for _, url := range os.Args[1:] {
		resp, err := http.Get(url) //Get函数创建HTTP请求，会在resp这个结构体中得到访问的请求结果
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch : %v\n", err)
			os.Exit(1)
		}
		//resp的Body字段包括一个可读的服务器响应流
		// cont,err := ioutil.ReadAll(resp.Body) 这是原先的。但是io/ioutil包自go 1.16后被逐渐弃用，相同功能被挪到os或io包
		cont, err := io.ReadAll(resp.Body) //从response中读取到全部内容
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch : reading : %s : %v\n", url, err)
			os.Exit(1)
		}
		fmt.Printf("%s\n", cont)
	}
}
