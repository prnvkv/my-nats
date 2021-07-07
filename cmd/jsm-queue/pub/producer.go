package main

import (
	"github.com/nats-io/nats.go"
	"github.com/prnvkv/my-nats/pkg/jsm-queue/pub"
	log "github.com/sirupsen/logrus"
)

const (
	stream  = "dns"
	subject = "dns.query"
)

func main() {

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

	log.Infoln("Calling the Jetstream publisher...")

	queue, err := pub.NewQueue(stream, subject, jsClient)
	noerr(err)
	m := pub.Message{Data: []byte("Hello world")}
	log.Println("Message struct is :: ", m)
	log.Println("Message data is :: ", string(m.Data))
	err = queue.Publish(&m)
	noerr(err)

}

func noerr(err error) {
	if err != nil {
		panic(err)
	}
}
