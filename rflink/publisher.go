package rflink

import (
	"encoding/json"
	"log"

	"github.com/yosssi/gmq/mqtt/client"
)

// Publisher takes input from a SensorReader and publishes the SensorData that
// has been read in an MQTT topic
type Publisher struct {
	c *client.Client

	Topic       string
	SensorInput *SensorReader
}

// NewPublisher return a Publisher according to the options specified
func NewPublisher(o *Options) (*Publisher, error) {
	cli := client.New(&client.Options{
		ErrorHandler: func(err error) {
			log.Printf("Error from MQTT client: %s", err)
		},
	})

	err := cli.Connect(&client.ConnectOptions{
		Network:  o.PublishURL.Scheme,
		Address:  o.PublishURL.Host,
		ClientID: []byte(o.PublishClientID),
	})
	if err != nil {
		return nil, err
	}
	p := &Publisher{
		c:     cli,
		Topic: o.PublishTopic,
	}
	return p, nil
}

// Publish formats the input SensorData into JSON and publishes it to the
// configured MQTT topic
func (p *Publisher) Publish(sd *SensorData) error {
	b, err := json.Marshal(sd)
	if err != nil {
		return err
	}
	log.Print(string(b))

	err = p.c.Publish(&client.PublishOptions{
		TopicName: []byte(p.Topic),
		Message:   b,
	})
	if err != nil {
		return err
	}

	return nil
}

// ReadAndPublish loops infinitely to read SensorData from the SensorReader and
// publish the output via the Publish() method
func (p *Publisher) ReadAndPublish() error {
	for {
		sd, err := p.SensorInput.ReadNext()
		if err != nil {
			log.Print(err)
			continue
		}

		err = p.Publish(sd)
		if err != nil {
			log.Print(err)
			continue
		}
	}

	return nil
}

// Disconnect properly disconnects the MQTT network connection
func (p *Publisher) Disconnect() error {
	return p.c.Disconnect()
}

// Terminate kills the MQTT client
func (p *Publisher) Terminate() {
	p.c.Terminate()
}
