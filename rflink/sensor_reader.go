package rflink

import (
	"bufio"
	"log"

	"github.com/tarm/serial"
)

type SensorReader struct {
	port   *serial.Port
	reader *bufio.Reader
}

func NewSensorReader(o *Options) (*SensorReader, error) {
	port, err := serial.OpenPort(&serial.Config{
		Name: o.SerialDevice,
		Baud: o.SerialBaud,
	})
	if err != nil {
		return nil, err
	}

	sr := &SensorReader{
		port:   port,
		reader: bufio.NewReader(port),
	}
	return sr, nil
}

func (sr *SensorReader) ReadNext() (*SensorData, error) {
	line, _, err := sr.reader.ReadLine()
	if err != nil {
		log.Printf("Cannot read from serial: %s", err)
		return nil, err
	}

	sd, err := SensorDataFromMessage(string(line))
	if err != nil {
		log.Printf("Error parsing message from rflink \"%s\": %s", line, err)
		return nil, err
	}

	return sd, nil
}

func (sr *SensorReader) Close() {
	sr.port.Close()
}
