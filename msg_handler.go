package main

type MsgHandler struct {
	topic    string
	channel  string
	usr      User
	msgChan  chan *ReciveMsg
	ExitChan chan int
}
