package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Daniorocket/RestApi-SQLite/handlers"
	"github.com/Daniorocket/RestApi-SQLite/routing"
	"github.com/Daniorocket/RestApi-SQLite/sqldb"
	"github.com/gomodule/redigo/redis"
)

func main() {
	var err error
	initCache()
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
func initCache() {
	// Initialize the redis connection to a redis instance running on your local machine
	conn, err := redis.DialURL("redis://localhost")
	if err != nil {
		log.Println("Unable to connect to a redis", err)
		return
	}
	// Assign the connection to the package level `cache` variable
	handlers.Cache = conn
	fmt.Println("Init cache done")
}
