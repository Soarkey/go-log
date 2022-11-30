package kafka

import (
	"github.com/Shopify/sarama"
	"fmt"
)

type Consumer struct {
	client sarama.Consumer
	addrs  []string
	topic  string
}

// Start 启动消费
func (c *Consumer) Start() {
	partitions, err := c.client.Partitions(c.topic)
	if err != nil {
		panic(fmt.Sprintf("kafka 消费者获取分区列表失败 err: %+v\n", err))
	}
	fmt.Println("所有分区: ", partitions)
	for partition := range partitions {
		// 每一个分区启动一个消费者
		pc, err := c.client.ConsumePartition(c.topic, int32(partition), sarama.OffsetOldest)
		if err != nil {
			fmt.Printf("分区 %d 启动消费者失败 err: %+v", pc, err)
		}
		defer pc.AsyncClose()
		// 启动异步消费
		go func() {
			for msg := range pc.Messages() {
				fmt.Printf("分区 %d, Offset: %d Key: %s Value: %+v\n", msg.Partition, msg.Offset, msg.Key, string(msg.Value))
			}
		}()
	}
	defer c.client.Close()
	// 避免函数退出
	select {}
}

func NewConsumer(addrs []string, topic string) *Consumer {
	consumer, err := sarama.NewConsumer(addrs, nil)
	if err != nil {
		panic(fmt.Sprintf("kafka 消费者连接失败 err: %+v\n", err))
	}
	return &Consumer{
		client: consumer,
		addrs:  addrs,
		topic:  topic,
	}
}
