# WarehouseControl

![Go](https://img.shields.io/badge/Go-backend-00ADD8?logo=go&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-database-4169E1?logo=postgresql&logoColor=white)
![JWT](https://img.shields.io/badge/Auth-JWT-black)
![RBAC](https://img.shields.io/badge/Access-RBAC-success)
![Audit Log](https://img.shields.io/badge/Audit-enabled-orange)

WarehouseControl — это backend + web UI для управления складскими товарами с аутентификацией, ролевой моделью доступа и журналом изменений.

Проект решает базовые задачи складского учёта:
- хранение товаров;
- создание, редактирование и удаление записей;
- разграничение доступа по ролям;
- отслеживание истории изменений через audit log.

---

## Содержание

- [Возможности](#возможности)
- [Стек](#стек)
- [Архитектура](#архитектура)
- [Роли и права доступа](#роли-и-права-доступа)
- [API](#api)
- [Быстрый старт](#быстрый-старт)
- [Конфигурация](#конфигурация)
- [База данных](#база-данных)
- [Аудит изменений](#аудит-изменений)
- [Web-интерфейс](#web-интерфейс)
- [Пример сценария использования](#пример-сценария-использования)
- [Почему этот проект полезен](#почему-этот-проект-полезен)

---

## Возможности

- JWT-аутентификация
- Ролевой доступ (`admin`, `manager`, `viewer`)
- CRUD-операции для товаров
- Просмотр списка товаров
- История изменений по каждому товару
- PostgreSQL в качестве основной БД
- Web-интерфейс для работы с системой
- Хеширование паролей
- Audit log для `INSERT`, `UPDATE`, `DELETE`

---

## Стек

**Backend**
- Go

**Database**
- PostgreSQL

**Auth**
- JWT
- bcrypt

**Frontend**
- HTML / CSS / JavaScript

**Инфраструктура**
- Docker Compose для PostgreSQL

---

## Архитектура

Проект организован по слоям:

```text
WarehouseControl/
├── cmd/                    # точка входа приложения
│   └── main.go
├── internal/
│   ├── appCfg/             # загрузка конфигурации
│   ├── auth/               # JWT и работа с паролями
│   ├── customErrs/         # пользовательские ошибки
│   ├── handlers/           # HTTP handlers
│   ├── middleware/         # auth и role middleware
│   ├── migrations/         # SQL-миграции
│   ├── models/             # DTO и модели ответа
│   ├── repository/         # доступ к данным
│   └── service/            # бизнес-логика
├── web/
│   └── index.html          # web UI
├── config.yaml.example
└── docker-compose.yml.example
```

Такое разделение упрощает поддержку проекта: HTTP-слой не смешивается с бизнес-логикой, а работа с БД изолирована в `repository`.

---

## Роли и права доступа

В системе есть 3 роли:

### `admin`
Полный доступ:
- просмотр товаров
- создание товаров
- редактирование товаров
- удаление товаров
- просмотр audit log
- просмотр `/me`

### `manager`
Ограниченный доступ к управлению складом:
- просмотр товаров
- создание товаров
- редактирование товаров

### `viewer`
Доступ только на чтение:
- просмотр списка товаров

---

## API

### Аутентификация

#### `POST /auth/login`

Авторизация пользователя.

**Request**
```json
{
  "username": "admin",
  "password": "secret"
}
```

**Response**
```json
{
  "username": "admin",
  "role": "admin",
  "token": "jwt-token"
}
```

---

### Пользователь

#### `GET /me`

Информация о текущем пользователе.

> Доступно только для `admin`

---

### Товары

#### `POST /items`
Создание товара.

> Доступно для `admin`, `manager`

**Request**
```json
{
  "title": "Ручка",
  "sku": "PEN-001",
  "quantity": 100
}
```

---

#### `GET /items`
Получение списка товаров.

> Доступно для `admin`, `manager`, `viewer`

---

#### `PUT /items/:id`
Обновление товара.

> Доступно для `admin`, `manager`

**Request**
```json
{
  "title": "Ручка синяя",
  "sku": "PEN-001",
  "quantity": 120
}
```

---

#### `DELETE /items/:id`
Удаление товара.

> Доступно только для `admin`

---

### История изменений

#### `GET /items/:id/audit`

Получение истории изменений конкретного товара.

> Доступно только для `admin`

---

## Быстрый старт

### 1. Клонировать репозиторий

```bash
git clone https://github.com/A13xSuX/L3.git
cd L3/WarehouseControl
```

### 2. Поднять PostgreSQL

```bash
cp docker-compose.yml.example docker-compose.yml
docker compose up -d
```

### 3. Создать конфиг

```bash
cp config.yaml.example config.yaml
```

### 4. Заполнить параметры подключения

Пример `config.yaml`:

```yaml
server:
  addr: ":8080"

logger:
  level: debug

postgres:
  max_open_conns: 5
  max_idle_conns: 3
  conn_max_lifetime: 10s
  port: 5433
  master_dsn: "postgres://user:password@host/dbname?sslmode=disable"
  slave_dsn: []

jwt:
  secret_key: "secret_key"
  ttl: 24h
```

### 5. Применить миграции

Миграции лежат в каталоге:

```text
internal/migrations/
```

Там создаются таблицы:
- `items`
- `users`
- `audit`

А также триггер для автоматического аудита изменений.

### 6. Запустить приложение

```bash
go run ./cmd
```

После запуска приложение:
- читает конфиг;
- подключается к PostgreSQL;
- проверяет соединение запросом `SELECT 1`;
- поднимает HTTP-сервер.

---

## Конфигурация

Основные параметры приложения:

| Параметр | Описание |
|---|---|
| `server.addr` | адрес HTTP-сервера |
| `logger.level` | уровень логирования |
| `postgres.master_dsn` | основное подключение к PostgreSQL |
| `postgres.slave_dsn` | список read-replica подключений |
| `jwt.secret_key` | секрет для подписи JWT |
| `jwt.ttl` | время жизни токена |

---

## База данных

### Таблица `items`
Хранит товары:
- `id`
- `title`
- `sku`
- `quantity`
- `created_at`
- `updated_at`

Ограничения:
- `sku` — уникальный
- `quantity >= 0`

### Таблица `users`
Хранит пользователей:
- `username`
- `password_hash`
- `role`

Допустимые роли:
- `admin`
- `manager`
- `viewer`

### Таблица `audit`
Хранит историю изменений:
- тип действия (`INSERT`, `UPDATE`, `DELETE`)
- предыдущее состояние
- новое состояние
- кто выполнил изменение
- когда изменение произошло

---

## Аудит изменений

В проекте реализован audit log через PostgreSQL trigger.

Что логируется:
- создание товара;
- обновление товара;
- удаление товара.

Что сохраняется:
- `old_data`
- `new_data`
- `changed_by_username`
- `changed_by_role`
- `changed_at`

Это позволяет:
- отслеживать, кто изменил данные;
- восстанавливать историю изменений;
- анализировать действия пользователей.

---

## Web-интерфейс

В проекте есть встроенная web-страница:

```text
/web
```

Интерфейс позволяет:
- войти в систему;
- просматривать товары;
- создавать и редактировать товары;
- смотреть историю изменений.

Ограничения по ролям также отражены в интерфейсе:
- создание и редактирование доступны только `admin` и `manager`.

---

## Пример сценария использования

1. Пользователь входит в систему через `/auth/login`
2. Получает JWT-токен
3. Использует токен для запросов к `/items`
4. `manager` создаёт или обновляет товар
5. PostgreSQL trigger автоматически пишет изменение в `audit`
6. `admin` может посмотреть историю через `/items/:id/audit`

---

## Почему этот проект полезен

WarehouseControl — хороший пример backend-приложения, в котором сочетаются:
- аутентификация;
- RBAC;
- работа с PostgreSQL;
- слоистая архитектура;
- аудит действий пользователей;
- простой UI для ручного тестирования API.

Такой проект хорошо подходит для:
- портфолио;
- pet-project / учебного проекта;
- демонстрации backend-навыков на Go.

---