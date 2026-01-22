package main

import (
	"database/sql"
	"log"
	"net/http"
	"notes-api/internal/db"
	"notes-api/internal/notes"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func setupDB() (*db.Queries, string) {
	dbUrl := os.Getenv("DB_URL")
	tokenString := os.Getenv("TOKEN_SECRET")
	database, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal(err)
	}
	dbQuery := db.New(database)
	return dbQuery, tokenString
}

func main() {
	godotenv.Load()
	dbQuery, tokenSecret := setupDB()
	notesService := &notes.Service{
		Q:           dbQuery,
		TokenSecret: tokenSecret,
	}
	mux := http.NewServeMux()
	registerRouter(mux, notesService)
	server := http.Server{
		Addr:    ":5027",
		Handler: mux,
	}
	println("Server started...")
	server.ListenAndServe()

}
