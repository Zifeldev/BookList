package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

var Conn *pgx.Conn

func ConnectDB() {
	var err error
	

	err = godotenv.Load()
	if err != nil {
		fmt.Println("Error with connecting", err)
		os.Exit(1)
	}
	dbURL := os.Getenv("DB_URL")
	fmt.Println("DB_URL:", os.Getenv("DB_URL"))
	if dbURL == "" {
		fmt.Println("DBURL is empty")
	}

	Conn, err = pgx.Connect(context.Background(), dbURL)
	if err != nil {
		fmt.Println("fail", err)
		os.Exit(1)
	}

	fmt.Println("Connected")

}
