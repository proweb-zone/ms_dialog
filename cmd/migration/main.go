package main

import (
	"database/sql"
	"flag"
	"fmt"

	"ms_dialog/internal/config"
	"ms_dialog/internal/db/postgres"
	"ms_dialog/internal/migrator"
	"ms_dialog/internal/utils"
)

func main() {

	currentDir := utils.GetProjectPath()
	configPath := config.PathDefault(currentDir, nil)
	config := config.MustInit(configPath)

	var action string
	flag.StringVar(&action, "action", "up", "path to config file")
	flag.Parse()

	conn := postgres.Connect(config)
	defer postgres.Close(conn)

	migrator := migrator.MustGetNewMigrator(config.Db.Name)
	switchAndExecMigrateAction(action, migrator, conn)
}

func switchAndExecMigrateAction(action string, migrator *migrator.Migrator, conn *sql.DB) {
	switch action {
	case "up":
		if err := migrator.Up(conn); err != nil {
			fmt.Println("Up migration failed", err)
		} else {
			fmt.Println("Up migration success.")
		}
	case "down":
		if err := migrator.Down(conn); err != nil {
			fmt.Println("Down migration failed", err)
		} else {
			fmt.Println("Down migration success.")
		}
	default:
		fmt.Println("parameter -action up or down not found")
	}
}
