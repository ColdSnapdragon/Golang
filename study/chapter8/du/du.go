package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath" // 提供对路径的操作
	"sync"
	"time"
)

// 1 计算传入文件的总大小，并发实现
// 2 当指定-v选项时，每隔一秒打印当前扫描文件的数量以及总大小
// 3 按下回车键立即结束程序，并妥善做清理

var flag_v = flag.Bool("v", false, "输出详细执行状态") // -v要作为第一个参数才会被flag认为是选项

var done = make(chan struct{}) // 通过关闭done，实现向所有goroutine广播

// 检查程序是否被取消
func cancelled() bool {
	select { // 通过select轮询，可以在不产生阻塞的情况下检查channel是否可读
	case <-done: // 在关闭done后才可触发
		return true
	default:
		return false
	}
}

func main() {
	flag.Parse()

	go func() {
		os.Stdin.Read(make([]byte, 1))
		close(done)
	}()

	var filesizes = make(chan int64)
	var wg sync.WaitGroup // 统计所有正在进行的walkDir函数数量

	for _, arg := range flag.Args() {
		wg.Add(1)
		go walkDir(arg, &wg, filesizes)
	}

	go func() { // closer
		wg.Wait()
		close(filesizes) // 必须及时关闭filesizes，避免之后的for range一直阻塞
	}()

	var tick <-chan time.Time
	if *flag_v { // 很巧妙的设计。只有当-v为true时，才设置Tick向channel发数据
		tick = time.Tick(time.Second)
	}
	var filenum int
	var totsize int64
loop: // 只有当标签后紧随for时，可用break label跳出
	for {
		select { // 进行*轮询*(顺序随机) (不搞多路复用的话，这里也可以开goroutine)
		case sz, ok := <-filesizes: // 注意，即使关闭channel，依旧能读数据(得到零值)，只能用第二个返回值来判断关闭
			if !ok {
				break loop // 跳出外部循环
			}
			filenum++
			totsize += sz
		case <-tick: // -v为false时永不触发
			output(filenum, totsize)
		case <-done: // 检测到中止后，停止计时，进入结束流程
			for range filesizes { // “排空”操作。唤醒被filesizes阻塞的walkDir，使其正确地、尽快地完成 （但是由于要等walkDir都退出，实际上比不排空要慢）
				<-filesizes
			}
			return
		}
	}

	output(filenum, totsize)
}

func output(filenum int, totsize int64) {
	fmt.Printf("%d files, %.2f MB\n", filenum, float64(totsize)/float64(1<<20))
}

// 遍历一个目录
func walkDir(path string, wg *sync.WaitGroup, filesizes chan<- int64) { // <- chan: 只读  chan <-: 只写
	defer wg.Done()

	for _, de := range dirents(path) {
		if cancelled() {
			return
		}
		newpath := filepath.Join(path, de.Name()) // 使用filepath.Join来拼接路径更灵活与兼容
		if de.IsDir() {
			wg.Add(1)
			go walkDir(newpath, wg, filesizes)
		} else {
			info, err := de.Info()
			if err != nil {
				fmt.Fprintf(os.Stderr, "du: %v\n", err)
				continue
			}
			filesizes <- info.Size()
		}
	}
}

var sema = make(chan struct{}, 20) // 用一个channel作为semaphone来限制并发量

// 获取目录下的所有目录项
func dirents(path string) []os.DirEntry { // string类型是一个只读的字节数组切片
	if cancelled() {
		return nil
	}

	sema <- struct{}{}
	defer func() { <-sema }()

	ds, err := os.ReadDir(path) // os.ReadDir()返回目录项列表
	if err != nil {
		fmt.Fprintf(os.Stderr, "du: %v\n", err)
		return nil
	}
	return ds
}
