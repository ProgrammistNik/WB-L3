# Stage 1: build
FROM --platform=linux/amd64 golang:1.22.6-bookworm AS builder

WORKDIR /app

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        gcc \
        librdkafka-dev \
        pkg-config && \
    rm -rf /var/lib/apt/lists/*

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Собираем бинарник
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /app/bin/app ./cmd

# Stage 2: runtime
FROM --platform=linux/amd64 debian:bookworm-slim

WORKDIR /app

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        librdkafka1 \
        ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/bin/app .
COPY --from=builder /app/config ./config
COPY --from=builder /app/web ./web
COPY --from=builder /app/uploads ./uploads

RUN useradd -m appuser && chown -R appuser:appuser /app
USER appuser

EXPOSE 8080

ENTRYPOINT ["./app"]
