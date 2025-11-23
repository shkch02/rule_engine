FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy && go mod download

COPY . .

RUN CGO_ENABLE=0 go build -o rule-engine ./cmd/engine/

FROM alpine:latest


COPY --from=builder /app/rule-engine /usr/local/bin/rule-engine
COPY --from=builder /app/test_logs.txt /app/test_logs.txt

ENTRYPOINT ["rule-engine", "/app/test_logs.txt"]   