package pub

import (
	// "bytes"
	// "encoding/gob"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/nats-io/nats.go"
	"github.com/prnvkv/my-nats/pkg/util"
)

const (
	//nats.default.svc.cluster.local
	srvAddr = "0.0.0.0"
	srvPort = "4222"
)

func Publish(subject string, message string) ([]byte, error) {
	serverAddr := util.GetEnv("NATS_URL", srvAddr)
	serverPort := util.GetEnv("NATS_PORT", srvPort)

	if len(serverAddr) == 0 {
		serverAddr = srvAddr
	}

	if len(serverPort) == 0 {
		serverPort = srvPort
	}

	var natsConnection = "nats://" + serverAddr + ":" + serverPort

	log.Infof("Connecting the nats server: %s with subject %s\n", natsConnection, subject)
	nc, err := nats.Connect(natsConnection)
	if err != nil {
		return nil, err
	}

	defer nc.Close()

	// var buf bytes.Buffer
	// enc := gob.NewEncoder(&buf)
	// err = enc.Encode(message)
	// if err != nil {
	// 	return nil, err
	// }

	log.Infof("Publishing the message to the subject: '%s'", subject)

	log.Info("Message::", message)

	response, err := nc.Request(subject, []byte(message), 5*time.Second)
	if err != nil {
		return nil, err
	}

	return response.Data, nil
}
