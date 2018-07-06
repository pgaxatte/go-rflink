package rflink

import (
	"bufio"
	"log"

	"github.com/tarm/serial"
)

// SensorReader reads SensorData from the serial connection with RFLink
type SensorReader struct {
	port   *serial.Port
	reader *bufio.Reader
}

// NewSensorReader returns a SensorReader according to the options specified
func NewSensorReader(o *Options) (*SensorReader, error) {
	port, err := serial.OpenPort(&serial.Config{
		Name: o.Serial.Device,
		Baud: o.Serial.Baud,
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

// ReadNext reads a line from RFLink and returns it in the form of a SensorData
// struct
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

// Close closes the serial port
func (sr *SensorReader) Close() {
	sr.port.Close()
}
