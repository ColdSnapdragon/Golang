//所有以_test.go为后缀名的源文件在执行go build时不会被构建成包的一部分，它们是go test测试的一部分

package test1

import "testing"

//go test -v -run="."  run接受正则表达式，只执行与之匹配的测试函数(这里是匹配所有)
//测试函数的名字必须以Test开头，可选的后缀名必须以大写字母开头

func TestAdd(t *testing.T) {
	var cases = []struct {
		x    int
		y    int
		want int
	}{
		{1, 1, 2},
		{2, 2, 4},
		{3, 3, 6}, // 故意给一个错的
		{4, 4, 8},
	} // 这种表格驱动的测试在Go语言中很常见

	for _, test := range cases {
		if res := add(test.x, test.y); res != test.want {
			// t.Errorf调用也没有引起panic异常或停止测试的执行。需要的话可以用t.Fatal
			t.Errorf("add(%v, %v) = %v, want: %v\n", test.x, test.y, res, test.want)
		}
	}
}

// 如果某一个测试函数失败，那么测试过程将会停止，不会继续执行基准测试。

// 基准测试函数和普通测试函数写法类似，但是以Benchmark为前缀名
// go test -bench="Add|Sub"  bench接受正则表达式

func BenchmarkAdd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		add(1, 1)
	}
}
