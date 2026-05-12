FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o server.bin ./server/cmd/main.go

FROM alpine:3.22
WORKDIR /app
COPY --from=builder /app/. .
EXPOSE 8080
CMD ["./server.bin"]