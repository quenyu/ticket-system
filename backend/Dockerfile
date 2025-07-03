# syntax=docker/dockerfile:1
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN cd cmd/api && go build -o /app/api

FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/api ./api
EXPOSE 8080
CMD ["./api"] 