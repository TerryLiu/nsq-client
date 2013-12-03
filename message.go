package main

import ()

/*
消息格式：
{
	"Topic":"xxx",
	"Type":1,
	"body":{
		"From":{
			"Id":"xxx",
			"Nick":"xxx",
		},
		"Msg":"xxxxxxx"
	}
}
*/

const MSG_TYPE_CHAT = 1

type Message struct {
	Topic string
	Type  int //消息类型
	Body  MessageBody
}

type MessageBody struct {
	From User
	Msg  string
}
