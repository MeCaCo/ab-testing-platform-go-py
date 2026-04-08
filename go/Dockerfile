FROM golang:1.26-alpine AS builder

WORKDIR /build

# Копируем go.mod и go.sum
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь код
COPY . .

# Собираем бинарник
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/server/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Копируем бинарник из builder-образа
COPY --from=builder /build/main .

# Копируем конфиг
COPY config/config.yaml ./config/config.yaml

EXPOSE 8080

CMD ["./main"]