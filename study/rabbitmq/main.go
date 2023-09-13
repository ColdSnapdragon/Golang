package main

import (
	"github.com/streadway/amqp"
	"log"
	"time"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

const MQURL = "amqp://guest:guest@localhost:5672/"

// (3)(4)尚有点问题
func main() {
	conn, err := amqp.Dial(MQURL) // 连接到RabbitMQ服务器
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// 创建一个channel来传递消息
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// (1)工作队列。Producer直接往工作队列发消息
	// 声明队列是幂等的——只有在它不存在(名字Name)的情况下才会创建它。消息内容是一个字节数组，因此可以编写任何内容
	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable (是否持久化)
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	// 发送消息
	err = ch.Publish(
		"",     // exchange (空字符串代表默认或者匿名交换机，消息将会根据指定的routing_key分发到指定的队列)
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte("Hello World!"),
		})
	failOnError(err, "Failed to publish a message")

	go worker()

	// (2)订阅/发布。Producer往交换机发生消息，队列的声明由消费者实现

	go consumer0("logs")
	go consumer("logs_direct", "direct", "blue")
	go consumer("logs_topic", "topic", "red.#")
	go consumer("logs_topic", "topic", "*.*.big")
	time.Sleep(10 * time.Millisecond) // 要等消费者创建好队列，再让交换机发消息

	// 声明一个具名交换机
	err = ch.ExchangeDeclare(
		"logs",   // name
		"fanout", // type (fanout(散开)模式会把消息发给所有队列)
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	err = ch.Publish(
		"logs", // exchange
		"",     // routing key (fanout交换机会忽略routing key(绑定键))
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte("everything is ok"),
		})
	failOnError(err, "Failed to publish a message")
	// Worker和Consumer的队列都会收到消息

	// (3)路由(Routing)模式
	// 采用direct交换机

	err = ch.ExchangeDeclare(
		"logs_direct", // name
		"direct",      // type (直连交换机。交换机将会对binding key和routing key进行精确匹配，从而确定消息该分发到哪个队列)
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	err = ch.Publish(
		"logs_direct", // exchange
		"red",         // routing key (fanout交换机会忽略routing key(绑定键))
		false,         // mandatory
		false,         // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte("direct is ok"),
		})
	failOnError(err, "Failed to publish a message")

	// (4)主题(topic)模式
	// 采用topic交换机

	err = ch.ExchangeDeclare(
		"logs_topic", // name
		"topic",      // type (发送到topic交换机的消息携带的routing_key是一个由.分隔开的词语列表)
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	err = ch.Publish(
		"logs_topic",   // exchange
		"blue.123.big", // routing key (fanout交换机会忽略routing key(绑定键))
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte("topic is ok"),
		})
	failOnError(err, "Failed to publish a message")

	select {}

}

func Create(name string, conn *amqp.Connection) (amqp.Queue, *amqp.Channel) {

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	q, err := ch.QueueDeclare(
		name,  // name
		false, // durable
		false, // delete when usused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")
	return q, ch
}

func worker() {
	conn, err := amqp.Dial(MQURL)
	failOnError(err, "Failed to connect to RabbitMQ")
	q, ch := Create("hello", conn)
	defer conn.Close()
	defer ch.Close()

	msgs, err := ch.Consume( // 持续将消息发送到返回的chan Delivery (也就是这里的msgs) 中
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack (自动对消息做ack确认，即接收后立即从队列删除)
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to consume a message")
	for d := range msgs { // 消费消息
		log.Printf("Worker: Received a message: %s", string(d.Body))
	}
}

func consumer0(exchange string) {
	conn, err := amqp.Dial(MQURL)
	failOnError(err, "Failed to connect to RabbitMQ")
	ch, err := conn.Channel()
	defer conn.Close()
	defer ch.Close()

	// 消费者也必须声明一下交换机
	err = ch.ExchangeDeclare(
		exchange, // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	// 一个临时队列的例子
	q, err := ch.QueueDeclare(
		"",    // name (将由RabbitMQ生成随机队列名称)
		false, // durable
		false, // delete when usused
		true,  // exclusive (true表示，当声明它的连接关闭时，队列将被删除)
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	// 绑定队列到交换机。同一队列可绑定多次，每个绑定都会被视为独立的路由规则
	err = ch.QueueBind(
		q.Name,   // queue name
		"",       // routing key
		exchange, // exchange
		false,
		nil,
	)
	failOnError(err, "Failed to bind a queue")

	msgs, err := ch.Consume( // 持续将消息发送到返回的chan Delivery (也就是这里的msgs) 中
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack (自动对消息做ack确认，即接收后立即从队列删除)
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to consume a message")
	for d := range msgs { // 消费消息
		log.Printf("Consumer1: Received a message: %s", string(d.Body))
	}
}

func consumer(exchange string, kind string, routingKey string) {
	conn, err := amqp.Dial(MQURL)
	failOnError(err, "Failed to connect to RabbitMQ")
	q, ch := Create("hello", conn)
	defer conn.Close()
	defer ch.Close()

	err = ch.ExchangeDeclare(
		exchange, // name
		kind,     // type (直连交换机。交换机将会对binding key和routing key进行精确匹配，从而确定消息该分发到哪个队列)
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	err = ch.QueueBind(
		q.Name,     // queue name
		routingKey, // routing key
		exchange,   // exchange
		false,
		nil,
	)
	failOnError(err, "Failed to bind a queue")

	msgs, err := ch.Consume( // 持续将消息发送到返回的chan Delivery (也就是这里的msgs) 中
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack (自动对消息做ack确认，即接收后立即从队列删除)
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to consume a message")
	for d := range msgs { // 消费消息
		log.Printf("%s -> routing(%s): Received a message: %s", d.RoutingKey, routingKey, string(d.Body))
	}
}
