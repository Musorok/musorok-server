FROM golang:1.22 AS builder
WORKDIR /app
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/server

FROM gcr.io/distroless/base-debian12
WORKDIR /
ENV TZ=Asia/Almaty
COPY --from=builder /app/server /server
COPY .env /.env
COPY docs /docs
EXPOSE 8080
USER 65532:65532
ENTRYPOINT ["/server"]
