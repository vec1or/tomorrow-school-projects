package main

import (
	"database/sql"
	"flag"

	//"fmt"
	"log"
	"net/http"

	//"strconv"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	db       *sql.DB
}

func main() {

	dsn := flag.String("dsn", "file:forum.db?_foreign_keys=on", "SQLite data source name")
	addr := flag.String("addr", ":8080", "HTTP network address")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Llongfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	query, err := os.ReadFile("schema.sql")
	if err != nil {
		errorLog.Fatalf("Couldn't read schema.sql: %v", err)
	}

	_, err = db.Exec(string(query))
	if err != nil {
		errorLog.Fatalf("Couldn't execute schema.sql: %v", err)
	}
	infoLog.Println("DataBase scheme used successfuly")

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		db:       db,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on: %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(1)

	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
