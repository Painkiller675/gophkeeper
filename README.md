# Менеджер паролей gophkeeper
( a pet project by Painkiller675)

## Настройка и запуск серверной части

Запуск сервера осуществляется с использованием конфигурационного файла или переменных окружения.

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
## Настройка и запуск клиентской части 

Запуск клиента осуществляется с использованием конфигурационного файла или переменных окружения.

Пример конфигурационного файла клиента:

```
# client-config.yaml

grpc:
  address: 127.0.0.1:8081
encryption:
  key: WYJcWgkItShq513L21E1CFuz6uQWDy3r
```

Пример настройки клиента c использованием переменных окружения:

```
# .env

GRPC_ADDRESS=127.0.0.1:8081
ENCRYPTION_KEY=WYJcWgkItShq513L21E1CFuz6uQWDy3r
```