package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"EffectiveMobile/internal/api"
	"EffectiveMobile/internal/service"
	"EffectiveMobile/internal/storage"
)

var (
	err      error
	logLevel int
	logger   *zap.Logger
	db       *pgxpool.Pool
	stores   *storage.Storage
	services *service.Service
	handlers *api.ApiHandler
)

func init() {
	// Подключение к файлу переменных окружения
	if err = godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}

	logLevel, err = strconv.Atoi(os.Getenv("LOGGER_LEVEL"))
	if err != nil {
		log.Fatal("LOGGER_LEVEL is not set")
	}

	var level zapcore.Level
	level = zapcore.Level(logLevel)
	logCfg := zap.NewDevelopmentConfig()
	logCfg.Level = zap.NewAtomicLevelAt(level)
	logger, err = logCfg.Build()
	if err != nil {
		fmt.Println("failed to build")
		return
	}
	defer logger.Sync()

	logger.Info("Initializing success")
}

func main() {
	sugar := logger.Sugar()

	// Подключение к базе данных
	config := storage.Config{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		Username: os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Name:     os.Getenv("POSTGRES_NAME"),
		SSLMode:  os.Getenv("POSTGRES_SSLMODE"),
	}
	db, err = storage.NewPostgresDB(config)
	if err != nil {
		sugar.Fatalf("DB connection error: %s", err.Error())
	}
	defer db.Close()

	stores = storage.NewStorage(db, sugar)
	services = service.NewService(stores, sugar)
	handlers = api.NewHandler(services, sugar)

	// Запуск веб сервера
	httpServerExitDone := &sync.WaitGroup{}
	httpServerExitDone.Add(1)
	serverAddress := fmt.Sprintf("%s:%s", os.Getenv("WEBSERVER_HOST"), os.Getenv("WEBSERVER_PORT"))

	srv := api.StartHttpServer(httpServerExitDone, serverAddress, handlers.InitRoutes())
	sugar.Infof("Application started")

	// Setting up signal capturing
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = srv.Shutdown(ctx); err != nil {
		sugar.Fatalf("Web server shutting down error: %s", err.Error())
	}

	httpServerExitDone.Wait()
	sugar.Infof("Application shutdown")
}
