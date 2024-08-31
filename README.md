# SarkorTelekom
1) Настроить данные для подключения к БД в файле configs/config (верхнее значение порта используется для самого приложения)
2) Запустить БД, затем запустить cmd/main.go файл
3) Проверку API осуществлял через Postman (
   POST "/products" - добавляет запись в таблицу, возвращает id записи
   GET "/product" - указываете id, возвращает информацию о продукте
   PUT "/product" - указываете id, обновляет содержимое
   DELETE "/product" - указываете id, удаляет указанный продукт
   GET "/products" -  возвращает список всего содержимого таблицы в БД 
)


Использовал в проекте - Gin - github.com/gin-gonic/gin, 
  sqlx - github.com/jmoiron/sqlx, 
  logrus - github.com/sirupsen/logrus, 
  viper - github.com/spf13/viper
