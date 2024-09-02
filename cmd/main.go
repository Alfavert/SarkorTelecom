package main

import (
	"SarkorTelekom"
	api "SarkorTelekom/pkg/handler"
	"SarkorTelekom/pkg/repository"
	"SarkorTelekom/pkg/service"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	if err := initConfig(); err != nil {
		logrus.Fatalf("Error initializing config file: %s", err.Error())
	}

	db, err := repository.ConnectPostgres(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: viper.GetString("db.password"),
		DBName:   viper.GetString("db.DBName"),
		SSLMode:  viper.GetString("db.SSLMode"),
	})
	if err != nil {
		logrus.Fatalf("Error connecting to database: %s", err.Error())
	}
	repository.CreateProductTable(db)
	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := api.NewHandler(services)
	srv := new(SarkorTelekom.Server)
	if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		logrus.Fatalf("error while running http server", err.Error())
	}
}
func initConfig() error {
	viper.AddConfigPath("./configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
