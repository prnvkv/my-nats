package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/prnvkv/my-nats/cmd/config"
	"github.com/prnvkv/my-nats/pkg/pub-sub/sub"
	"github.com/prnvkv/my-nats/pkg/util"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// func init() {
// 	pflag.Parse()
// 	viper.BindPFlags(pflag.CommandLine)
// 	viper.SetConfigName("config")
// 	viper.SetConfigType("yaml")
// 	viper.AddConfigPath("../../../deploy/")
// 	viper.AddConfigPath(viper.GetString("config.source"))
// 	viper.AutomaticEnv()
// 	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

// 	err := viper.ReadInConfig()
// 	if err != nil {
// 		log.Errorf("Error cannot read config file : %s\n", err.Error())
// 		os.Exit(1)
// 	}

// }
func init() {
	config.LoadConfig()
	viper.BindPFlags(pflag.CommandLine)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	// viper.SetConfigName("config")
	// viper.SetConfigType("yaml")
	// viper.AddConfigPath("../../../deploy/")
	viper.AddConfigPath(viper.GetString("config.source"))

	if len(viper.GetString("config.file")) != 0 {
		err := viper.ReadInConfig()
		if err != nil {
			log.Errorf("Error cannot read config file : %s\n", err.Error())
			os.Exit(1)
		}

	}

}

const (
	defaultSubject = "test_subject"
	//Server address, in cluster it will be like: nats.default.svc.cluster.local
	srvAddr = "0.0.0.0"
	// Server Port used by the nats by default
	srvPort = "4222"
)

func main() {

	serverAddr := util.GetEnv("NATS_URL", srvAddr)
	serverPort := util.GetEnv("NATS_PORT", srvPort)

	if len(serverAddr) == 0 {
		serverAddr = srvAddr
	}

	if len(serverPort) == 0 {
		serverPort = srvPort
	}

	subjectName := util.GetEnv("NATS_SUBJECT", defaultSubject)

	fmt.Println("Server port and subject...", serverAddr, serverPort, subjectName)
	for {
		msg, err := sub.Subscribe(subjectName)
		if err != nil {
			log.Errorf("Error: %s", err)
			return
		}
		log.Infof("Recieved the message: %s", msg)

	}

	//runtime.Goexit()
}
