FROM golang:1.26 AS builder

RUN apt-get update && apt-get install -y gcc libc6-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o exporter .

FROM alpine:3.20

RUN apk add --no-cache ca-certificates

ARG UID=1000
ARG GID=1000

RUN addgroup -g $GID exporter && adduser -D -u $UID -G exporter exporter

WORKDIR /app

COPY --from=builder /app/exporter .

COPY .env.example .env

RUN chown exporter:exporter /app/exporter

USER exporter

EXPOSE 9100

CMD ["./exporter"]
