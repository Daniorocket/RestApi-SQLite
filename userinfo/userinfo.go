package userinfo

import (
	"errors"
	"time"
)

type Userinfo struct {
	Uid        int       `json:"uid"`
	Username   string    `json:"username"`
	Departname string    `json:"departname"`
	Created    time.Time `json:"created"`
}

func CheckUsername(newusername string, oldusername string) error {
	if newusername == "" {
		return errors.New("New username can't be empty")
	}
	if newusername == oldusername {
		return errors.New("New username must be changed")
	}
	return nil
}
func CheckDepartname(newdepartname string, olddepartname string) error {
	if newdepartname == "" {
		return errors.New("New departname can't be empty")
	}
	if newdepartname == olddepartname {
		return errors.New("New departname must be changed")
	}
	return nil
}
