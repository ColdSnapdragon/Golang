package main

import (
	"fmt"
)

func main() {
	ar1 := [3]int{1, 2: -1}              // 指定2号索引的元素值
	ar2 := [...]int{99: 1}               //编译时自动确定长度
	fmt.Println(len(ar2))                // 100
	fmt.Println(ar1 == [3]int{1, 0, -1}) //true。按值比较等类型等长度的数组(== !=)
	for i, v := range ar1 {              // range生成连续的(索引,值)
		fmt.Println(i, v)
	}

	ar := [10]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	fmt.Printf("%T\n", ar[1:2]) //[]int。即切片类型
	//切片是其底层数组的视图
	var sl1 []int
	fmt.Println(sl1 == nil)       //true。切片唯一支持的等于运算，表示指向底层数组的指针为空
	sl1 = ar[1:2]                 //左开右闭
	fmt.Println(sl1[0], len(sl1)) // 1 1
	ar[1] = 2
	fmt.Println(sl1[0]) // 2。切片的底层数组发生改变
	sl2 := sl1          //引用传递(浅拷贝)
	sl2[0] = 3
	fmt.Println(sl1[0])    // 3
	sl1 = sl1[:len(sl1)+1] // 拓宽视图
	fmt.Println(sl1)       // [3 2]

	sl4 := make([]int, 2, 4) // 开一个总容为4的数组，建立长度为2的切片
	sl4 = append(sl4, 1)     // 安全起见，应当接收返回值(因为可能会得到新视图)
	fmt.Println(sl4)         // [0 0 1]

	sl4 = sl4[0:0] // 一种清空方式，这样保留了底层数组，不必从头扩容
	sl4 = append(sl4, sl1...) // ...将对象拆解为一个个元素，而append支持一次添加多个对象
	fmt.Println(sl4) // [3 2]

	var mp1 map[string]int        // 定义一个map
	fmt.Printf("%d\n", mp1["a"])  // 不存在的键默认返回对应零值(允许访问值为nil的map，但不能修改)
	if val, ok := mp1["a"]; !ok { // 第二个值表示键是否存在(bool)
		fmt.Println("key not exist")
	} else {
		fmt.Println(val)
	}
	fmt.Println(mp1 == nil) // true。此时mp1还未真正创建(未引用哈希表)
	//mp1["a"] += 1 // 错误。

	mp1 = make(map[string]int) // make创建map
	mp1["a"]++
	fmt.Println(len(mp1))
	delete(mp1, "a") // 内置的delete函数，删除一个键

	for k, v := range mp1 { // 遍历map
		_ = k
		_ = v
	}

	//同slice一样，map拷贝时传递引用，且只能与nil比较，因为不能作为键
	graph := map[string]map[string]bool{ // 字面量创建map
		"school": {"home": true},
	}
	st, ed := "aa", "bb"  // 给图添加一条边
	if graph[st] == nil { // slice、map的“零值”是nil
		graph[st] = make(map[string]bool)
	}
	graph[st][ed] = true

	G := map[string][]string {
		"school" : {"home1", "home2"},
	}
	G["school"] = append(G["school"], "home3")

	/*-----------------------------------------------------------------*/

	type Person struct { // 定义一个结构体类型
		age  int    // 不需要var
		Name string // 首字母大写表示在外部包中可见
	}

	var p1 Person
	p1.age = 1
	var pp *Person = &p1
	pp.age = 2    // 结构体指针也可以直接用点运算符访问成员(就像c的->)
	(*pp).age = 2 // 与上句等价

	type A struct {
		x int
		y int
	}

	a := A{1, 2} // 结构体字面值
	_ = a
	a = A{x: 3} // 以成员名字和相应的值来初始化
	var pa *A = new(A)
	fmt.Println(pa == nil) // false。new新创建了变量
	pa = &A{3, 4}          // 取字面值的地址

	// 结构体的成员均可比较，则结构体便支持==，也就可以作为键
	var mp map[A]struct{} // 偶尔这样去实现set: 值类型为空结构体

	type B struct {
		A     // 匿名成员。A的属性被嵌入到B中
		pb *B // 成员不能是当前结构体，但可以是其指针
	}

	var b1 B
	b1.x = 1
	b1.A.x = 1 // 等价。A是可选的
	//b2 := B{1,2} // 错误。匿名成员不支持结构体字面值
	//b2 := B{x:1, y:2} // 错误。同上
	b2 := B{A: A{1, 2}, pb: &b1} // 正确写法
	_ = b2
	b2 = B{A{1, 2}, &b1} // 等价

	fmt.Printf("%#v\n", b2) // %v是打印变量(类型确定)，#表示对于结构体要携带结构体名和字段名
	//输出：main.B{A:main.A{x:1, y:2}, pb:(*main.B)(0xc0000080c0)}

	var t = struct { // 定义一个匿名结构给一个变量
		x int
		mp map[int]int
	} {
		mp : make(map[int]int),
	}

	_ = ar1
	_ = ar2
	_ = sl4
	_ = a
	_ = pa
	_ = mp
	_ = b2
	_ = t
}
