package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"unicode"
)

func main() {

	in := bufio.NewReader(os.Stdin)

	var cnt = make(map[rune]int)

	var invalid = 0

	for { // for省略了全部参数
		//
		r, n, err := in.ReadRune() // 分别表示returns rune, nbytes, error
		
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "error : %v", err)
			os.Exit(1)
		}
		if r == unicode.ReplacementChar { // 判断r是否为\uFFFD，无效码点都会被此字符替代
			invalid++
			continue
		}

		cnt[r]++
		_ = n
	}

	for k, v := range cnt {
		fmt.Printf("%q\t%v\n", k, v) //%q有%c和%s的作用，但会给字符加上单引号或双引号
	}
}
