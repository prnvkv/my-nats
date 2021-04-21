package main

import (
	"runtime"

	"github.com/prnvkv/my-nats/pkg/req-qsub/pub"
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
	msg1 := "Hello world first message 1111111111111111"
	log.Infof("MAIN: message111111111111: '%s'", msg1)
	sub1 := "nats.subject.dns" // api/v1/dns
	sub2 := "nats.subject.vpn" // api/v1/vpn
	// subjectName := viper.GetString("nats.subject.dns")

	done := make(chan bool, 2)

	go func() {
		response, err := pub.Publish(sub1, msg1)
		if err != nil {
			log.Errorf("ERROR: %s", err.Error())
			done <- true
			return
		}
		log.Infoln("Received the message111111111111: ", string(response))
		done <- true

	}()

	msg2 := "Hello world again 2222222222222222"
	log.Infof("MAIN: message2222222222222222: '%s'", msg2)
	go func() {
		response, err := pub.Publish(sub2, msg2)
		if err != nil {
			log.Errorf("ERROR: %s", err.Error())
			done <- true
			return
		}
		log.Infoln("Received the message2222222222: ", string(response))
		done <- true
	}()

	<-done
	<-done
	log.Infoln("Exitting")
	runtime.Goexit()
}
