package main

import (
	"flag"
	"fmt"
	"strings"
)

// 标志参数的名字,默认值,描述信息
// 返回的是传入值的指针(不会拷贝新值)
var sep = flag.String("s", " ", "分隔符")
var bl = flag.Bool("n", false, "是否打印双引号")

//一旦使用-n，*bl的值就是true
//使用-h或-help查看描述信息

func main() {
	flag.Parse() //更新传给标志参数的值(如果有的话)(之前是默认值)
	//由于存在标志参数，不能简单地用os.Args[1:]了，flag.Args()返回其中非标志参数的部分
	var s = strings.Join(flag.Args(), *sep) //使用join函数拼接一个字符串切片
	if *bl {
		fmt.Printf("%q", s)
	} else {
		fmt.Printf("%s", s)
	}
}
