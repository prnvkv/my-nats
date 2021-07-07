package main

import (
	"runtime"

	"github.com/prnvkv/my-nats/pkg/jsm/sub"
	log "github.com/sirupsen/logrus"
)

func main() {

	log.Infoln("Calling the JSM subscriber...")
	sub.JSMSubscribe()

	runtime.Goexit()
}
