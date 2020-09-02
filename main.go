package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
)

//Db pointer to db
var Db *sql.DB
var err error

func main() {
	Db, err = createDb()
	stmt, err := Db.Prepare("CREATE TABLE IF NOT EXISTS todos (id INTEGER PRIMARY KEY,name varchar(64) NULL,completed bool false,due date NULL)")
	stmt.Exec()
	checkErr(err)
	insertRow(Db, 1, "Daniel", false)
	insertRow(Db, 2, "Janusz", false)
	selectAllData(Db)
	deleteRow(Db, 1)
	selectAllData(Db)
	fmt.Println("Liczba rekordow:", dbCountOfTodos())
	router := NewRouter()

	log.Fatal(http.ListenAndServe(":8081", router))
	Db.Close()
}
