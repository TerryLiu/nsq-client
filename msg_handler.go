package main

import (
	"github.com/bitly/go-nsq"
	"log"
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
	reader  *nsq.Reader
}

func (handler *MsgHandler) HandleMessage(m *nsq.Message, responseChannel chan *nsq.FinishedMessage) {
	log.Printf("HandleMessage...topic=%s\n", handler.topic)
	log.Printf("HandleMessage...m.Body=%s\n", string(m.Body))
	handler.msgChan <- &NsqMsg{m, responseChannel}
}
