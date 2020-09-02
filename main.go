package main

import (
	"database/sql"
	"log"
	"net/http"
)

//Db pointer to db
var Db *sql.DB
var err error

func main() {
	router := NewRouter()
	log.Fatal(http.ListenAndServe(":8081", router))
	Db.Close()
}
