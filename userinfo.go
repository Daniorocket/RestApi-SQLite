package main

import (
	"time"
)

type Userinfo struct {
	Uid        int       `json:"uid"`
	Username   string    `json:"username"`
	Departname string    `json:"departname"`
	Created    time.Time `json:"created"`
}

// type Userinfos []Userinfo

// var userinfoTable Userinfos

// func searchUserinfoById(id int) (int, error) {
// 	for i, userinfo := range userinfoTable {
// 		if userinfo.Uid == id {
// 			return i, nil
// 		}
// 	}
// 	return 0, errors.New("User not found in database.")
// }
// func deleteElementFromTableById(id int) error {
// 	for i, userinfo := range userinfoTable {
// 		if userinfo.Uid == id {
// 			userinfoTable = append(userinfoTable[:i], userinfoTable[i+1:]...)
// 			return errors.New("Cant delete user(not exists)")
// 		}
// 	}
// 	return nil
// }
