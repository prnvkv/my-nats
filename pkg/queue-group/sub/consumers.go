package sub

import (
	"reflect"
	"runtime"

	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Subscribe receives the message using the subject and Queue group name
func Subscribe(subject string, queueGroupName string) ([]byte, error) {
	serverAddr := viper.GetString("nats.server.addr")
	serverPort := viper.GetString("nats.server.port")
	// subjectName := viper.GetString("nats.subject.dns")

	natsConnection := "nats://" + serverAddr + ":" + serverPort

	log.Infof("Subscriber connecting to nats messaging server: %s with subject %s Queue Group: %s\n", natsConnection, subject, queueGroupName)
	nc, err := nats.Connect(natsConnection)
	if err != nil {
		return nil, err
	}

	// defer nc.Close()

	log.Infof("Consuming the message from the topic: %s\n", subject)
	var receivedMsg []byte
	nc.QueueSubscribe(subject, queueGroupName, func(m *nats.Msg) {
		receivedMsg = m.Data
		err = nc.Publish(m.Reply, []byte("success"))

		if err != nil {
			log.Errorf("Error while ack : %s \n", err.Error())
			return
		}
		log.Printf("[Received] %s\n", receivedMsg)
	})

	if err != nil {
		log.Errorf("Error here!! \n ")
		return nil, err
	}

	if len(receivedMsg) == 0 || reflect.ValueOf(receivedMsg).Kind() == reflect.Ptr && reflect.ValueOf(receivedMsg).IsNil() {
		runtime.Goexit()
	}

	return receivedMsg, nil

}
