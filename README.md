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
## Deploy (Docker) - v0.2.0-unstable
```
docker-compose up -d
docker-compose ps
docker-compose logs -f bot
```

### PostgresSQL (dev)
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
# Хост Jabber'a 
Host = "mail.ru"

# Логин и пароль от пользователя Jabber
Login = "user@mail.ru"
Password = ""

# Дебаг уровень для логирования [debug, info]
DebugLevel = "debug"
DebugOn = true
LogFile = "/var/log/mvdbot/mvdbot.log"

# Через сколько минут обновлять бота
UpdateChunk = 30

# Секрет для админских команд в боте
RefreshSecret = ""

# Включенные плагины: ["zabbix"]
Plugins = ["zabbix"]


[SUPPORT]
# Хост и порт почты (поддержки)
Host = "mail.ru"
Port = "smtp"

# Логин без хоста и пароль от пользователя поддержки
# С этого пользователя идут письма
LoginWithoutHost = "user" # FROM
Password = ""

# Куда идут письма
SupportEmail = "user@mail.ru"

[CONTACTS]
# Адрес до API контактов
URL = ""

[DBCONF]
# Конфигурация для подключения к БД
Master = "host=localhost port=5432 user=user password=password dbname=backend sslmode=disable"
Slave = "host=localhost port=5432 user=user password=password dbname=backend sslmode=disable"

# Режим разработки [true, false]
Dev = true

# Мульти-режим для БД [true, false]
Multi = false


# Настройка плагинов
[ZABBIX]
Host = ""
User = ""
Password = ""
```

# Tags. Fix
v0.4.5-stable
* Изменена система логирования
* Зачатки CI\CD. Можно создать .deb пакет локально
* Разработан CLI для администрирования

v0.4.4 
* Исправления в плагинах
* Исправления в конфигурации
* TODO - Обновление токена

v0.4.3-unstable
* Разработана библиотека для работы с Zabbix
* Добавлена возможность подключение плагинов
* Исправления

v0.4.2
* Поддержка Zabbix
* Исправления в драйверах БД

v0.4.1
* Мелкие исправления

v0.4.0-unstable
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