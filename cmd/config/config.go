package config

import (
	"github.com/spf13/pflag"
)

const (
	//Configuration defaults for the demo and dev environment. Update these values in helm charts during deployment.

	// NATS-server related
	// Default NATS server address
	// docker : "0.0.0.0"
	// k3s : nats.<namespace_name>.svc.cluster.local
	defaultNATSServeraddress = "0.0.0.0"

	// Default NATS-server port
	defaultNATSServerPort = "4222"

	// Default NATS-server subject
	// NOTE: This is not required. NATS-server has the capability to create the subject in runtime.
	// But this configuration is necessary for the documentation and consistency of the service.
	defaultNATSSubject_DNS = "dns.service"

	// Configuration file path for local development
	defaultConfigSource = "../../../deploy/"

	// NATS Queue group name
	// Note: Subscribers with same name form a queue group
	defaultNATSQueueGroupName = "dns.subscribers"

	// // NATS Default config file for local
	defaultNATSConfigFile = ""
)

var (
	// Define flag overrides

	// Server related details
	flagNATSServerAddr = pflag.String("nats.server.addr", defaultNATSServeraddress, "NATS-server address")
	flagNATSServerPort = pflag.String("nats.server.port", defaultNATSServerPort, "NATS-server port")

	// Pub sub related details
	flagNATSSubject_DNS = pflag.String("nats.subject.dns", defaultNATSSubject_DNS, "Default Subject for NATS")

	// Config file details
	flagNATSClientConfig     = pflag.String("config.source", defaultConfigSource, "Default config file for NATS client")
	flagNATSClientConfigFile = pflag.String("config.file", defaultNATSConfigFile, "Default config file ")

	// Queue group related details
	flagNATSQueueGroup = pflag.String("nats.queue_group.name", defaultNATSQueueGroupName, "NATS Queue Group name")
)

func LoadConfig() {
	pflag.Parse()
}
