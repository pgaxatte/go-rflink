package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/lilorox/go-rflink/rflink"
)

var VERSION = "0.0.1"

func main() {
	log.Printf("Starting go-rflink v%s", VERSION)

	// Set up channel on which to send signal notifications.
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, os.Kill)

	// Parse options
	opts := rflink.ParseOptions()
	log.Printf("Options: %v", opts)

	// Setup the MQTT publisher
	p, err := rflink.NewPublisher(opts)
	if err != nil {
		log.Fatal(err)
	}
	defer p.Terminate()
	log.Print("MQTT publisher created")

	// Setup the sensor reader
	sr, err := rflink.NewSensorReader(opts)
	if err != nil {
		log.Fatal(err)
	}
	defer sr.Close()
	log.Print("Sensor reader created")

	// Start reading/publishing loop
	p.SensorInput = sr
	go p.ReadAndPublish()

	// Wait for interruption
	<-sigc
	log.Print("Interruption signal received")

	// Disconnect the Network Connection.
	if err := p.Disconnect(); err != nil {
		log.Fatal(err)
	}
	log.Print("The End.")
}
