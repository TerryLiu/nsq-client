package main

import (
	"encoding/json"
	"github.com/bitly/go-nsq"
	"log"
)

type MsgPublisher struct {
	usr    User
	writer *nsq.Writer
}

func NewMsgPublisher(_usr User) (*MsgPublisher, error) {
	msgPublisher := &MsgPublisher{
		usr:    _usr,
		writer: nsq.NewWriter("106.186.31.48:4150"),
	}
	return msgPublisher, nil
}

func (mp *MsgPublisher) Publish(topic, msg string) (int32, []byte, error) {
	m := Message{
		Topic: topic,
		Type:  MSG_TYPE_CHAT,
		Body: MessageBody{
			From: mp.usr,
			Msg:  msg,
		},
	}
	msgBytes, _ := json.Marshal(m)
	log.Println("Publish.msg = ", string(msgBytes))
	return mp.writer.Publish(topic, msgBytes)
}

func (mp *MsgPublisher) PublishAsync(topic, msg string, responseChan chan *nsq.WriterTransaction) error {
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
