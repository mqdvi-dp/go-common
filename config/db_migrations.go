package config

import (
	"context"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/mqdvi-dp/go-common/env"
	"github.com/mqdvi-dp/go-common/logger"
	"github.com/pressly/goose/v3"
)

func DbAutoMigrations() {
	folderPath := env.GetString("PATH_DB_MIGRATIONS", "db/migrations")
	defer logger.Blue("finish db migrations")

	// get the folder path
	fs, err := os.Stat(folderPath)
	if err != nil {
		if os.IsNotExist(err) {
			logger.RedItalic("no have migrations")
			return
		}

		logger.Red(fmt.Errorf("get folder path of db/migrations is error: %s", err))
		return
	}

	if !fs.IsDir() {
		logger.Red("db/migrations is not directory")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(3)*time.Minute)
	defer cancel()

	// open connection with goose driver
	db, err := goose.OpenDBWithDriver("postgres", env.GetString("DSN_MASTER"))
	if err != nil {
		logger.Red(fmt.Errorf("failed open connection database for migrations: %s", err))
		return
	}
	defer db.Close()

	// running the migrations
	if err = goose.RunContext(ctx, "up", db, folderPath); err != nil {
		logger.Red(fmt.Errorf("migrations error: %s", err))
		return
	}
}
