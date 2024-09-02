FROM golang:latest

RUN go version
ENV GOPATH=/

# Устанавливаем зависимости
RUN apt-get update && apt-get install -y netcat-openbsd postgresql-client && apt-get clean

# Копируем файлы в контейнер
COPY ./ ./

# Делаем wait-for-postgres.sh исполняемым
RUN chmod +x wait-for-postgres.sh

# Загружаем модули Go
RUN go mod download

# Собираем Go приложение
RUN go build -o SarkorTelecom ./cmd/main.go

# Используем wait-for-postgres.sh для ожидания и запускаем приложение
CMD ["./wait-for-postgres.sh", "db", "5432", "./SarkorTelecom"]
