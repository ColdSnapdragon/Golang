package main

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	reader *kafka.Reader
	topic  = "my_topic"
)

// 需要创建所谓的writer和reader来写入或读取数据

// 生产消息
func WriteKafka(ctx context.Context) {
	write := kafka.Writer{
		Addr:         kafka.TCP("localhost:9092"),
		Topic:        topic,           // 如果这里没有定义topic，那么在写消息时必须显式指定topic
		Balancer:     &kafka.Hash{},   // 负载均衡算法(哈希取余)
		RequiredAcks: -1,              // 三种返回ack的方式。默认0(不安全)
		WriteTimeout: 1 * time.Second, // 阻塞等待写入的时限
		Async:        false,           // 是否阻塞(true仅当不关心写入成功)
		// AllowAutoTopicCreation: true, //true表示没有topic时自动创建(但最好是运维在命令行上去操作)
	}
	defer write.Close()
	for i := 0; i < 3; i++ { // 失败时最多重复三次
		if err := write.WriteMessages(ctx,
			// 默认阻塞直至全部消息写入(但此函数没有事务性)
			kafka.Message{Value: []byte("hi")},                                                      // 轮换分区
			kafka.Message{Key: []byte("msg"), Value: []byte("welcome")},                             // Key哈希取余得分区
			kafka.Message{Topic: topic, Partition: 0, Key: []byte("nick"), Value: []byte("blover")}, // 指定分区
		); err != nil {
			if err == kafka.LeaderNotAvailable {
				fmt.Println(topic, "不存在") // 一种常见情况
			} else {
				fmt.Printf("写入Kafka失败：%v\n", err)
			}
		} else {
			break // 写入成功
		}
	}
}

// 消费消息
func ReadKafka(ctx context.Context) {
	reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{"localhost:9092"}, // 实际项目中会有多个broker
		Topic:          topic,
		GroupID:        "my_team",         // 启用消费者组
		CommitInterval: 1 * time.Second,   // 消费者组自动向kafka提交offset的时间间隔(越长越快)
		StartOffset:    kafka.FirstOffset, // 指定新的consumer第一次从哪开始消费(这里是从头开始)
	})
	defer reader.Close() // 在进程退出时调用Reader的Close()方法非常重要，要处理进程退出的情况

	for {
		if msg, err := reader.ReadMessage(ctx); err != nil { // 没有消息时默认阻塞等待
			fmt.Printf("%#v\n", msg)
		} else {
			fmt.Println("读取Kafka失败")
			break
		}
	}
}

func listenSignal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGTERM) // 监听中断信号和终止信号
	// 优雅终止: 发送终止信号并给予时间让进程做收拾工作，必要时才发送强制终止
	<-c
	_ = reader.Close()
	os.Exit(0)
}

func main() {
	ctx := context.Background()
	WriteKafka(ctx)
	go listenSignal()
	ReadKafka(ctx)
}
