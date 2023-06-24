package main

import (
	"fmt"
	"log"
)

// 以下三个函数具有相同的函数签名：func(int, int) int
func f1(a int, b int) int { return a + b }
func f2(int, int) int     { return 0 }
func f3(x, y int) (z int) { z = x + y; return z }

func f4() (bool, error) { // 多返回值
	ok := true
	if ok {
		log.Printf("no bug") // 2023/05/01 20:50:14 no bug (log包的函数会为没有换行符的字符串增加\n)
		return true, nil     // 以nil表示无错误
	} else {
		return false, fmt.Errorf("this is a error") // 返回error类型。格式字符串作为错误信息
	}
}

func f6(x, y int, f func(int, int) int) int { // 传入一个函数类型值
	return f(x, y)
}

func f7() func() int { // 返回一个匿名函数(闭包)
	var x int
	return func() int { // 匿名函数
		x++ // 引用了x。闭包延长了f7函数体中局部变量的生存期
		return x
	}
}

func max(x int, nums...int) int { // 可变参数函数。nums是切片类型，包含传入的多个参数
	if len(nums) == 0 {
		return x
	}
	var y = max(nums[0], nums[1:]...)
	if x > y {
		return x
	}
	return y
}


func main() {
	fmt.Printf("%T\n%T\n%T\n", f1, f2, f3)
	f4()

	var f5 func(int, int) int // 声明一个函数值
	// 函数类型仅支持与nil比较
	fmt.Println(f5 == nil) //true。此时不可调用
	f5 = f1
	//f5 = f4 // 错误。函数签名不同
	fmt.Println(f6(1, 2, f5)) // 3
	fmt.Printf("%*d\n", 5, 1)   // *表示先输出若干个空格，这里表示输出5个空格后再输出1

	lam := f7()
	// 每次的返回值增加1
	fmt.Println(lam()) // 1
	fmt.Println(lam()) // 2

	ar := [...]int{4:0}
	sl := []int{0,2,4} // 由sl中指定的下标，将ar数组某些值置为-1
	var fs []func() // 函数切片
	for _,v := range sl {
		// 错误做法，将得到 [0 0 0 0 -1 0]
		// fs = append(fs, func() {
		// 	ar[v] = -1  // 闭包中引用着v这个变量，而v在存在于整个循环
		// })

		val := v
		// v := v 也可以
		fs = append(fs, func() {
			ar[val] = -1 // val的生存期被延长
		})
	}
	for _,f:=range fs { 
		f() // 调用闭包(main函数中)
	}
	fmt.Println(ar) // [-1 0 -1 0 -1 0]

	fmt.Println(max(1,3,4,2))
}
