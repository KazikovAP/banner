[![Go](https://img.shields.io/badge/-Go-464646?style=flat-square&logo=Go)](https://go.dev/)
[![PostgreSQL](https://img.shields.io/badge/-PostgreSQL-464646?style=flat-square&logo=PostgreSQL)](https://www.postgresql.org/)
[![Redis](https://img.shields.io/badge/-Redis-464646?style=flat-square&logo=Redis)](https://www.developer.redis.com/)
[![docker](https://img.shields.io/badge/-Docker-464646?style=flat-square&logo=docker)](https://www.docker.com/)

# banner
# Сервис, показывающий пользователем баннеры, в зависимости от требуемой фичи и тега пользователя, а также управляемый баннерами и связанными с ними тегами и фичами.

---
## Технологии
* Golang 1.22.0
* PostgreSQL
* Cash Redis
* REST API
* Docker
* Postman

---
## Взаимодейастиве с сервисом

### Экспорировать путь до конфига
`export CONFIG_PATH="<path>\banner\config\config.yaml"` 

### Запуск приложения в хранилище in-memory
`go run cmd/banner/main.go`

---
## Разработал:
[Aleksey Kazikov](https://github.com/KazikovAP)
---
## Лицензия:
[MIT](https://opensource.org/licenses/MIT)
