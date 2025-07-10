package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func ConnectDB() (*pgx.Conn, error) {
	loadErr := godotenv.Load(".env")
	if loadErr != nil {
		fmt.Fprintf(os.Stderr, "unable to load variables")
	}
	postgresUrl := os.Getenv("POSTGRES_URL")
	conn, connErr := pgx.Connect(context.Background(), postgresUrl)
	if connErr != nil {
		fmt.Fprintf(os.Stderr, "unable to connect")
	}

	return conn, connErr
}
