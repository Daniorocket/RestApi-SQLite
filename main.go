package main

import (
	"log"
	"net/http"
)

func main() {
	var err error
	db, err := createDb()
	if err != nil {
		log.Println("failed to open connection")
		return
	}
	defer db.Close()
	if err != nil {
		log.Println("failed to close connection")
		return
	}
	if err = initDb(db); err != nil {
		log.Println("Failed to init db: ", err)
		return
	}
	router := NewRouter(db)
	if err = http.ListenAndServe(":8081", router); err != nil {
		log.Println("failed to close server: ", err)
		return
	}

}
