# Применение In-Memory СУБД

В ходе работы для повышения отказоустойчивости был подключен In-memory СУБД Redis и перенесена бизнес логика в UDF функцию на LUA, что позволило сократить сетевые издержки и повысить производительность и скорость выполнения запроса.

## Развертывание

Открываем терминал на нашем сервере и в cli выполняем команды согласно инструкции.

1) Скачиваем 3 репозитория (otus, ms_dialog, event-getaway)

```
git clone https://github.com/proweb-zone/event-getaway.git

git clone https://github.com/proweb-zone/otus_social_network.git

git clone https://github.com/proweb-zone/ms_dialog.git
```
2) Переходим в проект event-getaway и запускаем docker контейнер

```
cd /my_path/event-getaway/
docker compose up -d
```
2) Переходим в проект otus_social_network и запускаем docker контейнер и делаем миграцию БД

```
#  переходим в проект
cd /my_path/event-getaway/

# поднимаем docker
docker compose up -d

# проводим миграцию БД
docker exec -it ms-dialog make migration-up
```
3) Переходим в проект ms_dialog и запускаем docker контейнер и делаем миграцию БД

```
#  переходим в проект
cd /my_path/ms_dialog/

# поднимаем docker
docker compose up -d

# проводим миграцию БД
docker exec -it app_socnet make migration-up
```

4) Перенести бизнес логику на UDF (LUA)

Я написал LUA скрипт для redis который получает все сообщения друзей.
Для регистрации LUA скрипта необходимо в контейнере redis выполнить следующую команду

```
docker exec -it redis-container bash
cd /
redis-cli -x FUNCTION LOAD < dialog_functions.lua
```

### Нагрузочное тестирование до подключения In-memory СУБД Redis
p 95 - выдерживает 4тыс. RPC
p 95 http_req_duration p(95)<1000 - 4тыс. RPC

### Нагрузочное тестирование после подключения In-memory СУБД Redis
p 95 - выдерживает 4тыс. RPC
p 95 http_req_duration p(95)<1000 - 4тыс. RPC
