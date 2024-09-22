package main

import (
	"github.com/hari134/comet/migration"
	"github.com/hari134/comet/internal/db"
	"github.com/uptrace/bun/extra/bundebug"
	"log"
)

func main() {
	db := db.NewDB()
	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))

	log.Println("Starting migrations...")
	migration.RunMigrations(db)
	log.Println("Migrations completed.")
}
