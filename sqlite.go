package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func selectAllData(db *sql.DB) {
	rows, err := db.Query("SELECT * FROM userinfo")
	checkErr(err)
	var uid int
	var username string
	var departname string
	var created time.Time
	fmt.Println("\nSelect all rows")
	for rows.Next() {
		err = rows.Scan(&uid, &username, &departname, &created)
		checkErr(err)
		fmt.Println("\n", uid)
		fmt.Println(username)
		fmt.Println(departname)
		fmt.Println(created)
	}
	rows.Close() //good habit to close
	fmt.Println("Done!")
}
func selectRowById(db *sql.DB, id int) Userinfo {
	var userinfo Userinfo
	rows, err := db.Query("select * from userinfo where uid = ?", id)
	if err != nil {
		checkErr(err)
	}
	var uid int
	var username string
	var departname string
	var created time.Time
	for rows.Next() {
		err := rows.Scan(&uid, &username, &departname, &created)
		if err != nil {
			panic(err)
		}
		userinfo.Uid = uid
		userinfo.Username = username
		userinfo.Departname = departname
		userinfo.Created = created
	}
	return userinfo
}
func dbCountOfUserinfo() int64 {
	rows, err := Db.Query("SELECT COUNT(*) as count FROM  userinfo")
	checkErr(err)
	return checkCount(rows)
}

func checkCount(rows *sql.Rows) (count int64) {
	for rows.Next() {
		err := rows.Scan(&count)
		checkErr(err)
	}
	return count
}
func insertRow(db *sql.DB, uid int64, username string, departname string) Userinfo {
	// insert
	var user Userinfo
	stmt, err := db.Prepare("INSERT INTO userinfo(username, departname,created) values(?,?,?)")
	checkErr(err)

	res, err := stmt.Exec(username, departname, time.Now())
	res.LastInsertId()
	checkErr(err)
	user.Uid = int(uid)
	user.Username = username
	user.Departname = departname
	user.Created = time.Now()
	userinfoTable = append(userinfoTable, user)
	return user
}
func updateName(db *sql.DB, id int64, name string) {
	stmt, err := db.Prepare("update userinfo set username=? where uid=?")
	checkErr(err)

	res, err := stmt.Exec(name, id)
	checkErr(err)

	affect, err := res.RowsAffected()
	checkErr(err)

	fmt.Println("Updated: ", affect)

}
func updateDepartname(db *sql.DB, id int64, departname string) {
	stmt, err := db.Prepare("update userinfo set departname=? where uid=?")
	checkErr(err)

	res, err := stmt.Exec(departname, id)
	checkErr(err)

	affect, err := res.RowsAffected()
	checkErr(err)

	fmt.Println("Updated: ", affect)
}
func deleteRow(db *sql.DB, id int64) int64 {
	// delete
	stmt, err := db.Prepare("delete from userinfo where uid=?")
	checkErr(err)

	res, err := stmt.Exec(id)
	checkErr(err)

	affect, err := res.RowsAffected()
	checkErr(err)

	fmt.Println("Deleted: ", affect)
	return affect
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
