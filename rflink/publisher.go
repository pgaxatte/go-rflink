package rflink

import (
	"encoding/json"
	"log"

	"github.com/yosssi/gmq/mqtt/client"
)

type Publisher struct {
	c *client.Client

	Topic       string
	SensorInput *SensorReader
}

type PublisherOptions struct {
	ConnectOptions *client.ConnectOptions
	Topic          string
}

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

func (p *Publisher) Disconnect() error {
	return p.c.Disconnect()
}

func (p *Publisher) Terminate() {
	p.c.Terminate()
}
