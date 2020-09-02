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
