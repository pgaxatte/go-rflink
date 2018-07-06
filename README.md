[![Godoc](https://godoc.org/github.com/pgaxatte/go-rflinki/rflink?status.svg)](https://godoc.org/github.com/pgaxatte/go-rflink/rflink)

# go-rflink
Publish rflink temperature and humidity measurement to an MQTT topic.


## Installation

```bash
go get -u github.com/pgaxatte/go-rflink
```

## Usage

Optionnal environment variable can be used to override the default configuration, for example:

```bash
PUBLISH_HOST=192.168.0.1:1883 SERIAL_DEVICE=/dev/ttyACM0 go run main.go
```

See the [Options struct definition](https://godoc.org/github.com/pgaxatte/go-rflink/rflink#Options) for a complete list of supported options.

### Within a docker container

It is possible to build go-rflink as a container and run it:

```bash
docker build -t go-rflink .

docker run \
    --device=/dev/ttyACM0 \
    --env PUBLISH_HOST="192.168.0.1:1883" \
    --env PUBLISH_TOPIC="myrflink" \
    go-rflink:latest
```

# TODO
- [ ] Try to reconnect (indefinetly or a limited number of times) when MQTT or USB connection fails
- [ ] Add tests
