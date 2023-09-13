package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", handler) // 所有发送到/路径下(/即当前站点，这里为localhost:8000)的请求和handler函数关联起来
	log.Fatal(http.ListenAndServe("localhost:8080", nil)) // 启动服务
	// ListenAndServe always returns a non-nil error.
}

func handler(w http.ResponseWriter, r *http.Request) {
	// 写入回应
	fmt.Fprintf(w, "URl.Path = %s", r.URL.Path) 
}