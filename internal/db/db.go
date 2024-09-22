package db

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)


func NewDB() *bun.DB{
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	dbHost := os.Getenv("POSTGRES_HOST")
	dbPort := os.Getenv("POSTGRES_PORT")

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