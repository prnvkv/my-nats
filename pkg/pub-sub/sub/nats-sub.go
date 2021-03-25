package sub

import (
	"reflect"
	"runtime"

	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Subscribe(subject string) ([]byte, error) {
	serverAddr := viper.GetString("nats.server.addr")
	serverPort := viper.GetString("nats.server.port")
	// subjectName := viper.GetString("nats.subject.dns")

	natsConnection := "nats://" + serverAddr + ":" + serverPort

	log.Infof("Subscriber connecting to nats messaging server: %s with subject %s\n", natsConnection, subject)
	nc, err := nats.Connect(natsConnection)
	if err != nil {
		return nil, err
	}

	log.Infof("Consuming the message from the topic: %s\n", subject)
	var receivedMsg []byte
	nc.Subscribe(subject, func(m *nats.Msg) {
		receivedMsg = m.Data
		log.Printf("[Received] %s\n", receivedMsg)
	})

	if len(receivedMsg) == 0 || reflect.ValueOf(receivedMsg).Kind() == reflect.Ptr && reflect.ValueOf(receivedMsg).IsNil() {
		runtime.Goexit()
	}

	return receivedMsg, nil

}
