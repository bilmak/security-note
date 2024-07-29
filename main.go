package main

import (
	"database/sql"
	_ "embed"
	"flag"
	"fmt"
	"net/http"
	"os"
	"password/hendlers"
	"password/model"

	_ "github.com/lib/pq"
)

//go:embed postNotes.html
var postNotes string

//go:embed getNotes.html
var getNotes string

type NotesModel struct {
	DB  *sql.DB
	Dsn string
}

func main() {
	db := connecttionToPostgres()
	lastRow := model.NotesModel{DB: db}
	counter, err := lastRow.GetLastRowFromDB()
	if err != nil {
		fmt.Println("main, counter erorr", err)
	}

	note := hendlers.NotesServer{
		NoteModel: lastRow,
		Counter:   counter,
	}
	http.HandleFunc("/passpost", note.PassPost)
	http.HandleFunc("/passget", note.PassGet)

	http.HandleFunc("/post", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Content-Type", "text/html")
		w.Write([]byte(postNotes))
	})

	http.HandleFunc("/get", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Content-Type", "text/html")
		w.Write([]byte(getNotes))
	})

	http.ListenAndServe(":8989", nil)

}

func connecttionToPostgres() *sql.DB {
	var config NotesModel

	flag.StringVar(&config.Dsn, "db-dsn", os.Getenv("READINGLIST_DB_DSN"), "postgresSQL DSN")
	flag.Parse()

	db, err := sql.Open("postgres", config.Dsn)
	if err != nil {
		fmt.Print(err)

	}

	err = db.Ping()
	if err != nil {
		fmt.Print(err)

	}
	fmt.Println("\nSuccessfully connected!")
	return db
}
