# Ticket System

Система управления тикетами с backend на Go и frontend на Qt/C++.

---

## Архитектура

- **Backend:** Go (Gin, GORM, PostgreSQL, JWT)
- **Frontend:** Qt 6 (C++), QML, CMake
- **База данных:** PostgreSQL
- **Docker:** Поддержка контейнеризации для быстрого старта

---

## Возможности

- Регистрация и аутентификация пользователей (JWT)
- CRUD для тикетов, комментариев, вложений
- Справочники: департаменты, статусы, приоритеты
- Роли и права (админ, пользователь)
- Загрузка и удаление вложений (attachments) с предпросмотром изображений
- Lightbox-модалка для просмотра изображений
- Скачивание файлов вложений
- Современный UI на английском языке
- Подробное логирование и валидация

---

## Структура проекта

- `backend/` — исходный код backend (Go)
  - `cmd/api/main.go` — точка входа
  - `internal/` — бизнес-логика, delivery, middleware, usecase, repository, domain, model
  - `migrations/` — SQL-миграции для PostgreSQL
  - `Dockerfile` — сборка backend
- `frontend/` — исходный код frontend (Qt/C++)
  - `src/` — исходники, модели, views, network, utils
  - `config.ini` — конфигурация клиента
  - `CMakeLists.txt` — сборка
- `docker-compose.yml` — запуск всех сервисов
- `start.sh` — удобный запуск через Docker

---

## Быстрый старт (Docker)

```bash
cd ticket-system
./start.sh
```

- Backend: http://localhost:8080
- Frontend: соберите вручную (см. ниже)
- Adminer (UI для БД): http://localhost:8081

---

## Ручной запуск

### Backend

```bash
cd backend
# Установить зависимости
go mod tidy
# Запустить сервер
go run cmd/api/main.go
```

### Frontend

```bash
cd frontend
mkdir -p build && cd build
cmake ..
make
./ticket_frontend
```

---

## Конфигурация

### Backend (env)
- `DB_URL` — строка подключения к PostgreSQL
- `JWT_SECRET` — секрет для JWT
- `SERVER_PORT` — порт (по умолчанию 8080)
- `LOG_LEVEL` — уровень логирования (info, debug)

### Frontend (`config.ini`)
```ini
[api]
base_url=http://localhost:8080
version=v1
[ui]
theme=light
language=en
[network]
timeout=30000
retry_attempts=3
```

---

## Основные API эндпоинты

### Аутентификация
- `POST /api/v1/auth/login` — вход
- `POST /api/v1/auth/register` — регистрация

### Тикеты
- `GET /api/v1/tickets` — список
- `POST /api/v1/tickets` — создать
- `GET /api/v1/tickets/:id` — получить по ID
- `PATCH /api/v1/tickets/:id` — обновить
- `DELETE /api/v1/tickets/:id` — удалить

### Комментарии
- `GET /api/v1/tickets/:id/comments` — список
- `POST /api/v1/tickets/:id/comments` — добавить
- `DELETE /api/v1/comments/:comment_id` — удалить

### Вложения (attachments)
- `GET /api/v1/tickets/:id/attachments` — список вложений тикета
- `POST /api/v1/tickets/:id/attachments` — загрузить файл (multipart/form-data, поле file)
- `GET /api/v1/attachments/:att_id/download` — скачать файл
- `DELETE /api/v1/attachments/:att_id` — удалить (только свои или если админ)

### Справочники
- `GET /api/v1/departments` — департаменты
- `GET /api/v1/ticket_statuses` — статусы
- `GET /api/v1/ticket_priorities` — приоритеты

---

## Модели данных (основные)

### Ticket
- id, title, description, status, priority, department, assignee, creator, createdAt, updatedAt

### Comment
- id, ticketId, authorId, authorName, content, createdAt

### Attachment
- id, ticketId, filename, filePath, uploadedBy, uploadedAt
- Для изображений поддерживается предпросмотр и lightbox

---

## Работа с вложениями
- Загрузка через drag&drop или диалог выбора файла
- Предпросмотр миниатюр изображений (jpg, png, gif)
- Lightbox-модалка для увеличения изображения (клик по миниатюре)
- Скачивание любого файла вложения
- Удаление вложения (только свои или если админ)
- Проверка прав на backend

---

## Миграции и база данных

- Используется PostgreSQL
- Миграции: `backend/migrations/`
- Для применения вручную:
```bash
cd backend
# Установить migrate
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
# Применить миграции
migrate -path migrations -database "postgres://user:password@localhost:5432/ticket_system?sslmode=disable" up
```

---

## Безопасность
- Хеширование паролей (bcrypt)
- JWT для аутентификации
- Валидация входных данных
- Защита от SQL-инъекций (GORM)
- Проверка прав на действия (удаление вложений, тикетов и т.д.)
- Логирование ошибок и действий

---

## Известные проблемы и рекомендации
- Все русские сообщения и комментарии переведены на английский
- Рекомендуется работать через feature-ветки, сливать в main после тестирования
- Для разработки используйте переменные окружения и отдельные конфиги

---

## Требования
- Go 1.19+
- Qt 6.0+
- PostgreSQL 12+
- CMake 3.16+
- Docker (для быстрого старта)

---

## Контакты и поддержка

Для вопросов и предложений — создавайте issue или pull request.
