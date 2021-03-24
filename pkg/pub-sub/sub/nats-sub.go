package sub

import (
	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Subscribe() (interface{}, error) {
	serverAddr := viper.GetString("nats.server.addr")
	serverPort := viper.GetString("nats.server.port")
	subjectName := viper.GetString("nats.subject.dns")

	natsConnection := "nats://" + serverAddr + serverPort

	log.Infof("Subscriber connecting to nats messaging server: %s:%s\n", serverAddr, serverPort)
	nc, err := nats.Connect(natsConnection,
		nats.UserInfo("foo", "secret"),
	)
	if err != nil {
		log.Errorf("Error: %s", err)
		return nil, err
	}

	log.Infof("Consuming the message from the topic: %s\n", subjectName)
	var receivedMsg interface{}
	nc.Subscribe(subjectName, func(m *nats.Msg) {
		receivedMsg = m.Data
		log.Printf("[Received] %s\n", string(m.Data))
	})

	return receivedMsg, nil
}
