package main

import (
	"fmt"
	"math"
)

type Point struct {
	X, Y float64
}

// 方法声明。接收器(receiver)类型为Point
func (p Point) newScale(c float64) Point { // 注意，p是调用者的拷贝
	p.X *= c
	p.Y *= c
	return p
}

func (p *Point) Scale(c float64) { // p是调用者的指针，类似其他语言的self
	if p == nil { // 允许接收器为nil
		return
	}
	p.X *= c
	p.Y *= c
}

func (p *Point) Dis(q Point) float64 {
	return math.Sqrt((p.X-q.X)*(p.X-q.X) + (p.Y-q.Y)*(p.Y-q.Y))
}

type PointX struct {
	*Point // 内嵌匿名命名类型指针，访问需要通过该指针指向的对象去取
	val    float64
}

func main() {
	var p1 = Point{1, 2}
	_ = p1.newScale(2)
	fmt.Println(p1) // p1没有改变
	p1.Scale(2)
	fmt.Println(p1) // {2 4}

	pp := &p1
	_ = pp.newScale(1) // 隐式转为(*pp).
	p1.Scale(1)        // 隐式转为(&p1). (与上一个都是语法糖而已,T类型的值不拥有*T类型的方法)

	px1 := PointX{&Point{1, 2}, 3}
	fmt.Println(px1.Dis(p1)) // 会从调用者的类型出发，bfs式逐层在所有内嵌类型中查找Dis方法
	// 这与继承、多态无关，函数参数仍必须传对应的Point类型，而不是PointX

	px2 := px1         // 共享指针所指对象
	px2.X = 4          // (*px2.Point).X = 4
	fmt.Println(px1.X) // 4

	var p2 = Point{1, 1}
	ft := p2.Scale         // 就像函数值一样，方法也有方法值
	fmt.Printf("%T\n", ft) // func(float64) 。该方法值的接收器已经被绑定为p2
	ft(2)
	ft(3)
	fmt.Println(p2) // {6 6}

	f2 := (*Point).Scale   // 方法表达式。必须指明是*Point类型的Scale方法
	fmt.Printf("%T\n", f2) // func(*main.Point, float64) 。方法表达式的第一个参数为接收器
	f2(&p2, 2)
	fmt.Println(p2) // {12 12}

}
