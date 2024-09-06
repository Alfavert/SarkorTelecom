FROM golang:latest

RUN go version
ENV GOPATH=/

RUN apt-get update && apt-get install -y netcat-openbsd postgresql-client && apt-get clean

COPY ./ ./

RUN chmod +x wait-for-postgres.sh

RUN go mod download

RUN go build -o SarkorTelecom ./cmd/main.go

CMD ["./wait-for-postgres.sh", "db", "5432", "./SarkorTelecom"]
