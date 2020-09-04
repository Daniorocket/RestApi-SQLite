package main

import (
	"log"
	"net/http"

	"github.com/Daniorocket/RestApi-SQLite/routing"
	"github.com/Daniorocket/RestApi-SQLite/sqldb"
)

func main() {
	var err error
	db, err := sqldb.CreateDb()
	if err != nil {
		log.Println("Failed to open connection:", err)
		return
	}
	defer db.Close()
	if err = sqldb.InitDb(db); err != nil {
		log.Println("Failed to init db: ", err)
		return
	}
	router := routing.NewRouter(db)
	if err = http.ListenAndServe(":8081", router); err != nil {
		log.Println("Failed to close server: ", err)
		return
	}
}
