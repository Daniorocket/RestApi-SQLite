package sqldb

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Daniorocket/RestApi-SQLite/userinfo"

	_ "github.com/mattn/go-sqlite3"
)

func InitDb(db *sql.DB) error {
	stmt, err := db.Prepare("CREATE TABLE IF NOT EXISTS userinfo (uid INTEGER PRIMARY KEY AUTOINCREMENT,username varchar(64) NULL,departname varchar(64) NULL,created date NULL)")
	if err != nil {
		return err
	}
	stmt.Exec()
	stmt, err = db.Prepare("CREATE TABLE IF NOT EXISTS users (username TEXT PRIMARY KEY,password text)")
	if err != nil {
		return err
	}
	stmt.Exec()
	fmt.Println("Init db done!")
	return nil
}
func SelectAllData(db *sql.DB) ([]userinfo.Userinfo, error) {
	var usersinfo []userinfo.Userinfo
	rows, err := db.Query("SELECT uid,username,departname,created FROM userinfo")
	if err != nil {
		return nil, err
	}
	defer rows.Close() //good habit to close
	for rows.Next() {
		var user userinfo.Userinfo
		err = rows.Scan(
			&user.Uid, &user.Username, &user.Departname, &user.Created)
		if err != nil {
			return nil, err
		}
		usersinfo = append(usersinfo, user)
	}
	return usersinfo, nil
}

func SelectRowFromUserinfoById(db *sql.DB, id int) (userinfo.Userinfo, error) {
	var userinf userinfo.Userinfo
	err := db.QueryRow("select uid,username,departname,created from userinfo where uid = ?", id).Scan(&userinf.Uid, &userinf.Username, &userinf.Departname, &userinf.Created)
	if err != nil {
		return userinf, err
	}
	return userinf, nil
}
func SelectPasswordFromUserByName(db *sql.DB, username string) (string, error) {
	var password string
	err := db.QueryRow("select password from users where username = ?", username).Scan(&password)
	if err != nil {
		return "", err
	}
	return password, nil
}
func DbCountOfUserinfo(db *sql.DB) (int64, error) {
	var count int64
	err := db.QueryRow("SELECT COUNT(*) as count FROM  userinfo").Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func CheckCount(rows *sql.Rows) (int64, error) {
	var count int64
	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			return 0, err
		}
	}
	return count, nil
}
func InsertRowIntoUserinfo(db *sql.DB, uid int64, username string, departname string) (userinfo.Userinfo, error) {
	var user userinfo.Userinfo
	stmt, err := db.Prepare("INSERT INTO userinfo(username, departname,created) values(?,?,?)")
	if err != nil {
		return user, err
	}

	now := time.Now().UTC()
	res, err := stmt.Exec(username, departname, now)
	if err != nil {
		return user, err
	}
	res.LastInsertId()
	user.Uid = int(uid)
	user.Username = username
	user.Departname = departname
	user.Created = now
	return user, nil
}
func InsertRowIntoUsers(db *sql.DB, username string, password string) error {
	fmt.Println(username, password)
	stmt, err := db.Prepare("INSERT INTO users(username, password) values(?,?)")
	if err != nil {
		return err
	}
	res, err := stmt.Exec(username, password)
	if err != nil {
		return err
	}
	res.LastInsertId()
	return nil
}
func UpdateRowById(db *sql.DB, id int64, userinf userinfo.Userinfo) error {
	stmt, err := db.Prepare("update userinfo set username=?, departname=? where uid=?")
	if err != nil {
		return err
	}

	res, err := stmt.Exec(userinf.Username, userinf.Departname, userinf.Uid)
	if err != nil {
		return err
	}

	affect, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affect == 0 {
		return errors.New("Row not found.")
	}
	fmt.Println("Updated: ", affect, " record.")
	return nil
}

func DeleteRow(db *sql.DB, id int64) error {
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

	if affect == 0 {
		return errors.New("Row not found.")
	}

	fmt.Println("Deleted: ", affect, " record.")
	return nil
}

func CreateDb() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./todos.db")
	if err != nil {
		return nil, err
	}
	return db, err
}
