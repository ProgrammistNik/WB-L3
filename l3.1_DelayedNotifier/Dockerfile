FROM golang:1.23.0-alpine

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o app ./cmd/main.go

EXPOSE 8080

CMD ["./app"]
