package main

import (
	"runtime"

	"github.com/nats-io/nats.go"

	"github.com/prnvkv/my-nats/pkg/jsm-queue/pub"
	log "github.com/sirupsen/logrus"
)

const (
	stream  = "dns"
	subject = "dns.query"
)

func main() {

	log.Infoln("Calling the Jetstream subscriber...")

	url := "nats://localhost:4222"
	nats, err := nats.Connect(
		url,
	)
	if err != nil {
		log.Fatalln(err)
	}
	jsClient, err := nats.JetStream()
	if err != nil {
		log.Fatalln(err)
	}

	queue, err := pub.NewQueue(stream, subject, jsClient)
	noerr(err)

	ch := make(chan *pub.Message)
	err = queue.Subscribe(ch)
	noerr(err)
	data := <-ch

	log.Println(string(data.Data))

	// Begin to unsubscribe here..
	log.Println("Beginning to unsubscribe... ")
	err = queue.Unsubscribe()
	noerr(err)

	runtime.Goexit()
}

func noerr(err error) {
	if err != nil {
		panic(err)
	}
}
