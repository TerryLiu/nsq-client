package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/bitly/go-nsq"
	"github.com/bitly/nsq/util"
	_ "github.com/lxn/walk/declarative"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	showVersion      = flag.Bool("version", false, "print version string")
	hostIdentifier   = flag.String("host-identifier", "", "value to output in log filename in place of hostname. <SHORT_HOST> and <HOSTNAME> are valid replacement tokens")
	topic            = flag.String("topic", "", "nsq topic")
	channel          = flag.String("channel", "chat", "nsq channel")
	maxInFlight      = flag.Int("max-in-flight", 1000, "max number of messages to allow in flight")
	verbose          = flag.Bool("verbose", false, "verbose logging")
	skipEmptyFiles   = flag.Bool("skip-empty-files", false, "Skip writting empty files")
	nsqdTCPAddrs     = util.StringArray{}
	lookupdHTTPAddrs = util.StringArray{}

	tlsEnabled            = flag.Bool("tls", false, "enable TLS")
	tlsInsecureSkipVerify = flag.Bool("tls-insecure-skip-verify", false, "disable TLS server certificate validation")
)

func init() {
	flag.Var(&nsqdTCPAddrs, "nsqd-tcp-address", "nsqd TCP address (may be given multiple times)")
	flag.Var(&lookupdHTTPAddrs, "lookupd-http-address", "lookupd HTTP address (may be given multiple times)")
}

type FileLogger struct {
	logChan  chan *Message
	ExitChan chan int
}

type Message struct {
	*nsq.Message
	returnChannel chan *nsq.FinishedMessage
}

type SyncMsg struct {
	m             *nsq.FinishedMessage
	returnChannel chan *nsq.FinishedMessage
}

func (l *FileLogger) HandleMessage(m *nsq.Message, responseChannel chan *nsq.FinishedMessage) {
	id := string(m.Id[:])
	body := string(m.Body[:])
	log.Printf("HandleMessage, id = %s, body = %s", id, body)
	l.logChan <- &Message{m, responseChannel}
}

func (f *FileLogger) router(r *nsq.Reader, termChan chan os.Signal, hupChan chan os.Signal) {
	ticker := time.NewTicker(time.Duration(30) * time.Second)
	exit := false
	for {
		select {
		case <-r.ExitChan:
			exit = true
		case <-termChan:
			ticker.Stop()
			r.Stop()
		case <-hupChan:

		case <-ticker.C:
		case m := <-f.logChan:
			m.returnChannel <- &nsq.FinishedMessage{m.Id, 0, true}
		}
		if exit {
			close(f.ExitChan)
			break
		}
	}
}

func NewFileLogger() (*FileLogger, error) {
	f := &FileLogger{
		logChan:  make(chan *Message, 1),
		ExitChan: make(chan int),
	}
	return f, nil
}

func StartReceiver() {
	flag.Parse()

	if *showVersion {
		fmt.Printf("nsq_to_file v%s\n", util.BINARY_VERSION)
		return
	}

	if *topic == "" || *channel == "" {
		log.Fatalf("--topic and --channel are required")
	}

	if *maxInFlight <= 0 {
		log.Fatalf("--max-in-flight must be > 0")
	}

	if len(nsqdTCPAddrs) == 0 && len(lookupdHTTPAddrs) == 0 {
		log.Fatalf("--nsqd-tcp-address or --lookupd-http-address required.")
	}
	if len(nsqdTCPAddrs) != 0 && len(lookupdHTTPAddrs) != 0 {
		log.Fatalf("use --nsqd-tcp-address or --lookupd-http-address not both")
	}

	hupChan := make(chan os.Signal, 1)
	termChan := make(chan os.Signal, 1)
	signal.Notify(hupChan, syscall.SIGHUP)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)

	f, err := NewFileLogger()
	if err != nil {
		log.Fatal(err.Error())
	}

	r, err := nsq.NewReader(*topic, *channel)
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

	r.AddAsyncHandler(f)
	go f.router(r, termChan, hupChan)

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
	<-f.ExitChan
}
