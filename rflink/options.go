package rflink

import (
	"log"

	"github.com/vrischmann/envconfig"
)

// Options stores the options needed to communicate with RFLink and the
// message queue
type Options struct {
	// MQTT options
	Publish struct {
		Host     string `envconfig:"default=localhost:1883"` // host:port
		Scheme   string `envconfig:"default=tcp"`
		ClientID string `envconfig:"default=rflink"`
		Topic    string `envconfig:"default=rflink"`
	}

	// Serial connection options
	Serial struct {
		Device string `envconfig:"default=/dev/ttyUSB0"`
		Baud   int    `envconfig:"default=57600"`
	}
}

// GetOptions reads the options from the environment and returns an Options
// struct
func GetOptions() *Options {
	var opts Options
	if err := envconfig.Init(&opts); err != nil {
		log.Fatal("Could not parse options:", err)
	}
	return &opts
}
