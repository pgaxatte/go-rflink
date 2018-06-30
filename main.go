package main

import (
	"bufio"
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/tarm/serial"
)

type SensorData struct {
	Manufacturer string
	Id           uint16
	Temperature  float32
	Humidity     uint16
}

func strToUint16(s string, base int) (uint16, error) {
	u, err := strconv.ParseUint(s, base, 16)
	if err != nil {
		return 0, err
	}
	return uint16(u), nil
}

func parse(msg string) (*SensorData, error) {
	log.Printf(msg)
	pieces := strings.Split(msg, ";")

	sd := SensorData{
		Manufacturer: pieces[2],
	}
	for i := 3; i < len(pieces); i++ {
		arr := strings.SplitN(pieces[i], "=", 2)
		switch arr[0] {
		case "ID":
			id, err := strToUint16(arr[1], 16)
			if err != nil {
				return nil, errors.New("Skipping message, id could not be parsed")
			}
			//log.Printf("ID: %d (%s)", id, arr[1])
			sd.Id = id
		case "TEMP":
			t, err := strToUint16(arr[1], 16)
			if err != nil {
				return nil, errors.New("Skipping message, temperature could not be parsed")
			}
			//log.Printf("Temperature: %d (%s)", t, arr[1])
			sd.Temperature = float32(t) / 10
		case "HUM":
			h, err := strToUint16(arr[1], 10)
			if err != nil {
				return nil, errors.New("Skipping message, humidity could not be parsed")
			}
			//log.Printf("Humidity: %d (%s)", h, arr[1])
			sd.Humidity = h
		}
	}
	return &sd, nil
}

func read(r *bufio.Reader) {
	for {
		line, _, err := r.ReadLine()
		if err != nil {
			log.Fatal(err)
		}

		sd, err := parse(string(line))
		if err != nil {
			log.Printf("Error parsing \"%s\": %s", line, err)
		}

		log.Printf("DATA: %+v", sd)
	}
}

func main() {
	c := &serial.Config{Name: "/dev/ttyUSB0", Baud: 57600}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}
	defer s.Close()

	r := bufio.NewReader(s)
	read(r)
}
