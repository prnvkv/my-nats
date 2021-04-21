package sub

import (
	"fmt"
	"math"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/prnvkv/my-nats/pkg/util"
	log "github.com/sirupsen/logrus"
)

// Stream name
const stream = "EXAMPLE"

// type callbackHandler func() error
// var cb

func JSMSubscribe() {
	opts := []nats.Option{nats.Name("NATS JetStream Transfer")}
	opts = util.SetupConnOptions(opts)

	urls := "nats://0.0.0.0:4222"
	log.Infof("Attempting connection to nats server at %q", urls)
	// Connect to NATS
	nc, err := nats.Connect(urls, opts...)
	if err != nil {
		log.Fatal("Error while connecting... ", err)
	}
	defer nc.Close()

	js, err := nc.JetStream()
	if err != nil {
		log.Fatalf("Jetstream : %v", err)
	}

	si, err := js.StreamInfo(stream)
	if err != nil {
		log.Fatalf("Could not find stream: %s", stream)
	}

	createSub := func(startSeq uint64) *nats.Subscription {
		sub, err := js.SubscribeSync(
			si.Config.Subjects[0],
			// stream,

			nats.AckNone(),
			nats.MaxDeliver(1),
			nats.StartSequence(startSeq),
			nats.EnableFlowControl(),
		)
		if err != nil {
			log.Fatalf("Error creating consumer: %v", err)
		}
		return sub
	}

	sub := createSub(1)
	defer sub.Unsubscribe()

	start := time.Now()
	bytes, last, eseq := 0, si.State.Msgs, uint64(1)

	// Loop over our inbound messages
	for m, err := sub.NextMsg(5 * time.Second); err == nil; m, err = sub.NextMsg(time.Second) {
		meta, err := m.Metadata()
		if err != nil {
			log.Fatalf(err.Error())
		}

		if eseq != meta.Sequence.Stream {
			log.Printf("Missed chunk sequence, expected %d got %d, resetting ", eseq, meta.Sequence.Stream)
			sub = createSub(eseq)
		}
		eseq++
		if eseq > last {
			break
		}
	}

	log.Printf("Completed retrieval of %v in %v", friendlyBytes(bytes), time.Since(start))

}

func friendlyBytes(bytes int) string {
	fbytes := float64(bytes)
	base := 1024
	pre := []string{"K", "M", "G", "T", "P", "E"}
	if fbytes < float64(base) {
		return fmt.Sprintf("%v B", fbytes)
	}
	exp := int(math.Log(fbytes) / math.Log(float64(base)))
	index := exp - 1
	return fmt.Sprintf("%.2f %sB", fbytes/math.Pow(float64(base), float64(exp)), pre[index])
}
