package main

import (
	"encoding/json"
	"github.com/bitly/go-nsq"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"log"
	"os"
	"strings"
	"time"
)

type GroupChatWindow struct {
	*walk.MainWindow
	usrModel   *UsrModel
	usrList    *walk.ListBox
	chatView   *ChatMsgView
	msgEdit    *walk.TextEdit
	sendBtn    *walk.PushButton
	msgHandler *MsgHandler
	usr        User
}

func NewGroupChatWindow(_usr User) {
	walk.SetPanicOnError(true)
	myWindow, _ := walk.NewMainWindow()

	mw := &GroupChatWindow{
		MainWindow: myWindow,
		usr:        _usr,
		usrModel:   NewUsrModel(),
	}

	mw.SetTitle("简易群聊：" + _usr.Nick)

	usrList, _ := walk.NewListBox(mw)
	mw.usrList = usrList
	mw.usrList.SetModel(mw.usrModel)
	mw.usrList.SetSize(walk.Size{100, 450})
	mw.usrList.SetX(10)
	mw.usrList.SetY(5)
	mw.usrList.ItemActivated().Attach(mw.userlist_ItemActivated)
	mw.usrList.CurrentIndexChanged().Attach(mw.userlist_CurrentIndexChanged)

	msgEdit, _ := walk.NewTextEdit(mw)
	mw.msgEdit = msgEdit
	mw.msgEdit.SetSize(walk.Size{500, 100})
	mw.msgEdit.SetX(120)
	mw.msgEdit.SetY(310)
	mw.msgEdit.SetReadOnly(false)

	chatView, _ := NewChatMsgView(mw)
	mw.chatView = chatView
	mw.chatView.SetSize(walk.Size{500, 300})
	mw.chatView.SetX(120)
	mw.chatView.SetY(5)

	sendBtn, _ := walk.NewPushButton(mw)
	mw.sendBtn = sendBtn
	mw.sendBtn.SetText("发送")
	mw.sendBtn.SetX(560)
	mw.sendBtn.SetY(420)
	mw.sendBtn.SetSize(walk.Size{60, 30})
	mw.sendBtn.Clicked().Attach(mw.sendBtn_OnClick)

	mw.MainWindow.Show()

	mw.msgEdit.SetFocus()
	mw.SetMinMaxSize(walk.Size{645, 500}, walk.Size{645, 500})
	mw.SetSize(walk.Size{645, 500})

	mw.msgHandler = &MsgHandler{
		topic:   "imtech",
		channel: mw.usr.Id,
		msgChan: make(chan *NsqMsg, 1),
	}
	go Receiver.AddMsgHandler(mw.msgHandler)
	go mw.msgRouter()

	mw.MainWindow.Run()
	mw.msgHandler.reader.Stop()
	Publisher.Stop()
	os.Exit(0)

}

func (mw *GroupChatWindow) userlist_CurrentIndexChanged() {
	i := mw.usrList.CurrentIndex()
	item := &mw.usrModel.items[i]
	log.Println("CurrentIndex: ", i)
	log.Println("CurrentName: ", item.Nick)
}

func (mw *GroupChatWindow) userlist_ItemActivated() {
	partner := mw.usrModel.items[mw.usrList.CurrentIndex()]
	//walk.MsgBox(mw, "单聊:"+partner.nick, "单聊功能正在开发中...", walk.MsgBoxIconInformation)
	go NewSingleChatWindow(mw.usr, partner)
}

func (mw *GroupChatWindow) sendBtn_OnClick() {
	text := mw.msgEdit.Text()
	if strings.EqualFold(text, "") {
		return
	}
	mw.chatView.PostAppendTextln(mw.usr.Nick + " " + time.Now().Format("2006-01-02 15:04:05") + " :")
	mw.chatView.PostAppendTextln("  " + text)
	Publisher.Write("imtech", text)
	mw.msgEdit.SetText("")
}

func (mw *GroupChatWindow) msgRouter() {
	for {
		select {
		case m := <-mw.msgHandler.msgChan:
			log.Printf("msgRouter, id = %s, body = %s", string(m.Id[:]), string(m.Body[:]))

			var chatMsg Message
			err := json.Unmarshal(m.Body, &chatMsg)
			if err == nil && !strings.EqualFold(chatMsg.Body.From.Id, mw.usr.Id) {
				mw.chatView.PostAppendTextln(chatMsg.Body.From.Nick + " " + chatMsg.Time + " :")
				mw.chatView.PostAppendTextln("  " + chatMsg.Body.Msg)
			}
			m.returnChannel <- &nsq.FinishedMessage{m.Id, 0, true}
		}
	}
}
