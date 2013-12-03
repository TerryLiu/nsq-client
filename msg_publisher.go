package main

import (
	"encoding/json"
	"github.com/bitly/go-nsq"
	"log"
	"time"
)

var (
	Publisher *MsgPublisher
)

type MsgPublisher struct {
	usr    User
	writer *nsq.Writer
}

func init() {
	Publisher = &MsgPublisher{
		writer: nsq.NewWriter("106.186.31.48:4150"),
	}
}

func (publisher *MsgPublisher) SetLoginUsr(_usr User) error {
	if Publisher == nil {
		Publisher = &MsgPublisher{
			writer: nsq.NewWriter("106.186.31.48:4150"),
		}
	}
	Publisher.usr = _usr
	return nil
}

func (mp *MsgPublisher) Write(topic, msg string) (int32, []byte, error) {
	m := Message{
		Topic: topic,
		Type:  MSG_TYPE_CHAT,
		Time:  time.Now().Format("2006-01-02 15:04:05"),
		Body: MessageBody{
			From: mp.usr,
			Msg:  msg,
		},
	}
	msgBytes, _ := json.Marshal(m)
	log.Println("Publish.msg = ", string(msgBytes))
	return mp.writer.Publish(topic, msgBytes)
}

func (mp *MsgPublisher) WriteAsync(topic, msg string, responseChan chan *nsq.WriterTransaction) error {
	m := Message{
		Topic: topic,
		Type:  MSG_TYPE_CHAT,
		Body: MessageBody{
			From: mp.usr,
			Msg:  msg,
		},
	}
	msgBytes, _ := json.Marshal(m)
	log.Println("PublishAsync.msg = ", string(msgBytes))
	return mp.writer.PublishAsync(topic, msgBytes, responseChan, "")
}

func (mp *MsgPublisher) Stop() {
	mp.writer.Stop()
}
