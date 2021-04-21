package sub

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/nats-io/nats.go"
)

func JetstreamConsumer() {
	closed := make(chan struct{}, 1)
	nc, err := nats.Connect("nats://0.0.0.0:4222",
		nats.ClosedHandler(func(_ *nats.Conn) {
			close(closed)
		}),
		nats.DrainTimeout(5*time.Second))
	noerr(err)
	defer nc.Close()

	js, err := nc.JetStream()
	noerr(err)

	// _, err = js.AddStream(&nats.StreamConfig{
	// 	Name:     "ORDERS",
	// 	Subjects: []string{"ORDERS.received"},
	// })
	// noerr(err)

	// _, err = js.Publish("ORDERS.received", []byte("hello world"))
	// noerr(err)

	sub, err := js.SubscribeSync("ORDERS.received",
		nats.Durable("NEW"),
		// nats.PullMaxWaiting(1)
	)
	fmt.Println("Error here ??????", err)
	noerr(err)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	drain := func() {
		sub.Drain()
		cancel()
	}

	var received int32 = 0
	var processed int32 = 0
	go func() {
		msg, err := sub.NextMsg(2 * time.Second)
		noerr(err)
		atomic.AddInt32(&received, 1)
		time.Sleep(2 * time.Second)
		atomic.AddInt32(&processed, 1)
		fmt.Println("Message is :: ", string(msg.Data))
		err = msg.AckSync()
		noerr(err)
		drain()
	}()

	// Wait until goroutine has finished processing.
	<-ctx.Done()

	// Call drain and close the connection.
	nc.Drain()
	<-closed

	fmt.Printf("Received: %d, Processed: %d\n", received, processed)
}

func noerr(err error) {
	if err != nil {
		panic(err)
	}
}
