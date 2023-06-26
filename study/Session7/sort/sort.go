package main

import (
	"fmt"
	"sort"
)

type Point struct {
	X, Y int
}

// 不影响Point的前提下，为Point切片排序

//sort.Interface需要三个函数: Len()int、Less(i,j int)bool、Swap(i,j int)

type seque []Point

func (sq seque) Len() int {
	return len(sq)
}

func (sq seque) Less(i,j int) bool {
	if sq[i].X==sq[j].X {
		return sq[i].Y<sq[j].Y
	}
	return sq[i].X<sq[j].X
}

func (sq seque) Swap(i,j int) {
	sq[i],sq[j] = sq[j],sq[i]
}

func main() {

	var sl = []Point{
		{3,1},{5,2},{2,4},{1,7},{2,3},
	}

	sort.Sort(seque(sl)) // 要先包装成seque类型
	fmt.Println(sl) // [{1 7} {2 3} {2 4} {3 1} {5 2}]

}