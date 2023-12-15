package main

import (
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"snippetbox.markian.com/internal/models"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "bobbyshmurda66"
	dbname   = "snippetbox"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	snippets      *models.SnippetModel
	templateCache map[string]*template.Template
}

func main() {

	addr := flag.String("addr", ":4000", "HTTP network address")
	databaseConnection := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "EERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(databaseConnection)

	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	app := application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		snippets:      &models.SnippetModel{DB: db},
		templateCache: templateCache,
	}
	server := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Println("server runnig on http://localhost:4000")
	err = server.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(databaseConnection string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseConnection)
	if err != nil {
		return nil, err

	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
