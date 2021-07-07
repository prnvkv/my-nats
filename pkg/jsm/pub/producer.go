package pub

import (
	"time"

	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
)

// Name of the stream
const stream = "EXAMPLE"

// setupConnOptions sets the connection options
func setupConnOptions(opts []nats.Option) []nats.Option {
	totalWait := 10 * time.Second
	reconnectDelay := time.Second

	opts = append(opts, nats.ReconnectWait(reconnectDelay))
	opts = append(opts, nats.MaxReconnects(int(totalWait/reconnectDelay)))
	opts = append(opts, nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
		log.Printf("Disconnected due to: %s, will attempt reconnects for %.0fs", err, totalWait.Minutes())
	}))
	opts = append(opts, nats.ReconnectHandler(func(nc *nats.Conn) {
		log.Printf("Reconnected [%s]", nc.ConnectedUrl())
	}))
	opts = append(opts, nats.ClosedHandler(func(nc *nats.Conn) {
		log.Fatalf("Exiting: %v", nc.LastError())
	}))
	return opts
}

func JSMPublisher() {
	opts := []nats.Option{nats.Name("NATS JetStream Transfer")}
	opts = setupConnOptions(opts)

	urls := "nats://0.0.0.0:4222"

	log.Infof("Attempting connection to nats server at %q", urls)

	// Connect to NATS
	nc, err := nats.Connect(urls, opts...)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	// Create our jetstream context.
	// On an error we will just exit.

	// errHandler := func(_ nats.JetStream, _ *nats.Msg, err error) {
	// 	log.Fatalf("Error sending chunk to JetStream: %v", err)
	// }

	// We will use a sliding window and async publishes to maximize performance.
	// const maxPending = 8 // 8 * 64k
	js, err := nc.JetStream(
	// nats.PublishAsyncMaxPending(maxPending),
	// nats.PublishAsyncErrHandler(errHandler),
	)
	if err != nil {
		log.Fatalf("%v", err)
	}
	// if _, err := js.StreamInfo(stream); err != nil {
	// 	log.Fatalf("Stream %q already exists", stream)
	// }

	// Delivery subject as an inbox to avoid accidentally interfering with other subjects
	subj := nats.NewInbox()
	log.Infoln("Subject is %s", subj)

	// Create the stream
	_, err = js.AddStream(&nats.StreamConfig{
		Name:     stream,
		Subjects: []string{subj},
	})
	if err != nil {
		log.Fatalf("Unexpected error creating the stream: %q", err.Error())
	}

	if _, err := js.PublishAsync(subj, []byte(`Hello people, how are you doing?`)); err != nil {
		log.Fatalf("Error while publishing the message: %q", err.Error())
	}

	log.Printf("Completed... ")

}
