# Ticket System

Система управления тикетами с backend на Go и frontend на Qt/C++.

## Исправления и улучшения

### Backend (Go)
- ✅ Добавлены JSON теги для всех структур домена
- ✅ Улучшено логирование с кастомным middleware
- ✅ Добавлена валидация входных данных
- ✅ Исправлена структура APIError для лучшей обработки ошибок
- ✅ Добавлены проверки UUID и числовых параметров

### Frontend (Qt/C++)
- ✅ Создана система конфигурации (config.h/config.cpp)
- ✅ Убраны хардкод URL, теперь используются настройки
- ✅ Улучшена обработка ошибок JSON
- ✅ Добавлена валидация на стороне клиента
- ✅ Исправлены проблемы с загрузкой ComboBox
- ✅ Добавлены дополнительные проверки null pointer

## Запуск

### Backend

1. Перейдите в папку backend:
```bash
cd ticket-system/backend
```

2. Установите зависимости:
```bash
go mod tidy
```

3. Запустите сервер:
```bash
go run cmd/api/main.go
```

Сервер будет доступен по адресу: http://localhost:8080

### Frontend

1. Перейдите в папку frontend:
```bash
cd ticket-system/frontend
```

2. Создайте build директорию:
```bash
mkdir build
cd build
```

3. Сконфигурируйте проект:
```bash
cmake ..
```

4. Соберите проект:
```bash
make
```

5. Запустите приложение:
```bash
./ticket_frontend
```

## Конфигурация

### Backend
Настройки загружаются из переменных окружения:
- `DB_URL` - URL базы данных PostgreSQL
- `JWT_SECRET` - секретный ключ для JWT токенов
- `SERVER_PORT` - порт сервера (по умолчанию 8080)

### Frontend
Настройки загружаются из файла `config.ini` или переменных окружения:
- `API_BASE_URL` - базовый URL API сервера
- `API_VERSION` - версия API

## API Endpoints

### Аутентификация
- `POST /api/v1/auth/login` - вход в систему
- `POST /api/v1/auth/register` - регистрация

### Тикеты
- `GET /api/v1/tickets` - список тикетов
- `POST /api/v1/tickets` - создание тикета
- `GET /api/v1/tickets/:id` - получение тикета по ID
- `PATCH /api/v1/tickets/:id` - обновление тикета
- `DELETE /api/v1/tickets/:id` - удаление тикета

### Справочники
- `GET /api/v1/departments` - список департаментов
- `GET /api/v1/ticket_statuses` - список статусов
- `GET /api/v1/ticket_priorities` - список приоритетов

## База данных

Система использует PostgreSQL. Миграции находятся в папке `migrations/`.

Для инициализации базы данных:
```bash
# Установите golang-migrate
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Примените миграции
migrate -path migrations -database "postgres://user:password@localhost:5432/ticket_system?sslmode=disable" up
```

## Безопасность

- ✅ Пароли хешируются с использованием bcrypt
- ✅ JWT токены для аутентификации
- ✅ Валидация входных данных
- ✅ Защита от SQL инъекций (используется GORM)
- ✅ Логирование ошибок и запросов

## Логирование

Backend логирует:
- Все HTTP запросы с временем выполнения
- Ошибки с детальной информацией
- Ошибки базы данных
- Ошибки валидации

## Отладка

### Backend
Логи выводятся в stdout. Для включения детального логирования:
```bash
export GIN_MODE=debug
```

### Frontend
Логи выводятся в консоль. Для включения детального логирования:
```bash
export QT_LOGGING_RULES="*.debug=true"
```

## Известные проблемы

1. **Решено**: Крэш при создании тикета - исправлено добавлением отложенной загрузки ComboBox
2. **Решено**: Проблемы с JSON binding - исправлено добавлением правильных JSON тегов
3. **Решено**: Хардкод URL - исправлено созданием системы конфигурации

## Требования

- Go 1.19+
- Qt 6.0+
- PostgreSQL 12+
- CMake 3.16+
