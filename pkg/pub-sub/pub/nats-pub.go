package pub

import (
	"bytes"
	"encoding/gob"

	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Publish(subject string, message interface{}) error {
	serverAddr := viper.GetString("nats.server.addr")
	serverPort := viper.GetString("nats.server.port")
	// subjectName := viper.GetString("nats.subject.dns")

	var natsConnection = "nats://" + serverAddr + ":" + serverPort

	log.Infof("Connecting the nats server: %s with subject %s\n", natsConnection, subject)
	nc, err := nats.Connect(natsConnection)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err = enc.Encode(message)
	if err != nil {
		return err
	}

	log.Infof("Publishing the message to the subject: '%s'", subject)

	nc.Publish(subject, buf.Bytes())
	return nil
}
