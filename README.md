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
docker run -d \
    --name postgres \
    -e POSTGRES_DB=DATABASE \
    -e POSTGRES_USER=USER \
    -e POSTGRES_PASSWORD=PASSWORD \
    -e PGDATA=/var/lib/postgresql/data/pgdata \
    -v /custom/mount:/var/lib/postgresql/data \
    postgres:14-alpine
```

## Configuration (Example)
```toml
[DEFAULT]
Host = "Хост учетной записи" # Example: "xmpp.mail.ru"
Login = "Логин учетной записи"
Password = "Пароль учетной записи" 
DebugLevel = "Уровень дебаг-лога" # [debug, info]
DebugOn = true или false # Включить или выключить дебаг-лог
RefreshSecret = "1234"
ZabbixOn = true

[SUPPORT]
Host = "Хост почты" # Example: "mail.ru"
Port = "Порт" # Example: "smtp"
LoginWithoutHost = "Логин аккаунта поддержки, без хоста, от кого идут письма" # Example: "user" 
Password = "Пароль аккаунта поддержки"
SupportEmail = "Полный логин поддержки, куда идут письма" # Example: "user@mail.ru"

[CONTACTS]
URL = "Адрес до контактов" # Example: "https://mail.ru/contacts"

# Значение из ENV в приоритете
[DBCONF]
Master = "host=localhost port=5432 user=user password=password dbname=backend sslmode=disable"
Slave = "host=localhost port=5432 user=user password=password dbname=backend sslmode=disable"
Dev = true
Multi = false

[ZABBIX]
Host = ""
User = ""
Password = ""
```

# Tags. Fix
v0.4.3-unstable
* Написана своя библиотека для работы с Zabbix
* Возможность подключение плагинов
* Исправления

v0.4.2
* Поддержка Zabbix
* Исправления в драйверах БД

v0.4.1
* Мелкие исправления

v0.4.0
* Изменена логика запоминия предыдущих команд пользователя

v0.3.2-stable
* Поддержка xslx 

v0.3.1-stable
* Удалена тема сообщений
* Настроена универсальность под разные клиенты
* TODO - Разобраться с отправкой\обработкой файлов

v0.3.0-stable
* Исправлена логика таймаута, по которому обновляется соединение (10 мин)
* Исправлены ошибки, возникающие при простаивании бота
* Запоминание команд не пересекаются с данными сервисов и крутятся в sqlite3(TODO:memory)

v0.2.1-stable
* Исправлены ошибки

v0.2.0
* Добавлена возможность мультиподключения к базам данным