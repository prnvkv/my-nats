package main

import (
	"os"
	"runtime"
	"strings"

	"github.com/prnvkv/my-nats/pkg/request-reply/pub"
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
	msg := "Hello world"
	log.Infof("MAIN: message: '%s'", msg)
	subjectName := viper.GetString("nats.subject.dns")

	response, err := pub.Publish(subjectName, msg)
	if err != nil {
		log.Errorf("ERROR: %s", err.Error())
		return
	}

	log.Infoln("Received the message: ", string(response))
	runtime.Goexit()
}
