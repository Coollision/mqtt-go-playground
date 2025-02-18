FROM golang:1.23.2 AS bebuilder

ARG buildtags="offline"
ARG version="dontknow"

WORKDIR /app
ENV CGO_ENABLED=0
ENV GOMAXPROCS=4
COPY go.mod go.sum ./
RUN go env -w GOCACHE=/go-cache
RUN go env -w GOMODCACHE=/gomod-cache

RUN --mount=type=cache,target=/gomod-cache go mod download

COPY ./ ./


RUN go build -v -tags "$buildtags" -ldflags="-X main.version=$version" -o /app/mqtt-go-playground ./cmd/srvstart/

FROM alpine

RUN addgroup -S umqtt-go-playground && adduser -S umqtt-go-playground -G umqtt-go-playground
RUN mkdir /app && chown umqtt-go-playground:umqtt-go-playground /app

USER umqtt-go-playground

COPY --from=bebuilder /app/mqtt-go-playground ./mqtt-go-playground

CMD ["sh", "-c", "./mqtt-go-playground"]
