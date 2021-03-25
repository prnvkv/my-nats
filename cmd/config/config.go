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
)

var (
	// Define flag overrides

	flagNATSServerAddr   = pflag.String("nats.server.addr", defaultNATSServeraddress, "NATS-server address")
	flagNATSServerPort   = pflag.String("nats.server.port", defaultNATSServerPort, "NATS-server port")
	flagNATSSubject_DNS  = pflag.String("nats.subject.dns", defaultNATSSubject_DNS, "Default Subject for NATS")
	flagNATSClientConfig = pflag.String("config.source", defaultConfigSource, "Default config file for NATS client")
)
