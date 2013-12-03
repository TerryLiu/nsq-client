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

type SingleChatWindow struct {
	*walk.MainWindow
	chatView   *ChatMsgView
	msgEdit    *walk.TextEdit
	sendBtn    *walk.PushButton
	msgHandler *MsgHandler
	usr        User
	partner    User
}

func NewSingleChatWindow(_usr, _partner User) {
	walk.SetPanicOnError(true)
	myWindow, _ := walk.NewMainWindow()

	mw := &SingleChatWindow{
		MainWindow: myWindow,
		usr:        _usr,
		partner:    _partner,
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

	mw.msgHandler = &MsgHandler{
		topic:   mw.usr.Id,
		channel: MSG_CHANNEL_CHAT,
		msgChan: make(chan *NsqMsg, 1),
	}
	go Receiver.AddMsgHandler(mw.msgHandler)
	go mw.msgRouter()

	mw.MainWindow.Run()
	mw.msgHandler.reader.Stop()
}

func (mw *SingleChatWindow) sendBtn_OnClick() {
	text := mw.msgEdit.Text()
	if strings.EqualFold(text, "") {
		return
	}

	mw.chatView.PostAppendTextln(mw.usr.Nick + " " + time.Now().Format("2006-01-02 15:04:05") + " :")
	mw.chatView.PostAppendTextln("  " + text)
	go Publisher.Write(mw.partner.Id, text)
	mw.msgEdit.SetText("")
}

func (mw *SingleChatWindow) msgRouter() {
	for {
		select {
		case m := <-mw.msgHandler.msgChan:
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
