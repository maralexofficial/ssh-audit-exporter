FROM golang:1.26 AS builder

RUN apt-get update && apt-get install -y git ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o ssh-audit-exporter .

FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y \
    ca-certificates \
    systemd \
    && rm -rf /var/lib/apt/lists/*

ARG UID=1000
ARG GID=1000

RUN groupadd -g $GID exporter && useradd -m -u $UID -g exporter exporter

WORKDIR /app

COPY --from=builder /app/ssh-audit-exporter /app/ssh-audit-exporter

RUN chmod +x /app/ssh-audit-exporter

COPY .env.example /app/.env

RUN chown exporter:exporter /app/ssh-audit-exporter

USER exporter

EXPOSE 9100

CMD ["/app/ssh-audit-exporter"]
