FROM golang:1.26.2-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./cmd ./cmd
COPY ./internal ./internal

RUN CGO_ENABLED=0 GOOS=linux go build -o /url_saver ./cmd/api

FROM alpine:3.19
WORKDIR /

COPY --from=builder /url_saver /url_saver

EXPOSE 8080

CMD ["/url_saver"]