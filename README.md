# Менеджер паролей gophkeeper
( a pet project by Painkiller675)

## Настройка и запуск серверной части

Перед запуском сервера необходимо создать конфигурационный файл с настройками
или задать переменные окружения. 
Пример конфигурационного файла:

```
# server-config.yaml

grpc:
  address: 127.0.0.1:8081
db:
  url: postgres://postgres:1234@localhost:5432/keeperdb?sslmode=disable
auth:
  key: xiuw1bi4r98vd1(&*7
  expiration_time: 24h
hasher:
  key: jc7YSHpH287)(*2bSw
```
godotenv
Пример настройки сервера через переменные окружения:

```
# .env

GRPC_ADDRESS=127.0.0.1:8081
DB_URL=postgres://postgres:1234@localhost:5432/keeperdb?sslmode=disable
AUTH_KEY=xiuw1bi4r98vd1(&*7
AUTH_EXPIRATION_TIME=24h
HASHER_KEY=jc7YSHpH287)(*2bSw
```

После настройки запуск серверной части осуществляется командной:

```
./gophkeeper server
```
