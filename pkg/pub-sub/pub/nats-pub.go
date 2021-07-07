package pub

import (
	"bytes"
	"encoding/gob"
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

// Publish publishes the message using subject name
func Publish(subject string, message interface{}) error {
	serverAddr := util.GetEnv("NATS_URL", srvAddr)
	serverPort := util.GetEnv("NATS_PORT", srvPort)

	if len(serverAddr) == 0 {
		serverAddr = srvAddr
	}

	if len(serverPort) == 0 {
		serverPort = srvPort
	}

	var natsConnection = "nats://" + serverAddr + ":" + serverPort
	// Connect Options.
	opts := []nats.Option{nats.Name("NATS Sample Publisher")}
	opts = setupConnOptions(opts)

	log.Infof("Connecting the nats server: %s with subject %s\n", natsConnection, subject)
	nc, err := nats.Connect(natsConnection, opts...)
	if err != nil {
		return err
	}

	defer nc.Flush()

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err = enc.Encode(message)
	if err != nil {
		return err
	}

	log.Infof("Publishing the message to the subject: '%s'", subject)

	err = nc.Publish(subject, buf.Bytes())
	if err != nil {
		log.Errorln("Error while publishing:: ", err)
		return err
	}
	log.Infoln("Published :: ", message)
	return nil
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
