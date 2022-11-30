package kafka

import (
	"fmt"

	"github.com/Shopify/sarama"
)

type Producer struct {
	client sarama.SyncProducer
	addrs  []string
}

// SendToKafka 处理日志并投递到kafka
func (p *Producer) SendToKafka(topic, data string) {
	// 构造消息
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(data),
	}
	// 发送消息
	pid, offset, err := p.client.SendMessage(msg)
	if err != nil {
		fmt.Printf("消息发送失败! err: %+v\n", err)
		return
	}
	fmt.Printf("消息发送成功! pid: %+v offset: %+v data: %+v\n", pid, offset, data)
	return
}

func NewProducer(addrs []string) *Producer {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true

	client, err := sarama.NewSyncProducer(addrs, config)
	if err != nil {
		panic(fmt.Sprintf("kafka 生产者连接失败 err: %+v\n", err))
	}
	return &Producer{
		client: client,
		addrs:  addrs,
	}
}
