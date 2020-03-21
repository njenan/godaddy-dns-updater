FROM golang:1.14.0-alpine3.11

ENV GOPROXY=athens.nathanjenan.me

WORKDIR /app
COPY . .

RUN go install -ldflags "-X main.Version=$(cat version) -X main.License=GPL-2.0" .

