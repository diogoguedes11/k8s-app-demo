# Build stage
FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY server.go .

RUN go build server.go

FROM alpine:latest

COPY --from=builder /app/server /server

EXPOSE 8080

CMD ["/server"]