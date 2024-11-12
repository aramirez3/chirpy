package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/aramirez3/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Printf("error connecting to db: %s\n", err)
		return
	}
	s := createServer("8080")
	s.Config.dbQueries = database.New(db)
	env, err := godotenv.Read()
	if err != nil {
		fmt.Printf("error reading .env file: %s\n", err)
		return
	}
	s.Config.Secret = env["chirpy_secret"]
	s.startServer()
}
