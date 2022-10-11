# syntax=docker/dockerfile:1

FROM golang:1.19.2-alpine

ENV address=":8080"
ENV GO111MODULE=on
WORKDIR /usr/src
COPY go.mod server/go.sum ./
COPY ../proto/service proto/service
COPY server.go ./
RUN go mod download && go mod verify
RUN  go build
CMD ["sh", "-c", "./server $address"]
