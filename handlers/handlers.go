package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/Daniorocket/RestApi-SQLite/sqldb"
	"github.com/Daniorocket/RestApi-SQLite/userinfo"

	"github.com/gorilla/mux"
)

type Handler struct {
	Db *sql.DB
}

const MaxBufferSize = 1048576

func (d *Handler) Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome!")
}

func (d *Handler) UserinfoIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	users, err := sqldb.SelectAllData(d.Db)
	if err != nil {
		log.Println("Failed to read rows from  table: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	bytes, err := json.Marshal(&users)
	if err != nil {
		log.Println("Failed to prepare json describe list of users: ", err)
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
	var userinfo userinfo.Userinfo
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, MaxBufferSize))
	if err != nil {
		log.Println("Failed to read body: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(body, &userinfo); err != nil {
		log.Println("Failed to unmarshal new user: ", err)
		w.WriteHeader(http.StatusUnprocessableEntity) // unprocessable entity
		return
	}
	index, err := sqldb.DbCountOfUserinfo(d.Db)
	index = index + 1
	userinfo, err = sqldb.InsertRow(d.Db, index, userinfo.Username, userinfo.Departname)
	if err != nil {
		log.Println("Failed to insert row on db: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := json.NewEncoder(w).Encode(userinfo); err != nil {
		log.Println("Failed to encode userinfo: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
func (d *Handler) EditUserinfo(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var userinf userinfo.Userinfo
	var oldUserinfo userinfo.Userinfo
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	uid := vars["uid"]
	UidInt, err := strconv.Atoi(uid)
	if err != nil {
		log.Println("Can't convert this ID to integer: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	oldUserinfo, err = sqldb.SelectRowById(d.Db, UidInt)
	if err != nil {
		log.Println("Unable to select row with this ID: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err = json.NewDecoder(r.Body).Decode(&userinf); err != nil {
		log.Println("Failed to decode body request: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := json.NewEncoder(w).Encode(err); err != nil {
		log.Println("Failed to encode userinfo: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//If user edited username update it

	if err := userinfo.CheckUsername(userinf.Username, oldUserinfo.Username); err == nil {
		oldUserinfo.Username = userinf.Username
	}

	//If user edited departname update it

	if err := userinfo.CheckDepartname(userinf.Departname, oldUserinfo.Departname); err == nil {
		oldUserinfo.Departname = userinf.Departname
	}
	if err := sqldb.UpdateRowById(d.Db, int64(UidInt), oldUserinfo); err != nil {
		log.Println("Failed to update row: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
func (d *Handler) DeleteUserinfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	uid := vars["uid"]
	UidInt, err := strconv.Atoi(uid)
	if err != nil {
		log.Println("Failed to convers userId to int: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err = sqldb.DeleteRow(d.Db, int64(UidInt)); err != nil {
		log.Println("Failed to delete row on database: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
