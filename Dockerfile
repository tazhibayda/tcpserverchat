FROM golang:1.22 AS builder

WORKDIR /app

COPY . .
COPY go.mod go.sum ./
RUN go mod download

COPY deployment .
RUN CGO_ENABLED=0 GOOS=linux go build -o tcpserverchat ./cmd/server


FROM alpine:latest AS runner

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app
COPY --from=builder /app/tcpserverchat .
EXPOSE 8000 9090 6060


ENTRYPOINT ["./tcpserverchat"]
