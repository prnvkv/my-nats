package main

import (
	"strings"

	"github.com/prnvkv/my-nats/pkg/pub-sub/pub"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}

func main() {
	msg := "Hello world"
	log.Infof("MAIN: message: '%s'", msg)

	err := pub.Publish(msg)
	if err != nil {
		log.Errorf("ERROR: %s", err.Error())
	}
	return
}
