package rflink

import (
	"errors"
	"fmt"
	"strings"
)

// SensorData represents one message received from RFLink.
// It can only contain Temperature or Humidity values and exposes the model
// id of the sending device.
// A friendly name can be added to the data so that the Publisher can tag it
// before sending it to the MQTT topic.
type SensorData struct {
	Model        string   `json:"model"`
	Id           string   `json:"id"`
	FriendlyName string   `json:"name,omitempty"`
	Temperature  *float32 `json:"t,omitempty"`
	Humidity     *uint16  `json:"h,omitempty"`
}

// SensorDataFromMessage crafts a SensorData struct from a message read from
// RFLink
func SensorDataFromMessage(msg string) (*SensorData, error) {
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

	if sd.Temperature == nil && sd.Humidity == nil {
		return nil, errors.New("Skipping message, no temperature nor humidity")
	}

	return &sd, nil
}

// String outputs a string representing the SensorData
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
