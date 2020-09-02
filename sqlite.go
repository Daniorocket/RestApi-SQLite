package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func selectAllData(db *sql.DB) {
	rows, err := db.Query("SELECT * FROM todos")
	checkErr(err)
	var id int
	var username string
	var department string
	var created time.Time
	fmt.Println("\nSelect all rows")
	for rows.Next() {
		err = rows.Scan(&id, &username, &department, &created)
		checkErr(err)
		fmt.Println("\n", id)
		fmt.Println(username)
		fmt.Println(department)
		fmt.Println(created)

	}
	rows.Close() //good habit to close
	fmt.Println("Done!")
}
func dbCountOfTodos() int {
	rows, err := Db.Query("SELECT COUNT(*) as count FROM  todos")
	checkErr(err)
	return checkCount(rows)
}

func checkCount(rows *sql.Rows) (count int) {
	for rows.Next() {
		err := rows.Scan(&count)
		checkErr(err)
	}
	return count
}
func insertRow(db *sql.DB, id int64, username string, department bool) {
	// insert
	stmt, err := db.Prepare("INSERT INTO todos(id, name, completed,due) values(?,?,?,?)")
	checkErr(err)

	res, err := stmt.Exec(id, username, department, time.Now())
	res.LastInsertId()
	checkErr(err)
}
func updateName(db *sql.DB, id int64, name string) {
	stmt, err := db.Prepare("update todos set name=? where id=?")
	checkErr(err)

	res, err := stmt.Exec(name, id)
	checkErr(err)

	affect, err := res.RowsAffected()
	checkErr(err)

	fmt.Println("Updated: ", affect)
}
func deleteRow(db *sql.DB, id int64) {
	// delete
	stmt, err := db.Prepare("delete from todos where id=?")
	checkErr(err)

	res, err := stmt.Exec(id)
	checkErr(err)

	affect, err := res.RowsAffected()
	checkErr(err)

	fmt.Println("Deleted: ", affect)
}
func createDb() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./todos.db")
	checkErr(err)
	return db, err
}

// 	insertRow(db, "Daniel", "ZUT")
// 	insertRow(db, "Janusz", "PWR")
// 	selectAllData(db)
// 	updateName(db, 1, "Andrzej")
// 	selectAllData(db)
// 	deleteRow(db, 2)
// 	selectAllData(db)
// 	db.Close()
// }

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
