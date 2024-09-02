# SarkorTelekom
1) Скачать репозиторий
2) В консоли перейти в директорию в которую вы скачали репозиторий
3) Запустить БД - postgres (я использовал docker container, дальше приложу команды для запуска контейнера
- docker run --name=pg-test -e POSTGRES_PASSWORD='secret' -p 5444:5432 -d --rm postgres (в конфиге указан этот пароль и порт)
- docker exec -ti pg-test psql -U postgres
- CREATE DATABASE sarkortest;
- \c sarkortest
- ALTER USER postgres WITH PASSWORD 'secret'; (на всякий случай поменять пароль повторно на secret, так как иначе может не сработать файлы миграции)
- (если создали докер контайнер используя эти настройки этапы 4-5 можно пропустить))
4) Создать отдельную Базу - sarkortest
5) Настроить данные для подключения к БД в файле configs/config (верхнее значение порта используется для самого приложения)
6) Запустить файлы миграции командой
migrate -database "postgres://username:password@localhost:port/sarkortest?sslmode=disable" -path ./schema up
7) Запустить cmd/main.go файл
8) Проверку API осуществлял через Postman
- POST "/products" - добавляет запись в таблицу, возвращает id записи
- GET "/product" - указываете id, возвращает информацию о продукте
- PUT "/product" - указываете id, обновляет содержимое
- DELETE "/product" - указываете id, удаляет указанный продукт
- GET "/products" -  возвращает список всего содержимого таблицы в БД 


Использовал в проекте - Gin - github.com/gin-gonic/gin, 
  sqlx - github.com/jmoiron/sqlx, 
  logrus - github.com/sirupsen/logrus, 
  viper - github.com/spf13/viper,
  golang-migrate
Для написания unit test's было использовано:
-	github.com/stretchr/testify/assert
- github.com/stretchr/testify/require
- github.com/zhashkevych/go-sqlxmock
- github.com/golang/mock/gomock
