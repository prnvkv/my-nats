package main

import (
	"runtime"

	"github.com/prnvkv/my-nats/pkg/jsm/pub"
	log "github.com/sirupsen/logrus"
)

func main() {

	log.Infoln("Calling the JSM publisher...")
	pub.JSMPublisher()

	runtime.Goexit()
}
