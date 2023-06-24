package main

import (
	"fmt"
	"net/http"
)

type database map[string]float32

func (db database) list(w http.ResponseWriter, req *http.Request){
	for item,price:= range db {
		fmt.Fprintf(w,"%s: %f\n",item,price)
	}
}

func (db database) price(w http.ResponseWriter, req *http.Request) {
    item := req.URL.Query().Get("name")
    price, ok := db[item]
    if !ok {// 常用状态码在http包中均有常量值，这里是404
        w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Item not found !\n")
        return
    }
    fmt.Fprintf(w, "%f\n", price)
}

func main() {
	db:=database{
		"apple":5.5,
		"melon":10.20,
	}
	mux:=http.NewServeMux(); // 获得一个路由(ServerMux类型)
	mux.Handle("/list", http.HandlerFunc(db.list)) // 将函数值强转为正确的接口类型
	mux.HandleFunc("/price", db.price) // ServerMux提供简便的HandleFunc方法，不用手动强转
	http.ListenAndServe(":8080",mux) // 表示localhost:8080。路由本身也是合法的handler接口
}

/*
http://127.0.0.1:8080/list
http://127.0.0.1:8080/price?name=apple
http://127.0.0.1:8080/price?name=abc
*/