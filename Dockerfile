FROM golang:1.25.1-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -ldflags="-s -w" -o /app/server cmd/server/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/server .

COPY web/ ./web/

EXPOSE 8080

CMD ["./server"]
