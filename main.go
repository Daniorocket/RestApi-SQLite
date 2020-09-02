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
	Db, err = createDb()
	stmt, err := Db.Prepare("CREATE TABLE IF NOT EXISTS userinfo (uid INTEGER PRIMARY KEY AUTOINCREMENT,username varchar(64) NULL,departname varchar(64) NULL,created date NULL)")
	stmt.Exec()
	checkErr(err)
	insertRow(Db, 1, "Daniel", "ZUT")
	insertRow(Db, 2, "Janusz", "PWR")
	selectAllData(Db)
	//	fmt.Println("Liczba rekordow:", dbCountOfUserinfo())
	router := NewRouter()

	log.Fatal(http.ListenAndServe(":8081", router))
	Db.Close()
}
