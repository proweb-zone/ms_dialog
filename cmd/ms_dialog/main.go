package main

import (
	"ms_dialog/internal/app"
	"ms_dialog/internal/config"
	"ms_dialog/internal/utils"
)

func main() {
	currentDir := utils.GetProjectPath()
	configPath := config.ParseConfigPathFromCl(currentDir)
	config := config.MustInit(configPath)
	app.InitApp(config)
}
