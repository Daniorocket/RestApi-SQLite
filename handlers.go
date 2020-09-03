package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Handler struct {
	db *sql.DB
}

const MaxBufferSize = 1048576

func (d *Handler) Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome!")
}

func (d *Handler) UserinfoIndex(w http.ResponseWriter, r *http.Request) {

	users, err := selectAllData(d.db)
	if err != nil {
		log.Println("failed to read rows from  table ", err)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	bytes, err := json.Marshal(&users)
	if err != nil {
		log.Println("failed to prepare json describe list of users ", err)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}

func (d *Handler) UserinfoShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	todoId := vars["todoId"]
	fmt.Fprintln(w, "Todo show:", todoId)
}

func (d *Handler) UserinfoCreate(w http.ResponseWriter, r *http.Request) { //Post
	defer r.Body.Close()
	var userinfo Userinfo
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, MaxBufferSize))
	if err != nil {
		log.Println("failed to read body: ", err)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(body, &userinfo); err != nil {
		log.Println("failed to unmarshal new user: ", err)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Println("failed to read body: ", err)
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	index, err := dbCountOfUserinfo(d.db)
	index = index + 1
	userinfo, err = insertRow(d.db, index, userinfo.Username, userinfo.Departname)
	if err != nil {
		log.Println("failed to insert row on db: ", err)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(userinfo); err != nil {
		log.Println("failed to encode userinfo ", err)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
func (d *Handler) EditUserinfo(w http.ResponseWriter, r *http.Request) {
	var userinfo Userinfo
	var oldUserinfo Userinfo
	vars := mux.Vars(r)
	uid := vars["uid"]
	UidInt, err := strconv.Atoi(uid)
	if err != nil {
		log.Println("Can't convert this ID to integer.", err)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	oldUserinfo, err = selectRowById(d.db, UidInt)
	if err != nil {
		log.Println("unable to select row with this ID ", err)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		log.Println("failed to read json body from edituserinfo ", err)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := r.Body.Close(); err != nil {
		log.Println("failed to close body from edituserinfo ", err)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(body, &userinfo); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Println("failed to encode userinfo ", err)
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	//If user edited username
	if userinfo.Username != "" && userinfo.Username != oldUserinfo.Username {
		err = updateName(d.db, int64(UidInt), userinfo.Username)
		if err != nil {
			log.Println("failed to update name of user ", err)
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	//If user edited departname
	if userinfo.Departname != "" && userinfo.Departname != oldUserinfo.Departname {
		err = updateDepartname(d.db, int64(UidInt), userinfo.Departname)
		if err != nil {
			log.Println("failed to update departname of user ", err)
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
}
func (d *Handler) DeleteUserinfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid := vars["uid"]
	UidInt, err := strconv.Atoi(uid)
	if err != nil {
		log.Println("failed to convers userId to int", err)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = deleteRow(d.db, int64(UidInt))
	if err != nil {
		log.Println("failed to delete row on database", err)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
