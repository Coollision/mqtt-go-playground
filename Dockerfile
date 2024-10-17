FROM golang:alpine AS builder
ARG buildtags=""
ARG version="none given"

RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go build -v -tags "$buildtags" -ldflags="-X main.version=$version" -o mqtt-go-playground .

FROM alpine
COPY --from=builder /app/mqtt-go-playground .
ENTRYPOINT ["/mqtt-go-playground"]
