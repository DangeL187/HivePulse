# build stage
FROM golang:1.24-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o auth "./cmd"

# run stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/auth .
COPY --from=builder /app/casbin_model.conf .
CMD ["./auth"]
