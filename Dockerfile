FROM golang:alpine

WORKDIR /go/src/github.com/pgaxatte/go-rflink
COPY . .

RUN go install -v ./...

CMD ["go-rflink"]
