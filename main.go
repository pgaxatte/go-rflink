package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"

	"github.com/tarm/serial"
	"github.com/yosssi/gmq/mqtt/client"
)

type SensorData struct {
	Model       string   `json:"model"`
	Id          string   `json:"id"`
	Temperature *float32 `json:"t,omitempty"`
	Humidity    *uint16  `json:"h,omitempty"`
}

func (sd *SensorData) String() string {
	format := "%s [%s]:"
	args := []interface{}{
		sd.Model,
		sd.Id,
	}

	if sd.Temperature != nil {
		format += " temp=%.1f°C"
		args = append(args, *sd.Temperature)
	}

	if sd.Humidity != nil {
		format += " hum=%d%%"
		args = append(args, *sd.Humidity)
	}

	return fmt.Sprintf(format, args...)
}

func (sd *SensorData) Publish(cli *client.Client) (err error) {
	b, err := json.Marshal(sd)
	if err != nil {
		return err
	}
	log.Print(string(b))

	err = cli.Publish(&client.PublishOptions{
		TopicName: []byte("rflink"),
		Message:   b,
	})
	if err != nil {
		return err
	}

	return nil
}

func strToUint16(s string, base int) (uint16, error) {
	u, err := strconv.ParseUint(s, base, 16)
	if err != nil {
		return 0, err
	}
	return uint16(u), nil
}

func parse(msg string) (*SensorData, error) {
	pieces := strings.Split(msg, ";")

	sd := SensorData{
		Model: strings.Replace(pieces[2], " ", "_", -1),
	}
	for i := 3; i < len(pieces); i++ {
		arr := strings.SplitN(pieces[i], "=", 2)
		switch arr[0] {
		case "ID":
			sd.Id = arr[1]
		case "TEMP":
			t, err := strToUint16(arr[1], 16)
			if err != nil {
				return nil, errors.New("Skipping message, temperature could not be parsed")
			}
			temp := float32(t) / 10
			sd.Temperature = &temp
		case "HUM":
			h, err := strToUint16(arr[1], 10)
			if err != nil {
				return nil, errors.New("Skipping message, humidity could not be parsed")
			}
			sd.Humidity = &h
		}
	}
	return &sd, nil
}

func read(r *bufio.Reader, cli *client.Client) {
	for {
		line, _, err := r.ReadLine()
		if err != nil {
			log.Fatal(err)
		}

		sd, err := parse(string(line))
		if err != nil {
			log.Printf("Error parsing \"%s\": %s", line, err)
		}

		err = sd.Publish(cli)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	// Set up channel on which to send signal notifications.
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, os.Kill)

	/*
	 * Setup the MQTT client
	 */
	cli := client.New(&client.Options{
		ErrorHandler: func(err error) {
			log.Fatal(err)
		},
	})
	defer cli.Terminate()

	err := cli.Connect(&client.ConnectOptions{
		Network:  "tcp",
		Address:  "10.1.0.4:1883",
		ClientID: []byte("rflink"),
	})
	if err != nil {
		log.Fatal(err)
	}

	/*
	 * Setup the serial connection and start the reader
	 */
	c := &serial.Config{Name: "/dev/ttyUSB0", Baud: 57600}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}
	defer s.Close()

	r := bufio.NewReader(s)
	go read(r, cli)

	/*
	 * Wait for interruption
	 */
	<-sigc

	// Disconnect the Network Connection.
	if err := cli.Disconnect(); err != nil {
		log.Fatal(err)
	}
}
