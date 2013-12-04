package main

import (
	"github.com/bitly/go-nsq"
	"log"
)

type NsqMsg struct {
	*nsq.Message
	returnChannel chan *nsq.FinishedMessage
}

type ChatMgr struct {
	topic   string
	channel string
	usr     User
	msgChan chan *NsqMsg
	reader  *nsq.Reader
}

func (mgr *ChatMgr) GetTopic() string {
	return mgr.topic
}
func (mgr *ChatMgr) GetChannel() string {
	return mgr.channel
}
func (mgr *ChatMgr) SetReader(r *nsq.Reader) {
	mgr.reader = r
}

func (mgr *ChatMgr) HandleMessage(m *nsq.Message, responseChannel chan *nsq.FinishedMessage) {
	log.Printf("ChatMgr.HandleMessage...topic=%s\n", mgr.topic)
	log.Printf("ChatMgr.HandleMessage...m.Body=%s\n", string(m.Body))
	mgr.msgChan <- &NsqMsg{m, responseChannel}
}
