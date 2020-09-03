package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func initDb(db *sql.DB) error {
	stmt, err := db.Prepare("CREATE TABLE IF NOT EXISTS userinfo (uid INTEGER PRIMARY KEY AUTOINCREMENT,username varchar(64) NULL,departname varchar(64) NULL,created date NULL)")
	stmt.Exec()
	fmt.Println("Init db done!")
	return err
}

func selectAllData(db *sql.DB) ([]Userinfo, error) {
	var usersinfo []Userinfo
	rows, err := db.Query("SELECT * FROM userinfo")
	if err != nil {
		return nil, err
	}
	defer rows.Close() //good habit to close
	for rows.Next() {
		var user Userinfo
		err = rows.Scan(
			&user.Uid, &user.Username, &user.Departname, &user.Created)
		if err != nil {
			return nil, err
		}
		usersinfo = append(usersinfo, user)
	}
	return usersinfo, nil
}

func selectRowById(db *sql.DB, id int) (Userinfo, error) {
	var userinfo Userinfo
	rows, err := db.Query("select * from userinfo where uid = ?", id)
	if err != nil {
		return userinfo, err
	}
	var uid int
	var username string
	var departname string
	var created time.Time
	for rows.Next() {
		err := rows.Scan(&uid, &username, &departname, &created)
		if err != nil {
			return userinfo, err
		}
		userinfo.Uid = uid
		userinfo.Username = username
		userinfo.Departname = departname
		userinfo.Created = created
	}
	return userinfo, nil
}
func dbCountOfUserinfo(db *sql.DB) (int64, error) {
	rows, err := db.Query("SELECT COUNT(*) as count FROM  userinfo")
	if err != nil {
		return 0, err
	}
	return checkCount(rows)
}

func checkCount(rows *sql.Rows) (int64, error) {
	var count int64
	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			return 0, err
		}
	}
	return count, nil
}
func insertRow(db *sql.DB, uid int64, username string, departname string) (Userinfo, error) {
	var user Userinfo
	stmt, err := db.Prepare("INSERT INTO userinfo(username, departname,created) values(?,?,?)")
	if err != nil {
		return user, err
	}

	res, err := stmt.Exec(username, departname, time.Now())
	if err != nil {
		return user, err
	}
	res.LastInsertId()
	user.Uid = int(uid)
	user.Username = username
	user.Departname = departname
	user.Created = time.Now()
	return user, nil
}
func updateName(db *sql.DB, id int64, name string) error {
	stmt, err := db.Prepare("update userinfo set username=? where uid=?")
	if err != nil {
		return err
	}

	res, err := stmt.Exec(name, id)
	if err != nil {
		return err
	}

	affect, err := res.RowsAffected()
	if err != nil {
		return err
	}

	fmt.Println("Updated: ", affect, " record.")
	return nil
}
func updateDepartname(db *sql.DB, id int64, departname string) error {
	stmt, err := db.Prepare("update userinfo set departname=? where uid=?")
	if err != nil {
		return err
	}

	res, err := stmt.Exec(departname, id)
	if err != nil {
		return err
	}

	affect, err := res.RowsAffected()
	if err != nil {
		return err
	}

	fmt.Println("Updated: ", affect, " record.")
	return nil
}
func deleteRow(db *sql.DB, id int64) error {
	// delete
	stmt, err := db.Prepare("delete from userinfo where uid=?")
	if err != nil {
		return err
	}

	res, err := stmt.Exec(id)
	if err != nil {
		return err
	}

	affect, err := res.RowsAffected()
	if err != nil {
		return err
	}

	fmt.Println("Deleted: ", affect, " record.")
	return nil
}

func createDb() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./todos.db")
	if err != nil {
		return nil, err
	}
	return db, err
}
