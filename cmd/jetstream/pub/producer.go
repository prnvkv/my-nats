package main

import (
	"runtime"

	"github.com/prnvkv/my-nats/pkg/jetstream/pub"
	log "github.com/sirupsen/logrus"
)

func main() {

	log.Infoln("Calling the Jetstream publisher...")
	pub.JetStreamPublish()
	runtime.Goexit()
}
