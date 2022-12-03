package es

import (
	"context"
	"fmt"
	"time"

	"go-log/model"

	"github.com/olivere/elastic"
)

type ES struct {
	client  *elastic.Client
	addr    string
	workers int
	ch      chan *model.Log
}

func (e *ES) Start() {
	for i := 0; i < e.workers; i++ {
		go e.NewWorker()
	}
}

func (e *ES) NewWorker() {
	for {
		select {
		case msg := <-e.ch:
			resp, err := e.client.Index().Index(msg.Topic).
				BodyJson(msg).
				Do(context.Background())
			if err != nil {
				fmt.Printf("发送消息至ES失败 msg: %+v err: %+v\n", msg, err)
			}
			fmt.Printf("发送消息至ES成功 msg: %+v id: %+v index: %+v type %+v\n", msg, resp.Id, resp.Index, resp.Type)
		default:
			// 避免CPU空转
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func (e *ES) SendToES(log *model.Log) {
	e.ch <- log
}

func New(addr string, workers, chanMaxSize int) *ES {
	client, err := elastic.NewSimpleClient(elastic.SetURL(addr))
	if err != nil {
		panic(fmt.Sprintf("ES 连接失败 err: %+v\n", err))
	}
	fmt.Printf("ES 连接成功 addr: %s workers: %d\n", addr, workers)
	return &ES{
		client:  client,
		addr:    addr,
		workers: workers,
		ch:      make(chan *model.Log, chanMaxSize),
	}
}
