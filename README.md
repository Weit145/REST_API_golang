# REST API Golang (Pet Project)

Простой REST API на Go для работы с заказами (Orders).  
Поддерживает CRUD операции: Create, Read, Update, Delete.  

> **Примечание:** это просто **pet project** для практики.


## CI/CD

Проект настроен с использованием **GitHub Actions** для автоматического тестирования:  

- Все изменения в ветке `master` проверяются через CI.
- Перед слиянием в защищённые ветки необходимо создавать **Pull Request**.
- В Workflow выполняются следующие шаги:
  1. Установка Go (`1.25.5`) и зависимостей.
  2. Кэширование модулей Go для ускорения сборки.
  3. Сборка проекта.
  4. Запуск **unit-тестов** для всех CRUD-обработчиков.
  5. Запуск интеграционных тестов из папки `tests`.

> Благодаря этому, изменения с ошибками не попадут в основную ветку.


## Запуск

1. Клонировать репозиторий и перейти в папку:

```bash
git clone <https://github.com/Weit145/REST_API_golang>
cd REST_API_golang
````

2. Установить зависимости и запустить сервер:

```bash
go mod tidy
go run cmd/main.go
```

Сервер будет доступен на: `http://0.0.0.0:8080`


## Аутентификация

Для POST, PUT, DELETE используется **Basic Auth**:

* Логин: `Weit`
* Пароль: `123`

GET запросы открыты для всех.



## Примеры запросов

### Создать заказ

```bash
curl -X POST http://127.0.0.1:8080/order \
  -H "Content-Type: application/json" \
  -u Weit:123 \
  -d '{"order_name":"TestOrder","price":123.45}'
```

### Получить заказ

```bash
curl -X GET http://127.0.0.1:8080/order/TestOrder
```

### Обновить заказ

```bash
curl -X PUT http://127.0.0.1:8080/order \
  -H "Content-Type: application/json" \
  -u Weit:123 \
  -d '{"order_name":"TestOrder","price":200.00}'
```

### Удалить заказ

```bash
curl -X DELETE http://127.0.0.1:8080/order \
  -H "Content-Type: application/json" \
  -u Weit:123 \
  -d '{"order_name":"TestOrder"}'
```


## Тесты

Для запуска тестов:

```bash
go test ./tests -v
go test ./internal/http-server/handler/order/read -v
go test ./internal/http-server/handler/order/create -v
go test ./internal/http-server/handler/order/delete -v
go test ./internal/http-server/handler/order/update -v
```


## Зависимости

* **Go 1.25+** — версия Go, используемая в проекте
* **[Chi Router](https://github.com/go-chi/chi)** — маршрутизация HTTP и middleware
* **[SQLite](https://www.sqlite.org/)** — встроенная база данных для хранения заказов
* **[Testify](https://github.com/stretchr/testify)** — библиотека для unit-тестов и моков
* **[httpexpect](https://github.com/gavv/httpexpect)** — тестирование REST API
* **[Validator v10](https://github.com/go-playground/validator)** — валидация входящих данных
* **[Render](https://github.com/go-chi/render)** — удобная работа с JSON ответами и HTTP статусами

