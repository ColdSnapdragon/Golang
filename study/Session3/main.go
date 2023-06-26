package main

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

func main() {
	fmt.Printf("%d %[1]d %[1]x %#[1]x\n", 10) // 10 10 a 0xa
	// [1]指定要用第1个操作数。#则会让八进制加上前缀0、十六进制加上前缀0x

	newline := '\n'
	unicode := '我'
	// %q打印单引号括起来的字符
	fmt.Printf("%d %[1]c %[1]q\n", unicode) // 25105 我 '我'
	fmt.Printf("%d %[1]q\n", newline)       // 10(ascii) '\n'

	var v1 uint8 = 255  // int8 int16 int32 int64及其对应的uint
	fmt.Println(v1 + 1) // 0

	v1 = v1 &^ 1 // 取反1，然后和v1相与(即，将v1的某些位置为0)

	var c1 complex128 = complex(1, 2) //complex128表明实部虚部由两个float64组成
	c2 := 1 + 2i                      // 默认complex128
	fmt.Println(c1 == c2)             //true

	//var v2 int = (1 != 2) // 错误。布尔值并不会隐式转换为数字值0或1
	//var v2 int = int(1 != 2) // 这样也不行。
	var v2 int = b2i(1 != 2) // 包装一个函数

	var s1 string = "hello"
	//s1[0] = 'a' // 错误。字符串不可修改
	s1 += "w"                    // s1成为一个新的字符串值
	fmt.Println(len(s1), s1[3:]) //6 low
	// 字符串切片与原串共享空间

	fmt.Println( // 原生的字符串面值
		` hello \n /* \n和本注释都不会起效*/
		world `,
	)

	var s2 string = "h你好"
	fmt.Println(len(s2))                    // 7。按字节
	fmt.Println(utf8.RuneCountInString(s2)) // 3。按Unicode字符

	var n = 0
	for i := 0; i < len(s2); {
		ch, sz := utf8.DecodeRuneInString(s2[i:]) // 解码器。得到一个unicode字符(可变长)及其长度
		// ch将是rune类型
		fmt.Printf("%d  %c\n", i, ch) // %c也可以处理unicode字符
		i += sz
		n++ // 统计实际长度
	}
	/*
		0  h
		1  你
		4  好
	*/

	n = 0
	for i, ch := range "你好" { // range循环可以自动隐式解码UTF8字符串
		_, _ = i, ch
		n++
	}
	fmt.Println(n) // 2

	var r = []rune(s2) // rune类型等价于int32
	// []rune类型的Unicode字符slice或数组转为string，会作UTF-8编码(总空间变小)
	fmt.Println(string(r))             // h你好
	fmt.Println(string(rune(1234567))) // "�"。无效字符用\uFFFD作为替换

	x, _ := strconv.Atoi("123")             // 字符串转数字(Ascii -> Int)
	y, _ := strconv.ParseInt("123", 10, 64) //转十进制int64
	fmt.Println(x, y)                       //123 123

	var s3 = fmt.Sprintf("%d", 456)    // 使用Sprintf生成带额外信息的字符串，更为灵活
	fmt.Println(s3, strconv.Itoa(456)) // 456 456 (Int -> Ascii)

	fmt.Println(comma("123456789"))   // 123,456,789
	fmt.Println(comma("-123456.789")) // -123,456.789

	const (
		a1 = 1 // 不可省略初始化
		a2     // 会复制上面最近的初始化表达式
		a3 = 2
		a4 // a4 = 2
	)

	fmt.Println(a1, a2, a3, a4) // 1 1 2 2

	const (
		b1 = 1 << iota // iota初始为0
		b2             // 仍会复制，但是iota++
		b3             // 值为 1<<2
	)

	fmt.Println(b1, b2, b3) // 1 2 4

	const pi float64 = 3.141592653589                                           // 指明类型的常量
	const PI = 3.14159265358979323846264338327950288419716939937510582097494459 // 无类型常量
	fmt.Println((1 << 100) / (1 << 90))                                         // 1024。无类型常量的运算精度很高 (>=256)
	//var f1 float32 = pi // 错误。类型不匹配
	var f1 float32 = PI // 赋予合理的无类型常量无需作显式类型转换
	f2 := PI            // f2将取无类型浮点数常量默认类型float64

	_ = v1
	_ = v2
	_ = f1
	_ = f2
	_ = pi
}

func b2i(b bool) int {
	if b {
		return 1
	}
	return 0
}

func comma(s string) string { // 为数字添加逗号分隔。支持负数和浮点数
	var buf bytes.Buffer                //bytes.Buffer类型可以动态增长，按字节存值
	var neg = strings.LastIndex(s, "-") // 若不存在，返回-1
	var dot = strings.LastIndex(s, ".")
	var n = len(s) // 非小数部分
	if dot != -1 { // 有小数部分
		n = dot
	}
	var i = n % 3  // 遍历起始点
	if neg != -1 { // 有负号
		i = (n-1)%3 + 1
	}
	for lst := 0; i <= n; i += 3 {
		// s[i]是byte类型，对应WriteByte
		buf.WriteString(s[lst:i]) // 追加字符串
		if i > neg+1 && i < n {
			buf.WriteRune(',') // WriteRune可添加任意unicode字符
		}
		lst = i
	}
	buf.WriteString(s[n:])
	return buf.String()
}
