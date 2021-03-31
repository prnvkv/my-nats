package sub

import (
	"reflect"
	"runtime"

	"os"

	"github.com/nats-io/nats.go"
	"github.com/prnvkv/my-nats/pkg/util"
	log "github.com/sirupsen/logrus"
)

const (
	// nats.default.svc.cluster.local
	srvAddr = "0.0.0.0"
	srvPort = "4222"
)

type callBackFunc func(msg []byte) error

func Subscribe(subject string, queueGroupName string, cb callBackFunc, ackMsg string) ([]byte, error) {
	serverAddr := util.GetEnv("NATS_URL", srvAddr)
	serverPort := util.GetEnv("NATS_PORT", srvPort)

	if len(serverAddr) == 0 {
		serverAddr = srvAddr
	}

	if len(serverPort) == 0 {
		serverPort = srvPort
	}

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
		log.Printf("[Received] %s \n Calling the function handler ... \n", receivedMsg)
		if cb != nil {

			err = cb(m.Data)
			if err != nil {
				return
			}

		}
		log.Printf("[Received] %s\n", receivedMsg)
		log.Infof("Sending the ack: %s \n", ackMsg)
		err = nc.Publish(m.Reply, []byte(ackMsg))

		if err != nil {
			log.Errorf("Error while ack : %s \n", err.Error())
			return
		}
		log.Infof("ACK Sent successfully")

	})

	if err != nil {
		log.Errorf("ERROR: %s \n ", err.Error())
		return nil, err
	}

	if len(receivedMsg) == 0 || reflect.ValueOf(receivedMsg).Kind() == reflect.Ptr && reflect.ValueOf(receivedMsg).IsNil() {
		runtime.Goexit()
	}

	return receivedMsg, nil

}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
