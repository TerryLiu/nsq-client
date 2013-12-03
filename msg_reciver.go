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
	maxInFlight           = 1000
	verbose               = true
	skipEmptyFiles        = false
	tlsEnabled            = false
	tlsInsecureSkipVerify = false
	nsqdTCPAddrs          = util.StringArray{}
	lookupdHTTPAddrs      = util.StringArray{}
	Receiver              *MsgReceiver
)

func init() {
	lookupdHTTPAddrs.Set("106.186.31.48:4161")
	Receiver = &MsgReceiver{
		ExitChan: make(chan int),
	}
}

type MsgReceiver struct {
	usr        User
	msgHandler *MsgHandler
	ExitChan   chan int
}

func (r *MsgReceiver) router(termChan chan os.Signal, hupChan chan os.Signal) {
	for {
		select {
		case <-r.ExitChan:
			//r.Stop()
			return
		case <-termChan:
			//r.Stop()
			return
		case <-hupChan:
			//r.Stop()
			return
		}
	}
}

func (receiver *MsgReceiver) SetLoginUsr(_usr User) error {
	if Receiver == nil {
		Receiver = &MsgReceiver{
			ExitChan: make(chan int),
		}
	}
	Receiver.usr = _usr
	return nil
}

func (receiver *MsgReceiver) AddMsgHandler(topic, channel string, msgHandler *MsgHandler) {
	r, err := nsq.NewReader(topic, channel)
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

	r.AddAsyncHandler(msgHandler)

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
}

func (receiver *MsgReceiver) StartReceiver() {
	hupChan := make(chan os.Signal, 1)
	termChan := make(chan os.Signal, 1)
	signal.Notify(hupChan, syscall.SIGHUP)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	receiver.router(termChan, hupChan)
}

func (r *MsgReceiver) Stop() {
	r.ExitChan <- 1
}
