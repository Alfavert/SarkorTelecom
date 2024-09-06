# Тестовое задание для SarkorTelecom на позицию Junior разработчика
тестовое - https://disk.yandex.ru/i/6v-8xGiG6_tQ0A
1) Скачать репозиторий
2) В консоли перейти в директорию в которую вы скачали репозиторий
3) Запустить команду docker-compose up --build sarkortelecom
4) Подождать запуска программы
5) Проверку API осуществлял через Postman
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
