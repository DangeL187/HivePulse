# build stage
FROM golang:1.25.3-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o consumer "./cmd"

# run stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/consumer .
CMD ["./consumer"]
