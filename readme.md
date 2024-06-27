### Микросервис в этом репозитории

В данном репозитории находится код, который предназначен для работы *offline-магазина* в **pet-проекте**.<br>
<br>Используемое при разработке ПО и технологии/форумы:<br>
![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
![MySQL](https://img.shields.io/badge/mysql-%2300f.svg?style=for-the-badge&logo=mysql&logoColor=white)
![Stack Overflow](https://img.shields.io/badge/-Stackoverflow-FE7A16?style=for-the-badge&logo=stack-overflow&logoColor=white)
![GoLand](https://img.shields.io/badge/GoLand-0f0f0f?&style=for-the-badge&logo=goland&logoColor=white)
![Notepad++](https://img.shields.io/badge/Notepad++-90E59A.svg?style=for-the-badge&logo=notepad%2b%2b&logoColor=black)
![Apache Kafka](https://img.shields.io/badge/Apache%20Kafka-000?style=for-the-badge&logo=apachekafka)
![Obsidian](https://img.shields.io/badge/Obsidian-%23483699.svg?style=for-the-badge&logo=obsidian&logoColor=white)
![Postman](https://img.shields.io/badge/Postman-FF6C37?style=for-the-badge&logo=postman&logoColor=white)
![Markdown](https://img.shields.io/badge/markdown-%23000000.svg?style=for-the-badge&logo=markdown&logoColor=white)
![Ubuntu](https://img.shields.io/badge/Ubuntu-E95420?style=for-the-badge&logo=ubuntu&logoColor=white)
![JWT](https://img.shields.io/badge/JWT-black?style=for-the-badge&logo=JSON%20web%20tokens)
![Prometheus](https://img.shields.io/badge/Prometheus-E6522C?style=for-the-badge&logo=Prometheus&logoColor=white)

#### Содержание

[1. Описание проекта](#описание)

[2. Ограничения, принятые в проекте](#ограничения)

[3. Доступные команды](#команды)

[4. Конфигурация](#конфигурация)

[5. JWT](#jwt)

[6. Для чего это всё написано?](#длячего)

#### Описание

Проект будет имитировать бэкэнд сети магазинов по продаже часов ⌚. Описание требований:

- продажи будут вестись как в обычных *offline-магазинах*, так и через интернет магазин
- возможна отправка товара не только со складов компании, но и из *offline-магазинов*
- доступна бронь товара через интернет для самовывоза из конкретного *offline-магазина*

#### Ограничения

+ номер заказа - положительное число
+ номер заказа через интернет должен быть больше десяти. Номера меньше десяти используются для касс, которые отмечают
  пробитые, но еще не оплаченные товары, как заказанные. Это необходимо для того, чтобы исключить бронирование товаров,
  находящихся в процессе продажи
+ Дефекты товаров или упаковки, влияющие на цену, шифруются в артикуле товара (а это значит, что товар без дефектов и с
  дефектом имеют разные артикулы). Дефекты кодируются следующим образом - за основу берется артикул неповрежденного
  товара, после ставится точка, а далее идут четыре цифры:
    1. 0 - корпус без повреждений, 1 - корпус имеет легкие царапины, 2 - корпус имеет сильные царапины
    2. 0 - дисплей/стекло без повреждений, 1 - дисплей/стекло имеет легкие царапины, 2 - дисплей/стекло имеет сильные
       царапины. (Для ремешков всегда указывается 0, так как у них нет дисплея)
    3. 0 - упаковка/коробка не вскрывалась, 1 - упаковка/коробка вскрывалась
    4. 0 - упаковка/коробка без повреждений, 1 - упаковка/коробка повреждена

#### Команды

В проекте содержится Makefile, содержащий полезные в процессе разработки и развертывания команды:

+ **make help** - выводит справку по доступным опциям команды make
+ **make test** - запускает тесты
+ **make cover** - выводит покрытие кода тестами в браузере по умолчанию

#### Конфигурация

Конфигурация приложения сохраняется в YAML-файлах, имеющих следующую структуру:

```yaml
# в каком окружении запускается программа. Есть три варианта - "local", для обычной разработки, "debug" - для  
# отладки/разработки с проверкой прав доступа и "production" - для запуска на боевом сервере
env: "local"
# название экземпляра запущенного приложения. Служит уникальным идентификатором приложения в системе
instance: "instance1"
# нужно ли использовать брокер сообщений kafka
use_kafka: true
# раздел настройки http
http_server:
  # адрес и порт http-сервера
  address: "localhost:8091"
  # таймаут чтения
  read_timeout: 5s
  # таймаут записи
  write_timeout: 10s
  # таймаут простоя  
  idle_timeout: 60s
  # таймаут на завершение работы http-сервера при gracefully shutdown
  shutdown_timeout: 15s
# раздел настройки хранилища
storage:
  # логин базы данных
  database_login: "login"
  # пароль базы данных
  database_password: "password"
  # адрес базы данных
  database_address: "localhost"
  # максимально доступное количество открытых соединений базы данных
  database_max_open_connections: 10
  # имя базы данных
  database_name: "db_name"
  # таймаут запроса
  query_timeout: 5s
  # порт для отображения таблиц БД
  database_viewer_port: 9123
# раздел настройки безопасности
secure:
  # секретная подпись для валидации JWT
  secure_signature: "secure signature"
  # адрес сервера, выдающего JWT токены и имеющего право знать секретную подпись
  secure_server: "localhost:8095"
  # таймаут запроса к сервису безопасности
  secure_request_timeout: 2s
  # количество попыток получения токена с сервису безопасности
  secure_attempts: 3
  # протокол для общения с сервисом безопасности (http или https)
  secure_protocol: "http"
  # логин в сервисе безопасности
  secure_username: "store1"
  # пароль в системе безопасности
  secure_password: "pwd"
  # нужно ли хранить в файле кешированные разрешения с их номерами
   secure_use_permissions_file_cache: true
  # путь и название файла кеша разрешений и их номеров
  secure_permissions_file: "./cache/secure_permissions.json"
# раздел настройки kafka
kafka:
  # адреса брокеров kafka
  kafka_brokers: [ "localhost:9092" ]
  # название топика с обновлениями цены
  kafka_topic_update_price: "store.update-price"
# раздел настройки Prometheus 
prometheus:
  # на каком порту собирать метрики. Если не задан, то по умолчанию порт 9323
  prometheus_port: "9099"
  # url для сбора метрик Prometheus. Если не задан, то по умолчанию используется /metrics 
  prometheus_metrics_url: "/metrics"
```

Есть возможность переопределять значения из конфигурационных файлов переменными окружения. Соответствие опций из
конфигурационного файла переменным окружения представлено в таблице ниже:

| В файле конфигурации              | Переменная окружения              |
|-----------------------------------|-----------------------------------|
| instance                          | INSTANCE                          |
| env                               | ENV                               |
| secure_signature                  | SECURE_SIGNATURE                  |
| secure_server                     | SECURE_SERVER                     |
| secure_request_timeout            | SECURE_REQUEST_TIMEOUT            |
| secure_attempts                   | SECURE_ATTEMPTS                   |
| secure_protocol                   | SECURE_PROTOCOL                   |
| secure_username                   | SECURE_USERNAME                   |
| secure_password                   | SECURE_PASSWORD                   |
| secure_use_permissions_file_cache | SECURE_USE_PERMISSIONS_FILE_CACHE |
| secure_permissions_file           | SECURE_PERMISSIONS_FILE           |
| address                           | ADDRESS                           |
| read_timeout                      | READ_TIMEOUT                      |
| write_timeout                     | WRITE_TIMEOUT                     |
| idle_timeout                      | IDLE_TIMEOUT                      |
| shutdown_timeout                  | SHUTDOWN_TIMEOUT                  |
| database_login                    | DATABASE_LOGIN                    |
| database_password                 | DATABASE_PASSWORD                 |
| database_address                  | DATABASE_ADDRESS                  |
| database_name                     | DATABASE_NAME                     |
| database_max_open_connections     | DATABASE_MAX_OPEN_CONNECTIONS     |
| query_timeout                     | QUERY_TIMEOUT                     |
| database_viewer_port              | DATABASE_VIEWER_PORT              |
| kafka_brokers                     | KAFKA_BROKERS                     |
| kafka_topic_update_price          | KAFKA_TOPIC_UPDATE_PRICE          |
| prometheus_port                   | PROMETHEUS_PORT                   |
| prometheus_metrics_url            | PROMETHEUS_METRICS_URL            |

Путь к файлу конфигурации можно указывать по ключу *config* при запуске приложения или в переменной окружения
*STORE_CONFIG_PATH*. При отсутствии конфигурации приложение завершится с ошибкой.

#### JWT

Если приложение запущено не с конфигурацией локального окружения, то при HTTP-запросах выполняется middleware,
проверяющий корректность JWT-токена, содержащегося в заголовке Authorization. Префикс токена - *"Bearer "*. Алгоритм -
*HS256*. В полезной нагрузке токена должны быть переданы номера разрешений по ключу 'perm'.

#### ДляЧего?

В данном репозитории содержится код, являющийся частью моего **pet-проекта**, цель которого - изучение языка Golang,
микросервисной архитектуры, взаимодействия микросервисов между собой.