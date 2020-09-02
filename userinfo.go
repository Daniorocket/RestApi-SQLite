package main

import "time"

type Userinfo struct {
	Uid        int       `json:"uid"`
	Username   string    `json:"username"`
	Departname string    `json:"departname"`
	Created    time.Time `json:"created"`
}

type Userinfos []Userinfo

var userinfoTable Userinfos

func searchUserinfoById(id int) int {
	for i, userinfo := range userinfoTable {
		if userinfo.Uid == id {
			return i
		}
	}
	return -1
}
