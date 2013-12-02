package main

import (
	"crypto/tls"
	"flag"
	"github.com/bitly/go-nsq"
	"github.com/bitly/nsq/util"
	_ "github.com/lxn/walk/declarative"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	showVersion           = flag.Bool("version", false, "print version string")
	hostIdentifier        = flag.String("host-identifier", "", "value to output in log filename in place of hostname. <SHORT_HOST> and <HOSTNAME> are valid replacement tokens")
	maxInFlight           = flag.Int("max-in-flight", 1000, "max number of messages to allow in flight")
	verbose               = flag.Bool("verbose", false, "verbose logging")
	skipEmptyFiles        = flag.Bool("skip-empty-files", false, "Skip writting empty files")
	nsqdTCPAddrs          = util.StringArray{}
	lookupdHTTPAddrs      = util.StringArray{}
	tlsEnabled            = flag.Bool("tls", false, "enable TLS")
	tlsInsecureSkipVerify = flag.Bool("tls-insecure-skip-verify", false, "disable TLS server certificate validation")
)

func init() {
	flag.Var(&nsqdTCPAddrs, "nsqd-tcp-address", "nsqd TCP address (may be given multiple times)")
	flag.Var(&lookupdHTTPAddrs, "lookupd-http-address", "lookupd HTTP address (may be given multiple times)")
	lookupdHTTPAddrs.Set("106.186.31.48:4161")
}

type MsgReceiver struct {
	topic    string
	channel  string
	msgChan  chan *Message
	ExitChan chan int
}

type Message struct {
	*nsq.Message
	returnChannel chan *nsq.FinishedMessage
}

func (receiver *MsgReceiver) HandleMessage(m *nsq.Message, responseChannel chan *nsq.FinishedMessage) {
	receiver.msgChan <- &Message{m, responseChannel}
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
		case m := <-receiver.msgChan:
			m.returnChannel <- &nsq.FinishedMessage{m.Id, 0, true}
		}
	}
}

func NewMsgReceiver(_topic, _channel string, _msgChan chan *Message) (*MsgReceiver, error) {
	receiver := &MsgReceiver{
		topic:   _topic,
		channel: _channel,
		msgChan: _msgChan,
	}
	return receiver, nil
}

func StartReceiver(topic, channel string, msgChan chan *Message) {
	receiver, err := NewMsgReceiver(topic, channel, msgChan)

	if err != nil {
		log.Fatal(err.Error())
	}

	if topic == "" || channel == "" {
		log.Fatalf("--topic and --channel are required")
	}

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
	r.SetMaxInFlight(*maxInFlight)
	r.VerboseLogging = *verbose

	if *tlsEnabled {
		r.TLSv1 = true
		r.TLSConfig = &tls.Config{
			InsecureSkipVerify: *tlsInsecureSkipVerify,
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
