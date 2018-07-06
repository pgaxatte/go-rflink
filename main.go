package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/lilorox/go-rflink/rflink"
)

func main() {
	// Set up channel on which to send signal notifications.
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, os.Kill)

	// Parse options
	opts := rflink.ParseOptions()

	// Setup the MQTT publisher
	p, err := rflink.NewPublisher(opts)
	if err != nil {
		log.Fatal(err)
	}
	defer p.Terminate()

	// Setup the sensor reader
	sr, err := rflink.NewSensorReader(opts)
	if err != nil {
		log.Fatal(err)
	}
	defer sr.Close()

	// Start reading/publishing loop
	p.SensorInput = sr
	go p.ReadAndPublish()

	// Wait for interruption
	<-sigc

	// Disconnect the Network Connection.
	if err := p.Disconnect(); err != nil {
		log.Fatal(err)
	}
}
