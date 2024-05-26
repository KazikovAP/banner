[![Go](https://img.shields.io/badge/-Go-464646?style=flat-square&logo=Go)](https://go.dev/)
[![PostgreSQL](https://img.shields.io/badge/-PostgreSQL-464646?style=flat-square&logo=PostgreSQL)](https://www.postgresql.org/)
[![Redis](https://img.shields.io/badge/-Redis-464646?style=flat-square&logo=Redis)](https://developer.redis.com/)
[![docker](https://img.shields.io/badge/-Docker-464646?style=flat-square&logo=docker)](https://www.docker.com/)

# banner
# Сервис, показывающий пользователем баннеры, в зависимости от требуемой фичи и тега пользователя, а также управляемый баннерами и связанными с ними тегами и фичами.

---
## Технологии
* Golang 1.22.0
* PostgreSQL
* Redis
* REST API
* Docker
* Postman
* JWT

---
## Взаимодейастиве с сервисом

### Экспорировать путь до конфига
`export CONFIG_PATH="<path>\banner\config\config.yaml"` 

### Запуск приложения локально
`go run cmd/banner/main.go`

### Запуск докер контейнера с Postgres
`docker compose -p banner -f ./build/docker-compose.yaml up -d`

### Остановка и удаление докер контейнер с Postgres
`docker compose -p banner -f ./build/docker-compose.yaml down`

## Примеры запросов
### Authorization
**Регистрация пользователя:** POST запрос `http://localhost:8080/auth/sign-up`:

**Request:**
```JSON
{
    "name": "admin",
    "username": "admin",
    "password": "123456"
}
```
**Response:**
```JSON
{
    "id": 1
}
```

**Авторизация пользователя:** POST запрос `http://localhost:8080/auth/login`:

**Request:**
```JSON
{
    "username": "admin",
    "password": "123456",
}
```
**Response:**
```JSON
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTMwNzM4OTUsImlhdCI6MTcxMzAzMDY5NSwidXNlcl9pZCI6Nn0.cIhKXB6nTFLlfZGt5z3cR6yQPu1aKbQmW0DoEcaT5zw"
}
```

**Выход пользователя:** DELETE запрос `http://localhost:8080/auth/logout`:

**Request:**
```JSON
{
    
}
```
**Response:**
```JSON
{
    "status": 200
}
```

### Banners API
**Получение всех баннеров:** GET запрос `http://localhost:8080/api/banners`:

**Request:**
```JSON
{
    
}
```
**Response:**
```JSON
{
    "data": [
        {
            "id": 1,
            "isActive": 1,
            "featureId": 2,
            "tagId_1": 1,
            "tagId_2": 2,
            "tagId_3": 5
        },
        {
            "id": 2,
            "isActive": 2,
            "featureId": 4,
            "tagId_1": 5,
            "tagId_2": 1,
            "tagId_3": 5
        },
        {
            "id": 3,
            "isActive": 1,
            "featureId": 3,
            "tagId_1": 2,
            "tagId_2": 1,
            "tagId_3": 2
        },
        {
            "id": 4,
            "isActive": 2,
            "featureId": 3,
            "tagId_1": 2,
            "tagId_2": 4,
            "tagId_3": 2
        },
        {
            "id": 5,
            "isActive": 1,
            "featureId": 5,
            "tagId_1": 8,
            "tagId_2": 7,
            "tagId_3": 2
        }
    ]
}
```

**Получение баннеров по ID:** GET запрос `http://localhost:8080/api/banner/feature=5`:

**Request:**
```JSON
{
    
}
```
**Response:**
```JSON
{
    "data": [
        {
            "id": 5,
            "isActive": 1,
            "featureId": 5,
            "tagId_1": 8,
            "tagId_2": 7,
            "tagId_3": 2
        }
    ]
}
```

**Получение баннеров по ID тега:** GET запрос `http://localhost:8080/api/banner/tag=1`:

**Request:**
```JSON
{
    
}
```
**Response:**
```JSON
{
    "data": [
        {
            "id": 1,
            "isActive": 1,
            "featureId": 2,
            "tagId_1": 1,
            "tagId_2": 2,
            "tagId_3": 5
        },
        {
            "id": 2,
            "isActive": 2,
            "featureId": 4,
            "tagId_1": 5,
            "tagId_2": 1,
            "tagId_3": 5
        },
        {
            "id": 3,
            "isActive": 1,
            "featureId": 3,
            "tagId_1": 2,
            "tagId_2": 1,
            "tagId_3": 2
        }
    ]
}
```

**Создание баннера:** POST запрос `http://localhost:8080/api/banner`:

**Request:**
```JSON
{
    "isActive": 1,
    "featureId": 2,
    "tagId_1": 1,
    "tagId_2": 2,
    "tagId_3": 5
}
```
**Response:**
```JSON
{
    "id": 4
}
```

**Обновление баннера:** PATCH запрос `http://localhost:8080/api/updateBanner/5`:

**Request:**
```JSON
{
    "isActive": 1,
    "featureId": 15,
    "tagId_1": 16,
    "tagId_2": 17,
    "tagId_3": 18
}
```
**Response:**
```JSON
{
   
}
```

**Удаление баннера:** DELETE запрос `http://localhost:8080/api/deleteByFeature/4`:

**Request:**
```JSON
{
    
}
```
**Response:**
```JSON
{
   "status": 200
}
```

---
## Разработал:
[Aleksey Kazikov](https://github.com/KazikovAP)
---
## Лицензия:
[MIT](https://opensource.org/licenses/MIT)
