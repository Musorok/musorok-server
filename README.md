# MusorOK Server (Go + Gin + Postgres + Redis)

Готовый стартовый сервер под **Gin** (почему Gin: самый популярный Go-фреймворк, высокая производительность, зрелая экосистема middleware и примеров).

## Быстрый старт (Docker)

```bash
docker compose up -d --build
make migrate
make seed
```

Открой: http://localhost:8080/docs

## Локально без Docker
1. Подними Postgres 15 и Redis 7
2. Настрой `.env`
3. `make migrate && make seed && make run`

## Что внутри
- JWT (access+refresh), bcrypt-хеширование
- Роли: USER/COURIER/ADMIN
- Адреса с геопроверкой против полигонов (ЖК 4YOU сид)
- Заказы (quote + создание), WebSocket задел, платежи (заглушка paynetworks)
- Swagger UI на `/docs` + `docs/openapi.yaml`
