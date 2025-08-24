package app

import (
	"ms_dialog/internal/config"
	"ms_dialog/internal/server"
)

func InitApp(config *config.Config) {
	server.StartServer(config) // запуск сервера для обработки клентских запросов
}
