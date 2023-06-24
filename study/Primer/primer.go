package main

import "fmt"

func main() {
	var a = 10 //自动识别为整数
	//a := 10  //与上一句等价
	const b = "hi" //自动识别为字符串常量
	//定义了变量或引入了包，但不使用，被认为是错误
	fmt.Println("a is", a, "; b is", b)                    //多参数(默认间隔空格)
	fmt.Printf("a is %v ; b is %v \n", a, b)               //格式化打印
	fmt.Printf("a's type is %T ; b's type is %T \n", a, b) //%T可以获取类型

	//var c  //没有初始值用于推导类型时，必须显式指明
	var c string

	fmt.Println("Input something :")
	fmt.Scan(&c)

	fmt.Print("your input : ", c, "\n") //Print的参数之间默认无空格

	//var ar [10] string //只定义
	var ar [10]string = [10]string{"aa", "bb", "cc", 5: "ee"} //定义并初始化(其中赋值了下标5)
	//var ar = [10] string {"aa", "bb", "cc", 5:"ee"}  //简化写法，与上一行等价
	ar[0] = b + " " + c //与js等不同，变量b,c必须定义后才可以使用

	fmt.Printf("ar's length : %v, ar's type : %T \n", len(ar), ar) //len()为内置函数

	var sl []int        //定义切片类型(注意与数组的区别)
	sl = append(sl, 10) //追加，并重新赋值

	fmt.Printf("sl's length : %v, sl's type : %T \n", len(sl), sl)
}

//go fmt -help
//gofmt -w .\primer.go  //格式化代码