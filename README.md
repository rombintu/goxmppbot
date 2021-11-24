# Bot (XMPP)

## Dependencies
```
make, git, sqlite3 (dev)
Docker or (golang => 1.16)
```
## Deploy (Local)
```
cp config.toml.bak config.toml
edit config.toml
make build
./main
```
## Deploy (Docker)
```
docker-compose up -d
docker-compose ps
docker-compose logs -f bot
```

## Configuration (Example)
```toml
[DEFAULT]
HOST = "Хост учетной записи" # Example: "xmpp.mail.ru"
LOGIN = "Логин учетной записи"
PASSWORD = "Пароль учетной записи" 
DEBUGLEVEL = "Уровень дебаг-лога" # [debug, info]
DEBUG = true или false # Включить или выключить дебаг-лог
REFRESH_SECRET = "1234"

[SUPPORT]
HOST = "Хост почты" # Example: "mail.ru"
PORT = "Порт" # Example: "smtp"
LOGIN = "Логин аккаунта поддержки, без хоста, от кого идут письма" # Example: "user" 
PASSWORD = "Пароль аккаунта поддержки"
SUPPORTEMAIL = "Полный логин поддержки, куда идут письма" # Example: "user@mail.ru"

[CONTACTS]
URL = "Адрес до контактов" # Example: "https://mail.ru/contacts"

[BACKENDCONF]
CONNECTION= "Путь до локальной базы данных" # Example: "sqlite.db"
DEV=true # Рабоать в режиме разработки

[LINKS] # Перечисление сервисов
service_name = "НАЗВАНИЕ_СЕРВИСА - ссылка"
```