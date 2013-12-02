package main

import (
	"crypto/tls"
	"github.com/bitly/go-nsq"
	"github.com/bitly/nsq/util"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	maxInFlight           int
	verbose               bool
	skipEmptyFiles        bool
	tlsEnabled            bool
	tlsInsecureSkipVerify bool
	nsqdTCPAddrs          = util.StringArray{}
	lookupdHTTPAddrs      = util.StringArray{}
)

func init() {
	maxInFlight = 1000
	verbose = true
	skipEmptyFiles = false
	tlsEnabled = false
	tlsInsecureSkipVerify = false
	lookupdHTTPAddrs.Set("106.186.31.48:4161")
}

type MsgReceiver struct {
	topic    string
	channel  string
	usr      User
	msgChan  chan *ReciveMsg
	ExitChan chan int
}

type ReciveMsg struct {
	*nsq.Message
	returnChannel chan *nsq.FinishedMessage
}

func (receiver *MsgReceiver) HandleMessage(m *nsq.Message, responseChannel chan *nsq.FinishedMessage) {
	receiver.msgChan <- &ReciveMsg{m, responseChannel}
}

func (receiver *MsgReceiver) router(r *nsq.Reader, termChan chan os.Signal, hupChan chan os.Signal) {
	for {
		select {
		case <-r.ExitChan:
			r.Stop()
			return
		case <-termChan:
			r.Stop()
			return
		case <-hupChan:
			r.Stop()
			return
		}
	}
}

func NewMsgReceiver(_topic string, _usr User, _msgChan chan *ReciveMsg) (*MsgReceiver, error) {
	if _topic == "" || _usr.Id == "" {
		log.Fatalf("topic and channel are required")
	}
	receiver := &MsgReceiver{
		topic:    _topic,
		usr:      _usr,
		channel:  _usr.Id,
		msgChan:  _msgChan,
		ExitChan: make(chan int),
	}
	return receiver, nil
}

func (receiver *MsgReceiver) StartReceiver() {
	if len(nsqdTCPAddrs) == 0 && len(lookupdHTTPAddrs) == 0 {
		log.Fatalf("--nsqd-tcp-address or --lookupd-http-address required.")
	}
	if len(nsqdTCPAddrs) != 0 && len(lookupdHTTPAddrs) != 0 {
		log.Fatalf("use --nsqd-tcp-address or --lookupd-http-address not both")
	}

	r, err := nsq.NewReader(receiver.topic, receiver.channel)
	if err != nil {
		log.Fatalf(err.Error())
	}
	r.SetMaxInFlight(maxInFlight)
	r.VerboseLogging = verbose

	if tlsEnabled {
		r.TLSv1 = true
		r.TLSConfig = &tls.Config{
			InsecureSkipVerify: tlsInsecureSkipVerify,
		}
	}

	r.AddAsyncHandler(receiver)

	for _, addrString := range nsqdTCPAddrs {
		err := r.ConnectToNSQ(addrString)
		if err != nil {
			log.Fatalf(err.Error())
		}
	}

	for _, addrString := range lookupdHTTPAddrs {
		log.Printf("lookupd addr %s", addrString)
		err := r.ConnectToLookupd(addrString)
		if err != nil {
			log.Fatalf(err.Error())
		}
	}

	hupChan := make(chan os.Signal, 1)
	termChan := make(chan os.Signal, 1)
	signal.Notify(hupChan, syscall.SIGHUP)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	receiver.router(r, termChan, hupChan)
}

func (r *MsgReceiver) Stop() {
	r.ExitChan <- 1
}
