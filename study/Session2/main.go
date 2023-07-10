package main //描述该源文件属于哪个包

import ( //导入依赖的其它包
	"fmt"       // $GOPATH/src/fmt
	"math/rand" // 导入math下的rand，之后用rand.调用函数

	"practise/Session2/color" // 使用go module导入本地包。第一个名字是go.mod中module对象(项目名or总包名)
	// "./color" // 有时候可以这样，有时候又不行
	// 对于红波浪线，直接启动go run
)

// import "github.com/username/project/utils" // 从远程仓库导包

// 四大声明：var、const、type和func
const pi = 3.14 //包一级声明语句声明的名字可在整个包对应的每个源文件中访问

func main() {
	fmt.Println(rand.Intn(10))

	//var 变量名字 类型 = 表达式。其中“类型”或“= 表达式”两个部分可以省略其中的一个
	//go语言中不存在未初始化的变量，只声明的变量会被零初始化

	var v1, v2, v3 int                 //显式声明三个int
	var v4, v5, v6 = true, pi, "hello" //推导三个类型

	//var变量通常先显式声明后赋值，简单变量声明常用于临时变量
	sh1, sh2 := 0, 1

	sh1, sh2 = sh2, sh1 //支持元组赋值

	sh2, sh3 := sh1, sh2 //同级作用域下已有sh2，则重新赋值，而sh3被声明与初始化（:=至少要声明一个变量）

	var p1 *int            // 指针类型
	fmt.Println(p1 == nil) //true。空指针
	p1 = &v1
	*p1 = 1
	fmt.Println(*p1) // 1

	p2 := new(int) // 创建一个T类型变量，返回其地址。
	// new只是个语法糖，相当于简化了 var x int; p = &x;

	type A float32 // 定义新类型
	type B float32 // 与A的本质类型相同

	var ( // 声明多个var
		a A = 1.1
		b B = 1.1
	)

	//println(a == b) // 错误。a b类型不匹配。这样可避免无意的类型混用
	a = A(b) // 每个类型都有转型函数，只改变语义，不影响值
	//a = float32(b) // 错误，赋值时左右类型不匹配
	//var y int = 1.1 //同上，非法语句

	fmt.Println(color.GetColor()) // 通过包名调用可见的函数

	_ = v1
	_ = v2
	_ = v3
	_ = v4
	_ = v5
	_ = v6
	_ = sh1
	_ = sh2
	_ = sh3
	_ = p2
	_ = p3
	_ = a
	_ = b
	f()
}

var p3 *int

func f() {
	var x = 1 //一个变量的有效周期只取决于是否每个包级是否有可达路径(通过指针或引用)
	p3 = &x   // x从函数中“逃逸”，编译器会自动选择将其分配在堆上，由GC回收
}
