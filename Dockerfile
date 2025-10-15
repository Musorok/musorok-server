FROM golang:1.22 AS builder
WORKDIR /app

# 1) Только модули — чтобы кэшировалось
COPY go.mod go.sum* ./
RUN go mod download

# 2) Код проекта
COPY . .

# 3) ВАЖНО: синхронизируем зависимости внутри контейнера
RUN go mod tidy

# 4) Сборка бинарника с подробным выводом
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o /server ./cmd/server

# --- Runtime слой (можно оставить distroless) ---
FROM gcr.io/distroless/base-debian12
WORKDIR /
ENV TZ=Asia/Almaty
COPY --from=builder /server /server
COPY .env /.env
COPY docs /docs
EXPOSE 8080
USER 65532:65532
ENTRYPOINT ["/server"]
