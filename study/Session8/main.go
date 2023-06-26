package main

import (
	"fmt"
	"os"
	"sync"
	"time"
)

func main() {
	ss := []string{"abc", "de"}
	fmt.Println(work1(ss))

	work2()
}

func work1(list []string) uint {
	sizes := make(chan uint) // 不带缓存的channel(同样需要用make分配空间，不然只能读不能写)
	var wg sync.WaitGroup    // sync.WaitGroup是go的一种并发控制方式，使用计数器，类似信号量

	for _, s := range list {
		wg.Add(1)
		go func(s string) {
			defer wg.Done() // 等价于wg.Add(-1)
			sizes <- handleStr(s)
		}(s) // 一定要注意闭包引用变量问题，这里拷贝了s给闭包
	}

	go func() { // closer
		wg.Wait()    // 阻塞等待计数器为0
		close(sizes) // 关闭通道(的写端)
	}() // 这段代码必须令开goroutine，否则放哪里都是错的

	var tot uint
	for sz := range sizes { // range可以从通道中取值，无值则阻塞，直至写端关闭
		tot += sz
	}

	return tot
}

func handleStr(str string) uint {
	return uint(len(str))
}

func work2() {
	// tick := time.Tick(time.Second) // 创建一个channel，每隔1s向通道发送值
	ticker := time.NewTicker(time.Second) // 功能全面的Ticker(我称作节拍)
	// time.Tick()源码仅仅是调用time.NewTicker后返回C成员

	abort := make(chan struct{}) // 空结构体类型的channel
	go func() {                  // 打断倒计时
		os.Stdin.Read(make([]byte, 1)) // 阻塞，从stdin读取一个byte
		abort <- struct{}{}
	}()

	defer fmt.Println("lauch")

	for count := 10; count > 0; count-- {
		select {
		case x := <-ticker.C:
			fmt.Println(count, "\t", x.Minute(), ":", x.Second())
		case <-abort:
			ticker.Stop() // 中止ticker，以免其继续工作
			return
			// 有多个case可以进入时，select将随机选择 (default在最后才考虑)
			// 由于没有default，select将阻塞直到有case可以进入
		}
	}
}
