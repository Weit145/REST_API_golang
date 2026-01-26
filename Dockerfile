# ---------- Stage 1: Build ----------
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Устанавливаем gcc и libc для go-sqlite3
RUN apk add --no-cache gcc musl-dev

# Кэш зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходники
COPY . .

# Включаем CGO
ENV CGO_ENABLED=1 GOOS=linux GOARCH=amd64

# Сборка бинарника
RUN go build -ldflags="-s -w" -o app ./cmd/app/main.go

# ---------- Stage 2: Runtime ----------
FROM alpine:latest

# Нужные пакеты для работы приложения и SQLite
RUN apk add --no-cache bash ca-certificates sqlite

WORKDIR /app

# Создаём папку для SQLite базы
RUN mkdir -p /app/storage

# Копируем бинарник из Stage 1
COPY --from=builder /app/app .

# Копируем конфиг
COPY config/local.yaml /app/config/local.yaml

ENV CONFIG_PATH=/app/config/local.yaml

EXPOSE 8080

CMD ["./app"]
