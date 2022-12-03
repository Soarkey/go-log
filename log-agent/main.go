package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"

	"go-log/log-agent/kafka"
	"go-log/log-agent/tail"
)

// AutoWriteLog 模拟写入日志
func AutoWriteLog(filename string) {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	w := bufio.NewWriter(f)

	logs := []string{
		"[DEBUG] 这是一条日志",
		"[INFO] 这是一条日志",
		"[WARN] 这是一条日志",
		"[ERROR] 这是一条日志",
	}

	ticker := time.NewTicker(2 * time.Second)
	for range ticker.C {
		// 每2s随机写入一条日志
		log := fmt.Sprintf("[%+v] %+v\n", time.Now(), logs[rand.Int()%len(logs)])
		w.WriteString(log)
		w.Flush()
	}
}

const (
	TOPIC    = "go-log"
	LOG_PATH = "demo.log"
)

func main() {
	os.Create(LOG_PATH)

	k := kafka.NewProducer([]string{"kafka:9092"})
	t := tail.New(LOG_PATH)

	go AutoWriteLog(LOG_PATH)

	for {
		select {
		case line := <-t.ReadChan():
			// 从通道中读取写入的日志并发送到kafka
			k.SendToKafka(TOPIC, line.Text)
		default:
			// 避免CPU空转
			time.Sleep(500 * time.Millisecond)
		}
	}
}
