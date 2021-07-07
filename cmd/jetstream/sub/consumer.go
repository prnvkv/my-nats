package main

import (
	"runtime"

	"github.com/prnvkv/my-nats/pkg/jetstream/sub"
	log "github.com/sirupsen/logrus"
)

func main() {

	log.Infoln("Calling the Jetstream subscriber...")
	sub.JetstreamConsumer()

	runtime.Goexit()
}
