package main

import (
	"flag"
	"fmt"
	"strconv"
	"unicode"
)

type Num int64

func (n *Num) String() string {
	return strconv.FormatInt(int64(*n), 10) // FormatInt接受的是int64
}

// 在不影响Num类型的前提下，使用flag包增添一种Num类型符号

// 定义一个类型，它实现flag.Value接口，并且可与Num类型相互转化

type numflag struct {
	Num
}
// Num已经实现了String()，后续找接口也会bfs地找到Num的String()，就不用再另定义了

// 就如flag包其他函数一样，接受标志参数名,默认值,描述信息，返回指针
func NumFlag(name string, num Num, usage string) *Num {
	var nf = numflag{num}
	// 第一个参数是flag.Value接口类型，要求有String() string和Set(string) error
	flag.CommandLine.Var(&nf, name, usage) // 注册新的flag，会根据-name把值传给&nf.Set()。配合flag.Phase()使用
	return &(nf.Num)
}

func (nf *numflag) Set(s string) error { // 解析字符串
	var val int64
	if len(s) > 1 && s[0] == '0' {
		if s[1] == 'o' || unicode.IsDigit(rune(s[1])) {
			if s[1] == 'o' {
				s = s[1:]
			}
			fmt.Sscanf(s[1:], "%o", &val)
			nf.Num = Num(val)
			return nil
		}
		switch s[1] {
		case 'x':
			fmt.Sscanf(s[2:], "%x", &val)
		case 'b':
			fmt.Sscanf(s[2:], "%b", &val)
		default:
			return fmt.Errorf("invalid Num value")
		}
		nf.Num = Num(val)
		return nil
	} else if unicode.IsDigit(rune(s[0])) {
		fmt.Sscanf(s, "%d", &val)
		nf.Num = Num(val)
		return nil
	}
	return fmt.Errorf("invalid Num value")
}
