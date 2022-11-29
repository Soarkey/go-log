package kafka

import (
	"fmt"

	"github.com/Shopify/sarama"
)

type Kafka struct {
	client sarama.SyncProducer
}

// SendToKafka 处理日志并投递到kafka
func (k *Kafka) SendToKafka(topic, data string) {
	// 构造消息
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(data),
	}
	// 发送消息
	pid, offset, err := k.client.SendMessage(msg)
	if err != nil {
		fmt.Printf("消息发送失败! err: %+v\n", err)
		return
	}
	fmt.Printf("消息发送成功! pid: %+v offset: %+v data: %+v\n", pid, offset, data)
	return
}

func New(addrs []string) *Kafka {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true

	client, err := sarama.NewSyncProducer(addrs, config)
	if err != nil {
		panic(fmt.Sprintf("kafka连接失败 err: %+v\n", err))
	}
	return &Kafka{client: client}
}
