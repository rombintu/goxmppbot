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

### Mongo (dev)
```
docker run --name mongodb \ 
    -e MONGO_INITDB_ROOT_USERNAME=<username> \ 
    -e MONGO_INITDB_ROOT_PASSWORD=<passowrd> \ 
    -e MONGO_INITDB_DATABASE=backend \ 
    -v <some_path>/mongo:/data/db \ 
    -d -p 27017:27017 mongo
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
CONNECTION="mongodb://root:root@192.168.0.2:27017" # Значение из ENV в приоритете
DBNAME="backend"
```