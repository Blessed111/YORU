package main

import (
	"Diploma"
	"Diploma/pkg/handlers"
	"Diploma/pkg/repository"
	service "Diploma/pkg/service"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	if err := initConfig(); err != nil {
		logrus.Fatalf("ConfigurationError:%s", err.Error())
	}
	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error:%s", err.Error())
	}
	// db, err := repository.NewPostgresDb(repository.Config{
	// 	Host:     viper.GetString("db.host"),
	// 	Port:     viper.GetString("db.port"),
	// 	Username: viper.GetString("db.username"),
	// 	Password: os.Getenv("DB_PASSWORD"),
	// 	DbName:   viper.GetString("db.dbname"),
	// 	Sslmode:  viper.GetString("db.sslmode"),
	// })
	db, err := repository.NewPostgresDb1()
	if err != nil {
		logrus.Fatalf("24 row error on db:%s", err.Error())
		return
	}

	repos := repository.NewConnection(db)

	service := service.NewService(repos)
	handler := handlers.NewHandler(service)
	server := new(Diploma.Server)
	if err := server.Run(viper.GetString("port"), handler.Init()); err != nil {
		logrus.Fatalf("Error: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("../configs/")
	viper.SetConfigName("config")

	return viper.ReadInConfig()
}
