# syntax=docker/dockerfile:1

FROM golang:1.19.2-alpine

ENV address=":8080"
ENV GO111MODULE=on
WORKDIR /usr/src
RUN apk update && apk add bash ca-certificates git gcc g++ libc-dev
COPY go.mod go.sum ./
COPY proto/service proto/service
COPY server ./server
RUN go mod download && go mod verify
RUN cd server && go build
CMD ["sh", "-c", "./server/server $address"]
