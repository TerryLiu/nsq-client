package main

import (
	"github.com/bitly/go-nsq"
)

type MsgPublisher struct {
	usr    User
	topic  string
	writer *nsq.Writer
}

func NewMsgPublisher(_topic string, _usr User) (*MsgPublisher, error) {
	msgPublisher := &MsgPublisher{
		topic:  _topic,
		usr:    _usr,
		writer: nsq.NewWriter("106.186.31.48:4150"),
	}
	return msgPublisher, nil
}

func (mp *MsgPublisher) Publish(msg string) (int32, []byte, error) {
	return mp.writer.Publish(mp.topic, []byte(msg))
}

func (mp *MsgPublisher) PublishAsync(msg string, responseChan chan *nsq.WriterTransaction) error {
	return mp.writer.PublishAsync(mp.topic, []byte(msg), responseChan, "")
}

func (mp *MsgPublisher) Stop() {
	mp.writer.Stop()
}
