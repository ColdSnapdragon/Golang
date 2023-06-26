package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sort"
)

type Count struct {
	n int
}

func (c *Count) Write(b []byte) (n int, err error) { // 实现io.Writer接口
	c.n = len(b)
	return c.n, nil
}

func main() {
	var c1 Count
	//Fprintf的第一个参数是io.Writer接口类型，要求参数拥有对应签名的方法
	fmt.Fprintf(&c1, "hi") // 注意Write方法是*Count类型的，必须传&c1
	fmt.Println(c1.n)      // 2

	var w io.Writer
	w = os.Stdout         // os.Stdout有Write方法
	w = new(bytes.Buffer) //*bytes.Buffer有Write方法
	w = (*Count)(nil)     // 定义一个值为nil的Count指针
	//fmt.Println(w == nil) // false。注意，虽然w的接口值为nil，但是w本身非nil，其已经有了动态类型

	type ReadWrite interface { // 接口类型具体描述了一系列方法的集合
		io.Writer // 接口内嵌，组合已有的接口
		Read([]byte) (int, error)
	}

	var rw ReadWrite
	w = rw // 把一个接口值赋值给另一个
	//rw = w // 错误。w没有Read()

	type Any interface{} // 空接口类型接受任何值
	var any Any
	any = 1
	any = "1"
	any = func() {}

	var s1 = []string{"bb", "cc", "abc"}
	var s2 = []int{3, 4, 1}
	sort.Strings(s1)                    // 切片s1没有实现sort的接口，不能直接用Sort()，于是专门有了Strings()给string切片用
	sort.Ints(s2)                       // O(nlogn)
	fmt.Println(s1)                     // [abc bb cc]
	fmt.Println(s2)                     // [1 3 4]
	var si = sort.IntSlice(s2)          // 显式将[]int转换为sort.IntSlice接口
	sort.Sort(si)                       // 排序
	fmt.Println(sort.IntsAreSorted(si)) //true。判断是否为升序
	sort.Sort(sort.Reverse(si))         //将一个接口的Less函数逆转，生成另一个接口(底层引用相同的切片)
	fmt.Println(s2)

	w = os.Stdout
	// 若断言成功，返回类型与断言类型一致
	t1 := w.(*os.File)      // 成功，t1为*os.File类型，值为w的接口值
	t2 := w.(io.ReadWriter) // 成功，t2为io.ReadWriter接口，其接口值(动态类型与动态值)和w一致
	//t3:=w.(sort.Interface) // 失败。会在编译时爆panic
	t3, ok := w.(sort.Interface) // 若用两个值接收，则把结果传给ok，代替报panic
	fmt.Println(ok)              // false

	_ = w
	_ = any
	_ = t1
	_ = t2
	_ = t3
}
