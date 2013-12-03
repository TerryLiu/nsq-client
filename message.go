package main

import ()

/*
消息格式：
{
	"Topic":"xxx",
	"Type":1,
	"Time":"2013-12-04 00:00:00"
	"body":{
		"From":{
			"Id":"xxx",
			"Nick":"xxx",
		},
		"Msg":"xxxxxxx"
	}
}
*/
const MSG_TYPE_SYSTEM = 0
const MSG_TYPE_CHAT = 1

const MSG_TOPIC_SYSTEM = "__SYSTEM__"
const MSG_CHANNEL_CHAT = "__CHAT__"

type Message struct {
	Topic string
	Type  int //消息类型
	Time  string
	Body  MessageBody
}

type MessageBody struct {
	From User
	Msg  string
}
