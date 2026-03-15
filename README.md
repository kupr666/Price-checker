## Price-checker (for sites with static html)
Чтобы добавить поддержку нужного вебсайта (со статическим html) необходимо
вручную зарегистрировать домен и соответствующие HTML теги в исходном коде. 
Инструкция по настройке и запуску ниже.

### stack
- golang go1.26.1
- Docker
- postgreSQL 18.3
- Redis 8.6

### Add your own website with static HTML code
Открыть файл по пути /price_checker/internal/features/price_tracker/scraper/scraper.go
В конструкторе (NewGoQueryScraper) в мапу sellectors добавить нужный вам домен, а также html тег.

### Launch
1. Клонировать репозиторий
2. Создать .env из .env.example
3. Запустить с помощью

### API
- POST /items - созадать ссылку на товар, цена которого будет отслеживаться
- DELETE /items/{id} - удалить ссылку на определённый товар
- GET /items - получить список всех отслеживаемых товаров

### Migrations
