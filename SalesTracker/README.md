# SalesTracker

<p align="center">
  <img src="./docs/images/salestracker-banner.svg" alt="SalesTracker banner" width="100%" />
</p>

<p align="center">
  <b>SalesTracker</b> — веб-приложение для учёта продаж, управления записями и быстрой аналитики по выручке.
</p>

<p align="center">
  <img alt="Go" src="https://img.shields.io/badge/Go-backend-00ADD8?logo=go&logoColor=white">
  <img alt="PostgreSQL" src="https://img.shields.io/badge/PostgreSQL-database-336791?logo=postgresql&logoColor=white">
  <img alt="REST API" src="https://img.shields.io/badge/API-REST-111827">
  <img alt="HTML" src="https://img.shields.io/badge/UI-HTML%2FCSS-orange">
</p>

---

## Что умеет проект

- создавать записи о продажах;
- обновлять существующие продажи;
- удалять записи;
- получать одну запись по `id` и список всех записей;
- считать аналитику за период с фильтрами по категории и названию;
- работать через простой встроенный web-интерфейс.

## Почему проект полезен

SalesTracker закрывает базовый сценарий учёта продаж без тяжёлого фронтенда и без лишней инфраструктуры:

- **CRUD API** для управления продажами;
- **аналитика по периоду** через отдельный endpoint;
- **валидация данных** на уровне приложения и БД;
- **PostgreSQL** как надёжное хранилище;
- **минималистичный UI** для ручной работы и быстрой проверки API.

---

## Стек

- **Go**
- **PostgreSQL 16**
- **REST API**
- **HTML/CSS**
- Вспомогательные пакеты:
  - `github.com/wb-go/wbf/dbpg`
  - `github.com/wb-go/wbf/ginext`
  - `github.com/wb-go/wbf/zlog`
  - `github.com/wb-go/wbf/config`

---

## Структура проекта

```text
SalesTracker/
├── cmd/
│   └── main.go
├── internal/
│   ├── appCfg/
│   │   └── config.go
│   ├── customErrs/
│   │   └── errs.go
│   ├── handlers/
│   │   └── handlers.go
│   ├── migrations/
│   │   └── 20260315092450_create_sales_table.sql
│   ├── models/
│   │   └── models.go
│   └── repository/
│       └── sales.go
├── web/
│   └── index.html
├── config-example.yaml
└── docker-compose.yml.example
```

### Как устроено внутри

- `cmd/main.go` — точка входа, инициализация конфига, логгера, БД и HTTP-роутов.
- `internal/appCfg` — загрузка конфигурации из `config.yaml`.
- `internal/handlers` — HTTP-обработчики для CRUD и аналитики.
- `internal/repository` — SQL-запросы к PostgreSQL.
- `internal/models` — модели данных и валидация.
- `internal/migrations` — создание таблицы `sales`.
- `web/index.html` — простой UI для работы с продажами и аналитикой.

---

## Быстрый старт

> Ниже — самый практичный путь для локального запуска: PostgreSQL в Docker, приложение локально.

### 1. Подними PostgreSQL

Скопируй пример compose-файла:

```bash
cp docker-compose.yml.example docker-compose.yml
docker compose up -d
```

По умолчанию база поднимется с параметрами:

- host: `localhost`
- port: `5433`
- db: `salestracker`
- user: `admin`
- password: `admin`

### 2. Создай конфиг приложения

Скопируй пример:

```bash
cp config-example.yaml config.yaml
```

Пример рабочего `config.yaml`:

```yaml
server:
  addr: ":8080"

logger:
  level: "debug"

postgres:
  max_open_conns: 5
  max_idle_conns: 3
  conn_max_lifetime: 10s
  port: 5433
  master_dsn: "postgres://admin:admin@localhost:5433/salestracker?sslmode=disable"
  slave_dsn:
    - "postgres://admin:admin@localhost:5433/salestracker?sslmode=disable"
```

### 3. Создай таблицу `sales`

В проекте уже есть SQL-миграция. Можно применить её вручную.

Сначала подключись к БД:

```bash
psql "postgres://admin:admin@localhost:5433/salestracker?sslmode=disable"
```

Затем выполни:

```sql
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS sales(
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  title VARCHAR(255) NOT NULL,
  category VARCHAR(255) NOT NULL,
  price DECIMAL(10,2) NOT NULL CHECK (price > 0.0),
  quantity INT NOT NULL CHECK (quantity > 0),
  sale_date TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

> Если у тебя уже настроен `goose`, можешь применить миграцию из `internal/migrations/`.

### 4. Запусти приложение

Важно: приложение ожидает `config.yaml` по пути `../config.yaml` и web-файл `../web/index.html`, поэтому запускать удобнее из папки `cmd`.

```bash
cd cmd
go run .
```

После запуска приложение будет доступно по адресу:

```text
http://localhost:8080
```

---

## API

### Создать продажу

`POST /items`

```json
{
  "title": "iPhone 15",
  "category": "electronics",
  "price": 999.99,
  "quantity": 2,
  "sale_date": "2026-03-15T12:30:00Z"
}
```

Если `sale_date` не передать, дата продажи будет выставлена автоматически на стороне БД.

### Получить все продажи

`GET /items`

### Получить продажу по ID

`GET /items/:id`

### Обновить продажу

`PUT /items/:id`

```json
{
  "title": "iPhone 15 Pro",
  "category": "electronics",
  "price": 1199.99,
  "quantity": 1,
  "sale_date": "2026-03-15T12:30:00Z"
}
```

### Удалить продажу

`DELETE /items/:id`

### Получить аналитику

`GET /analytics?from=2026-03-01&to=2026-03-31`

Фильтры:

- `from` — обязательный параметр;
- `to` — обязательный параметр;
- `category` — опционально;
- `title` — опционально.

Пример:

```http
GET /analytics?from=2026-03-01&to=2026-03-31&category=electronics
```

Пример ответа:

```json
{
  "sum": 4599.95,
  "avg": 1533.3166666667,
  "count": 3,
  "median": 1199.99,
  "p90": 2199.99
}
```

---

## Web UI

Встроенный интерфейс позволяет:

- добавить продажу;
- посмотреть список всех записей;
- фильтровать таблицу;
- запрашивать аналитику по диапазону дат;
- быстро тестировать backend без Postman.

> В текущем UI фильтрация таблицы выполняется на фронте после `GET /items`, а агрегированная аналитика считается через backend endpoint `/analytics`.

---

## Валидация

На уровне приложения запись считается корректной, если:

- `price > 0`
- `quantity > 0`
- `title` не пустой
- `category` не пустая

Дополнительно ограничения продублированы в PostgreSQL через `CHECK`.

---

## Пример сценария использования

1. Пользователь открывает web-интерфейс.
2. Добавляет несколько продаж в разные категории.
3. Получает список всех записей.
4. Фильтрует нужные продажи по названию или категории.
5. Запрашивает аналитику за период:
   - сумму,
   - среднее,
   - количество,
   - медиану,
   - 90-й перцентиль.

---