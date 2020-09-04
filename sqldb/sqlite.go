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
	stmt.Exec()
	fmt.Println("Init db done!")
	return err
}
func SelectAllData(db *sql.DB) ([]userinfo.Userinfo, error) {
	var usersinfo []userinfo.Userinfo
	rows, err := db.Query("SELECT * FROM userinfo")
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

func SelectRowById(db *sql.DB, id int) (userinfo.Userinfo, error) {
	var userinfo userinfo.Userinfo
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
func InsertRow(db *sql.DB, uid int64, username string, departname string) (userinfo.Userinfo, error) {
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

func UpdateName(db *sql.DB, id int64, name string) error {
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
	if affect == 0 {
		return errors.New("Row not found.")
	}
	fmt.Println("Updated: ", affect, " record.")
	return nil
}
func UpdateDepartname(db *sql.DB, id int64, departname string) error {
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
