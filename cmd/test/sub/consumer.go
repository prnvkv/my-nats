package main

import (
	"fmt"
	"runtime"

	"github.com/prnvkv/my-nats/pkg/req-qsub/sub"
	log "github.com/sirupsen/logrus"
)

// func init() {
// 	config.LoadConfig()
// 	viper.BindPFlags(pflag.CommandLine)
// 	viper.AutomaticEnv()
// 	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
// 	// viper.SetConfigName("config")
// 	// viper.SetConfigType("yaml")
// 	// viper.AddConfigPath("../../../deploy/")
// 	viper.AddConfigPath(viper.GetString("config.source"))

// 	if len(viper.GetString("config.file")) != 0 {
// 		err := viper.ReadInConfig()
// 		if err != nil {
// 			log.Errorf("Error cannot read config file : %s\n", err.Error())
// 			os.Exit(1)
// 		}

// 	}

// }

func main() {
	serverAddr := "0.0.0.0"
	serverPort := "4222"
	subj1 := "nats.subject.dns"
	subj2 := "nats.subject.vpn"
	queueGroupName := "test_group"

	// serverAddr := viper.GetString("nats.server.addr")
	// serverPort := viper.GetString("nats.server.port")
	// subjectName := viper.GetString("nats.subject.dns")
	// queueGroupName := viper.GetString("nats.queue_group.name")

	fmt.Println("Server port and subject...", serverAddr, serverPort, subj1, subj2, queueGroupName)
	done := make(chan bool, 2)

	go func() {
		log.Infoln("Func 1:: Sub1:: mesg1 ....")
		msg, err := sub.Subscribe(subj1, queueGroupName, nil, "receieved..")
		if err != nil {
			log.Errorf("Error: %s", err)
			done <- true
			return
		}

		log.Infof("Recieved the message111111111111: %s", msg)

		done <- true

	}()

	go func() {
		log.Infoln("Func 2:: Sub2:: mesg2 ....")
		msg, err := sub.Subscribe(subj2, queueGroupName, nil, "receieved..")
		if err != nil {
			log.Errorf("Error: %s", err)
			done <- true
			return
		}
		log.Infof("Recieved the message2222222222: %s", msg)
		done <- true

	}()

	runtime.Goexit()
}
