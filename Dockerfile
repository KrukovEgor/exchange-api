FROM golang:1.27 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY ./ ./

RUN CGO_ENABLED=0 GOOS=linux go build -o ./build/api ./cmd/api/main.go

FROM alpine:latest

WORKDIR /root

COPY --from=builder ./app ./

CMD ["./build/api", "-c", "./configs/backend.yml"]