package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/spf13/viper"
)

var (
	err error
)

func init() {
	// подключение к файлу конфигурации
	viper.AddConfigPath("config")
	viper.SetConfigName("config")
	err = viper.ReadInConfig()
	if err != nil {
		log.Fatal("Config file not found")
	}
}

func main() {

	host := viper.GetString("WEBSERVER.HOST") + ":" + viper.GetString("WEBSERVER.PORT")
	router := chi.NewRouter()

	/*
		- Получение данных библиотеки с фильтрацией по всем полям и
		пагинацией
		- Получение текста песни с пагинацией по куплетам
		- Удаление песни
		- Изменение данных песни
		- Добавление новой песни в формате
	*/

	if err = http.ListenAndServe(host, router); err != nil {
		fmt.Println("Web server error:", err.Error())
		return
	}
}
