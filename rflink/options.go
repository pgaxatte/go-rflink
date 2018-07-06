package rflink

import (
	"net/url"
)

type Options struct {
	// MQTT options
	PublishURL      *url.URL
	PublishClientID string
	PublishTopic    string

	// Serial connection options
	SerialDevice string
	SerialBaud   int
}

func ParseOptions() *Options {
	return &Options{
		PublishURL: &url.URL{
			Scheme: "tcp",
			Host:   "10.1.0.4:1883",
		},
		PublishClientID: "rflink",
		PublishTopic:    "rflink",

		SerialDevice: "/dev/ttyUSB0",
		SerialBaud:   57600,
	}
}
