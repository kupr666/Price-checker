## Price-checker (for sites with static html)
REST API сервис, который отслеживает цены товаров, а также уведомляет вас в telegram,
когда цена стала желаемой.
Инструкция по настройке и запуску ниже. Для быстрой проверки можно использовать сайт:
https://future-phone.ru/ , который работает по умолчанию.

### Stack
- Golang 1.26.1
- Docker & Docker compose
- postgreSQL 18.3 
- Redis 8.6 (для кеширования цен)
- Telegram BOT API

### Launch
1. Клонировать репозиторий
2. Создать .env из .env.example
3. Запустить с помощью make service-deploy. (Можно запустить в контейнерах только redis
и postgreSQL, а приложение запустить локально: make service-dev-db, make service-dev-redis, make service-run)

### Fast check (tests)

### Add your own website with static HTML code
Чтобы добавить поддержку нужного вебсайта (со статическим HTML) необходимо
вручную зарегистрировать домен и соответствующие HTML теги в исходном коде.
Открыть файл по пути /price_checker/internal/features/price_tracker/scraper/scraper.go
В конструкторе (NewGoQueryScraper) в мапу sellectors добавить нужный вам домен, а также HTML тег.

### API
- POST /items - созадать ссылку на товар, цена которого будет отслеживаться
- DELETE /items/{id} - удалить ссылку на определённый товар
- GET /items - получить список всех отслеживаемых товаров

### Migrations
Миграции запускаются автоматически при старте приложения
