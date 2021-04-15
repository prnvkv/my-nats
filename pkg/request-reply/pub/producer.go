package pub

import (
	"bytes"
	"encoding/gob"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/nats-io/nats.go"
	"github.com/spf13/viper"
)

// Publish used to publish the messages using subject name and message itself.
func Publish(subject string, message interface{}) ([]byte, error) {
	serverAddr := viper.GetString("nats.server.addr")
	serverPort := viper.GetString("nats.server.port")
	var natsConnection = "nats://" + serverAddr + ":" + serverPort

	log.Infof("Connecting the nats server: %s with subject %s\n", natsConnection, subject)
	nc, err := nats.Connect(natsConnection)
	if err != nil {
		return nil, err
	}

	defer nc.Close()

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err = enc.Encode(message)
	if err != nil {
		return nil, err
	}

	log.Infof("Publishing the message to the subject: '%s'", subject)

	response, err := nc.Request(subject, buf.Bytes(), 1*time.Second)
	if err != nil {
		return nil, err
	}

	return response.Data, nil
}
