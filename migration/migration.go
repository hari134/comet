package migration

import (
	"context"
	"embed"
	"log"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/migrate"
)

//go:embed scripts/*.sql
var sqlMigrations embed.FS

func RunMigrations(db *bun.DB) {
	migrations := migrate.NewMigrations()
	if err := migrations.Discover(sqlMigrations); err != nil {
		panic(err)
	}

	migrator := migrate.NewMigrator(db, migrations)

	// Initialize the migrations table if it doesn't exist
	ctx := context.Background()
	if err := migrator.Init(ctx); err != nil {
		log.Fatalf("Failed to initialize migrations table: %v", err)
	}

	// Run pending migrations
	group, err := migrator.Migrate(ctx)
	if err != nil {
		log.Fatalf("Failed to migrate: %v", err)
	}

	// Log the status of the migrations
	if group.ID == 0 {
		log.Println("There are no new migrations to run.")
	} else {
		log.Printf("Migrated to group %d.\n", group.ID)
	}
}
