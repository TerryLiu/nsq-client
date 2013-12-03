package main

import (
	"github.com/bitly/go-nsq"
)

type NsqMsg struct {
	*nsq.Message
	returnChannel chan *nsq.FinishedMessage
}

type MsgHandler struct {
	topic   string
	channel string
	usr     User
	msgChan chan *NsqMsg
}

func (handler *MsgHandler) HandleMessage(m *nsq.Message, responseChannel chan *nsq.FinishedMessage) {
	handler.msgChan <- &NsqMsg{m, responseChannel}
}
