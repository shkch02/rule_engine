FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy && go mod download

COPY . .

RUN CGO_ENABLE=0 go build -o rule-engine ./cmd/engine/

FROM alpine:latest

COPY --from=builder /app/rule-engine /usr/local/bin/rule-engine

ENTRYPOINT ["rule-engine", "./test_log.txt"]