package main

import (
	"fmt"
	"github.com/streadway/amqp"
)

// 消费者订阅
func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		fmt.Println("连接失败，err:", err)
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			fmt.Println("连接关闭失败，err:", err)
		}
	}()
	ch, err := conn.Channel()
	if err != nil {
		fmt.Println("channel建立失败，err:", err)
	}
	defer func() {
		err := ch.Close()
		if err != nil {
			fmt.Println("channel关闭失败，err:", err)
		}
	}()
	//声明队列
	_, err = ch.QueueDeclare("queue1", true, false, false, false, nil)
	if err != nil {
		fmt.Println("声明queue1队列失败，err:", err)
	}
	_, err = ch.QueueDeclare("queue2", true, false, false, false, nil)
	if err != nil {
		fmt.Println("声明queue2队列失败，err:", err)
	}
	err = ch.ExchangeDelete("exchange1", false, false)
	//声明交换机
	if err != nil {
		fmt.Println("删除交换机失败，err:", err)
	}
	err = ch.ExchangeDeclare("exchange1", "topic", true, false, false, false, nil)
	if err != nil {
		fmt.Println("声明交换机失败，err:", err)
	}
	//绑定队列
	err = ch.QueueBind("queue1", "*.key1.*", "exchange1", false, nil)
	if err != nil {
		fmt.Println("绑定queue1队列失败，err:", err)
	}
	err = ch.QueueBind("queue2", "*.key2.*", "exchange1", false, nil)
	if err != nil {
		fmt.Println("绑定queue2队列失败，err:", err)
	}
	//发布消息
	for {
		s := ""
		key := ""
		_, _ = fmt.Scan(&s, &key)
		err = ch.Publish("exchange1", key, false, false, amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(s),
		})
		fmt.Println("发布成功,msg:", s)
		if err != nil {
			fmt.Println("发布失败，err:", err)
		}
	}
}
