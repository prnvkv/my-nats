package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/prnvkv/my-nats/pkg/request-reply/sub"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../../../deploy/")
	viper.AddConfigPath(viper.GetString("config.source"))
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := viper.ReadInConfig()
	if err != nil {
		log.Errorf("Error cannot read config file : %s\n", err.Error())
		os.Exit(1)
	}

}

func main() {

	serverAddr := viper.GetString("nats.server.addr")
	serverPort := viper.GetString("nats.server.port")
	subjectName := viper.GetString("nats.subject.dns")

	fmt.Println("Server port and subject...", serverAddr, serverPort, subjectName)
	msg, err := sub.Subscribe(subjectName)
	if err != nil {
		log.Errorf("Error: %s", err)
		return
	}
	log.Infof("Recieved the message: %s", msg)

	runtime.Goexit()
}
