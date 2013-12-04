package main

import (
	"encoding/json"
	"github.com/bitly/go-nsq"
	"log"
)

var (
	pairChatMgr *PairChatManager
)

func init() {
	pairChatMgr = &PairChatManager{
		partnersMap: make(map[string]chan *NsqMsg),
	}
}

type PairChatManager struct {
	ChatMgr
	partnersMap map[string]chan *NsqMsg
}

func (mgr *PairChatManager) SetLoginUsr(_usr User) error {
	if pairChatMgr == nil {
		pairChatMgr = &PairChatManager{
			partnersMap: make(map[string]chan *NsqMsg),
		}
	}
	pairChatMgr.usr = _usr
	pairChatMgr.topic = _usr.Id
	pairChatMgr.channel = MSG_CHANNEL_CHAT
	return nil
}

func (mgr *PairChatManager) register(partner string, msgChan chan *NsqMsg) {
	mgr.partnersMap[partner] = msgChan
	log.Printf("PairChatManager.register[%s]=%v\n", partner, msgChan)
	log.Printf("PairChatManager.partnersMap=%v\n", mgr.partnersMap)
}

func (mgr *PairChatManager) unregister(partner string) {
	log.Printf("PairChatManager.unregister [%s]\n", partner)
	delete(mgr.partnersMap, partner)
}

func (mgr *PairChatManager) GetTopic() string {
	return mgr.topic
}
func (mgr *PairChatManager) GetChannel() string {
	return mgr.channel
}
func (mgr *PairChatManager) SetReader(r *nsq.Reader) {
	mgr.reader = r
}

func (mgr *PairChatManager) HandleMessage(m *nsq.Message, responseChannel chan *nsq.FinishedMessage) {
	log.Printf("PairChatManager...topic=%s\n", mgr.topic)
	log.Printf("PairChatManager...m.Body=%s\n", string(m.Body))

	var msg Message
	err := json.Unmarshal(m.Body, &msg)
	if err != nil {
		responseChannel <- &nsq.FinishedMessage{m.Id, 0, true}
		return
	}
	partner := msg.Body.From.Id
	log.Printf("partner...From=%s\n", msg.Body.From.Id)
	log.Printf("mgr.partnersMap[%s]=%v\n", msg.Body.From.Id, mgr.partnersMap[partner])
	if mgr.partnersMap[partner] != nil {
		mgr.partnersMap[partner] <- &NsqMsg{m, responseChannel}
	}
}

func (mgr *PairChatManager) release() {
	if mgr.reader != nil {
		mgr.reader.Stop()
	}
}
