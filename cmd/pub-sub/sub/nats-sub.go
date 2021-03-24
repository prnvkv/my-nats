package main

import (
	"runtime"
	"strings"

	"github.com/nats-io/nats.go"
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
	nc, err := nats.Connect("nats://0.0.0.0:4222",
		nats.UserInfo("foo", "secret"),
	)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	nc.Subscribe("greeting", func(m *nats.Msg) {
		log.Printf("[Received] %s", string(m.Data))
	})

	runtime.Goexit()
}
