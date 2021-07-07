package sub

import (
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/prnvkv/my-nats/pkg/util"
	log "github.com/sirupsen/logrus"
)

const (
	//Server address, in cluster it will be like: nats.default.svc.cluster.local
	srvAddr = "0.0.0.0"
	// Server Port used by the nats by default
	srvPort = "4222"
)

// Subscribe receives the message using the subject name
func Subscribe(subject string) ([]byte, error) {
	serverAddr := util.GetEnv("NATS_URL", srvAddr)
	serverPort := util.GetEnv("NATS_PORT", srvPort)

	if len(serverAddr) == 0 {
		serverAddr = srvAddr
	}

	if len(serverPort) == 0 {
		serverPort = srvPort
	}

	natsConnection := "nats://" + serverAddr + ":" + serverPort
	// Connect Options.
	opts := []nats.Option{nats.Name("NATS Sample Subscriber")}
	opts = setupConnOptions(opts)

	log.Infof("Subscriber connecting to nats messaging server: %s with subject %s\n", natsConnection, subject)
	nc, err := nats.Connect(natsConnection, opts...)
	if err != nil {
		return nil, err
	}

	defer nc.Flush()
	wg := sync.WaitGroup{}
	wg.Add(1)

	log.Infof("Consuming the message from the topic: %s\n", subject)
	var receivedMsg []byte
	nc.Subscribe(subject, func(m *nats.Msg) {
		defer wg.Done()
		receivedMsg = m.Data
		log.Printf("[Received] %s\n", receivedMsg)
	})

	wg.Wait()

	// if len(receivedMsg) == 0 || reflect.ValueOf(receivedMsg).Kind() == reflect.Ptr && reflect.ValueOf(receivedMsg).IsNil() {
	// 	log.Println("No message received")
	// 	runtime.Goexit()
	// }
	// if len(receivedMsg) == 0 {
	// 	log.Println("No message received")
	// 	runtime.Goexit()
	// }

	return receivedMsg, nil

}

func setupConnOptions(opts []nats.Option) []nats.Option {
	totalWait := 10 * time.Minute
	reconnectDelay := time.Second

	opts = append(opts, nats.ReconnectWait(reconnectDelay))
	opts = append(opts, nats.MaxReconnects(int(totalWait/reconnectDelay)))
	opts = append(opts, nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
		log.Printf("Disconnected due to:%s, will attempt reconnects for %.0fm", err, totalWait.Minutes())
	}))
	opts = append(opts, nats.ReconnectHandler(func(nc *nats.Conn) {
		log.Printf("Reconnected [%s]", nc.ConnectedUrl())
	}))
	opts = append(opts, nats.ClosedHandler(func(nc *nats.Conn) {
		log.Fatalf("Exiting: %v", nc.LastError())
	}))
	return opts
}
