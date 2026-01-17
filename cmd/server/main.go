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

func setupDB() *db.Queries {
	dbUrl := os.Getenv("DB_URL")
	database, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal(err)
	}
	dbQuery := db.New(database)
	return dbQuery
}

func main() {
	godotenv.Load()
	dbQuery := setupDB()
	notesService := &notes.Service{
		Q: dbQuery,
	}
	mux := http.NewServeMux()
	server := http.Server{
		Addr:    ":5027",
		Handler: mux,
	}
	mux.HandleFunc("/app/notes", notesService.NotesCollectionHandler)
	mux.HandleFunc("/app/notes/{id}", notesService.NoteItemHandler)
	println("Server started...")
	server.ListenAndServe()

}
