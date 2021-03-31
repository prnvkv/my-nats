package main

import (
	"os"
	"runtime"
	"strings"

	"github.com/prnvkv/my-nats/cmd/config"
	"github.com/prnvkv/my-nats/pkg/req-qsub/pub"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

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
