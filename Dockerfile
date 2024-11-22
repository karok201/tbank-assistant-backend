FROM golang:alpine AS builder

WORKDIR /app

COPY main.go go.mod go.sum ./

RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o server

COPY docker/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

COPY . /app

FROM alpine:latest

COPY main.go go.mod go.sum ./

ADD go.mod .

ADD go.sum .

RUN go mod download

ENTRYPOINT ["/entrypoint.sh"]
