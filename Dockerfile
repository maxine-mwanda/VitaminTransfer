# Build stage
FROM golang:alpine AS builder
RUN apk add --no-cache git

WORKDIR /go/app
COPY . .

# Ensure Go modules are used properly
RUN go mod tidy

# Build the application
RUN go build -o /go/bin/app .


# Final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates

COPY --from=builder /go/bin/app /app
COPY .env .env 
ENTRYPOINT ["/app"]

LABEL Name=vitamintransfer Version=0.0.1
EXPOSE 8080

