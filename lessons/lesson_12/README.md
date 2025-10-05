# Мониторинг

### Цель - В результате выполнения ДЗ вы организуете мониторинг своего сервиса чатов

### Отчет о проделанной работе

Для выполнения домашнего задания я поднял в docker сервисы zabbix, grafana, prometeus.
  zabbix - служит для сбора технических метрик сервера,
  grafana - для сбора бизнес метрик (по принципу RED),
  prometeus - промежуточная бд для организации и хранения данных (для удобства предоставления графане)

# Сбор бизнес метрик  RED
1) Для сбора бизнес метрик я добавил в микросервис ms_dialog дополнительный middleware для всех запросов, который собирет всю информацию и сохраняет ее по пути /metrics.
Ссылка на скриншот:
https://github.com/proweb-zone/ms_dialog/blob/main/lessons/lesson_12/scrin_metrics.png
2) Добавил конфигурационный файл в prometheus prometheus.yml

```
global:
  scrape_interval: 15s
  evaluation_interval: 15s

scrape_configs:
  - job_name: 'ms-dialog'
    static_configs:
      - targets: ['ms-dialog:3002'] # Важно: используем имя сервиса из docker-compose
        labels:
          service: 'ms-dialog'

  - job_name: 'prometheus'
    static_configs:
      - targets: ['localhost:9090']
```
3) В grafana тоже добавил конфиги для сбора данных из prometeus

```
apiVersion: 1

datasources:
  - name: Prometheus
    type: prometheus
    access: proxy
    url: http://prometheus:9090 # Имя сервиса из docker-compose
    isDefault: true
    version: 1
    editable: false
```

4) Настроил dashboads в grafana со следующими настройками

```
# Сумма всех звпросов
sum(http_requests_total{job="ms-dialog"})


# процент отказа
sum(http_requests_total{job="ms-dialog", status_code=~"500"})/sum(http_requests_total{job="ms-dialog"})*100

#Гистограмма показывает время выполнения в 95 процентиле
histogram_quantile(0.95,
  sum(rate(http_request_duration_seconds_bucket{job="ms-dialog"})) by (le, path)
)
```

По итогу у меня получилось настроить метрики для сбора бизнес метрик по принципу RED, которые показывают колличество успешных запросов, процент отказа, и время выполнения

Ссылка на скриншот из grafana:
https://github.com/proweb-zone/ms_dialog/blob/main/lessons/lesson_12/grafana_ms_dialog.png

### Сбор технических метрик

Поскольку весь проект находиться в docker контейнерах и postgres и redis и микросервис - для сбора метрик я добавил контейнер zabbix-agent-dialog (Zabbix Agent для сбора метрик со всех сервисов).
Далее написал shell скрипт для подключения и передачи данных в zabbix каждого контейнера.
Ссылка на shell скрипты: https://github.com/proweb-zone/ms_dialog/tree/main/zabbix/agent/scripts

В zabbix agent прописал конфигурационный файл.
```
Server=zabbix-server
ServerActive=zabbix-server
Hostname=dialog-service-cluster
HostMetadata=dialog,postgres,redis,golang,docker

# Базовые настройки
ListenPort=10050
ListenIP=0.0.0.0
Timeout=30

# Только один Include для UserParameters
Include=/etc/zabbix/zabbix_agent2.d/user_parameters.conf
```

Далее для настройки dasboards в zabbix необходимо настроить hosts, templates, items.
После чего я создал для каждого сервиса (postgres, redis и т.д.) отдельный dasboards и указал нужные items для отображения.

По результату у меня получились метрики которые предоставляют технические метрики по всем сервисам.
Скрины на метрики:
https://github.com/proweb-zone/ms_dialog/blob/main/lessons/lesson_12/zabbix_postgres_metrics.png
https://github.com/proweb-zone/ms_dialog/blob/main/lessons/lesson_12/zabbix_redis_metrics.png
