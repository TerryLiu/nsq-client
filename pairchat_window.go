package main

import (
	"encoding/json"
	"fmt"
	"github.com/bitly/go-nsq"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"log"
	"strings"
	"time"
)

type PairChatWindow struct {
	*walk.MainWindow
	chatView *ChatMsgView
	msgEdit  *walk.TextEdit
	sendBtn  *walk.PushButton
	msgChan  chan *NsqMsg
	usr      User
	partner  User
}

func NewPairChatWindow(_usr, _partner User) {
	walk.SetPanicOnError(true)
	myWindow, _ := walk.NewMainWindow()

	mw := &PairChatWindow{
		MainWindow: myWindow,
		usr:        _usr,
		partner:    _partner,
		msgChan:    make(chan *NsqMsg, 1),
	}

	mw.SetTitle(fmt.Sprintf("与%s私聊", _partner.Nick))

	msgEdit, _ := walk.NewTextEdit(mw)
	mw.msgEdit = msgEdit
	mw.msgEdit.SetSize(walk.Size{530, 100})
	mw.msgEdit.SetX(10)
	mw.msgEdit.SetY(360)
	mw.msgEdit.SetReadOnly(false)

	chatView, _ := NewChatMsgView(mw)
	mw.chatView = chatView
	mw.chatView.SetSize(walk.Size{530, 350})
	mw.chatView.SetX(10)
	mw.chatView.SetY(5)

	sendBtn, _ := walk.NewPushButton(mw)
	mw.sendBtn = sendBtn
	mw.sendBtn.SetText("发送")
	mw.sendBtn.SetX(480)
	mw.sendBtn.SetY(470)
	mw.sendBtn.SetSize(walk.Size{60, 30})
	mw.sendBtn.Clicked().Attach(mw.sendBtn_OnClick)

	mw.MainWindow.Show()

	mw.msgEdit.SetFocus()
	mw.SetMinMaxSize(walk.Size{565, 550}, walk.Size{565, 550})
	mw.SetSize(walk.Size{565, 550})

	pairChatMgr.register(mw.partner.Id, mw.msgChan)
	go mw.msgRouter()

	mw.MainWindow.Run()
	pairChatMgr.unregister(mw.partner.Id)
}

func (mw *PairChatWindow) sendBtn_OnClick() {
	text := mw.msgEdit.Text()
	if strings.EqualFold(text, "") {
		return
	}

	mw.chatView.PostAppendTextln(mw.usr.Nick + " " + time.Now().Format("2006-01-02 15:04:05") + " :")
	mw.chatView.PostAppendTextln("  " + text)
	go Publisher.Write(mw.partner.Id, text)
	mw.msgEdit.SetText("")
}

func (mw *PairChatWindow) msgRouter() {
	for {
		select {
		case m := <-mw.msgChan:
			log.Printf("msgRouter, id = %s, body = %s", string(m.Id[:]), string(m.Body[:]))

			var chatMsg Message
			err := json.Unmarshal(m.Body, &chatMsg)
			if err != nil {
				m.returnChannel <- &nsq.FinishedMessage{m.Id, 0, true}
				continue
			}
			if !strings.EqualFold(chatMsg.Body.From.Id, mw.partner.Id) && !strings.EqualFold(chatMsg.Body.From.Id, mw.usr.Id) {
				continue
			}
			mw.chatView.PostAppendTextln(chatMsg.Body.From.Nick + " " + chatMsg.Time + " :")
			mw.chatView.PostAppendTextln("  " + chatMsg.Body.Msg)
		}
	}
}
