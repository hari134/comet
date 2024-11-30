package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"log"
)

type DBOptions struct {
	DBUser     string
	DBPassword string
	DBName     string
	DBHost     string
	DBPort     string
}

func NewDB(dbOptions DBOptions) *bun.DB {

	dbUser := dbOptions.DBUser
	dbPassword := dbOptions.DBPassword
	dbName := dbOptions.DBName
	dbHost := dbOptions.DBHost
	dbPort := dbOptions.DBPort

	dsn := "postgres://" + dbUser + ":" + dbPassword + "@" + dbHost + ":" + dbPort + "/" + dbName
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatal("Error in postgres dsn")
	}
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer pool.Close()
	sqldb := stdlib.OpenDBFromPool(pool)

	db := bun.NewDB(sqldb, pgdialect.New())
	return db
}
