# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -v -o loxone-vo-json .

# Final image
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/loxone-vo-json .
EXPOSE 8080
ENTRYPOINT ["./loxone-vo-json"]