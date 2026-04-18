FROM golang:1.26 AS builder

RUN apt-get update && apt-get install -y git ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ssh-audit-exporter .

FROM alpine:3.20

RUN apk add --no-cache ca-certificates util-linux

ARG UID=1000
ARG GID=1000

RUN addgroup -g $GID exporter && adduser -D -u $UID -G exporter exporter

WORKDIR /app

COPY --from=builder /app/ssh-audit-exporter /app/ssh-audit-exporter

RUN chmod +x /app/ssh-audit-exporter

COPY .env.example /app/.env

RUN chown exporter:exporter /app/ssh-audit-exporter

USER exporter

EXPOSE 9100

CMD ["/app/ssh-audit-exporter"]
