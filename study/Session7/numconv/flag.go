package main

import (
	"flag"
	"fmt"
)

var num = NumFlag("n", 0, "一个2/8/10/16进制数字")

func main() {
	flag.Parse()
	fmt.Println(*num)
}
