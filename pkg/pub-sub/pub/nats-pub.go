package pub

import (
	"bytes"
	"encoding/gob"

	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Publish(message interface{}) error {
	serverAddr := viper.GetString("nats.server.addr")
	serverPort := viper.GetString("nats.server.port")
	subjectName := viper.GetString("nats.subject.dns")

	log.Infof("Connecting the nats server: %s:%s\n", serverAddr, serverPort)
	var natsConnection = "nats://" + serverAddr + serverPort
	nc, err := nats.Connect(natsConnection,
		nats.UserInfo("foo", "secret"),
	)
	if err != nil {
		log.Errorf("Error: %s", err)
		return err
	}

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err = enc.Encode(message)
	if err != nil {
		return err
	}

	log.Infof("Publishing the message to the subject: '%s'", subjectName)

	nc.Publish(subjectName, buf.Bytes())
	return nil
}
