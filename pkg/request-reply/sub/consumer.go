package sub

import (
	"reflect"
	"runtime"

	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type CallBackHandlerFunc func(dbString []byte) error

func Subscribe(subject string, cb CallBackHandlerFunc) ([]byte, error) {
	serverAddr := viper.GetString("nats.server.addr")
	serverPort := viper.GetString("nats.server.port")

	natsConnection := "nats://" + serverAddr + ":" + serverPort
	log.Infof("Subscriber connecting to nats messaging server: %s with subject %s\n", natsConnection, subject)

	var err error

	nc, err := nats.Connect(natsConnection)
	if err != nil {
		return nil, err
	}

	defer nc.Close()

	log.Infof("Consuming the message from the topic: %s\n", subject)
	var receivedMsg []byte
	nc.Subscribe(subject, func(m *nats.Msg) {
		receivedMsg = m.Data
		err = cb(m.Data)
		if err != nil {
			nc.Publish(m.Reply, []byte(err.Error()))
			return
		}

		log.Infof("[Received] %s\n", receivedMsg)
		log.Infof("Sending the reply to ack the recieved msg...")
		nc.Publish(m.Reply, []byte("success"))
	})
	if err != nil {
		log.Errorf("Error occurred: %s", err.Error())
		return nil, err
	}

	if len(receivedMsg) == 0 || reflect.ValueOf(receivedMsg).Kind() == reflect.Ptr && reflect.ValueOf(receivedMsg).IsNil() {
		runtime.Goexit()
	}

	return receivedMsg, nil
}
