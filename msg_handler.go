package main

import (
	"github.com/bitly/go-nsq"
)

type MsgHandler interface {
	GetTopic() string
	GetChannel() string
	SetReader(r *nsq.Reader)
	HandleMessage(m *nsq.Message, responseChannel chan *nsq.FinishedMessage)
}
