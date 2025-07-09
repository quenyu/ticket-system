## Swagger/OpenAPI документация

Документация API генерируется автоматически с помощью [swaggo/swag](https://github.com/swaggo/swag) и доступна по адресу:

    http://localhost:8080/swagger/index.html

### Как сгенерировать/обновить спецификацию

1. Установите swag (один раз):

    go install github.com/swaggo/swag/cmd/swag@latest

2. Сгенерируйте документацию:

    swag init --generalInfo cmd/api/main.go --output docs

3. Запустите backend:

    go run cmd/api/main.go

Swagger UI будет доступен на /swagger/index.html

### Как добавить описание к эндпоинтам

Добавляйте GoDoc-style комментарии над хендлерами, например:

```go
// @Summary Логин пользователя
// @Description Авторизация по логину и паролю
// @Tags auth
// @Accept json
// @Produce json
// @Param input body model.LoginRequest true "Данные для входа"
// @Success 200 {object} model.LoginResponse
// @Failure 401 {object} model.ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) { ... }
```

После изменений не забудьте снова выполнить `swag init`. 