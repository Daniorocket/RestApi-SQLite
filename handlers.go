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
	// w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	// w.WriteHeader(http.StatusOK)
	// if err := json.NewEncoder(w).Encode(userinfoTable); err != nil {
	// 	panic(err)
	// }

	// users, err := selectAllData(Db)

	users, err := selectAllData(d.db)
	if err != nil {
		log.Println("failed to read rows from  table ", err)
		return
	}
	bytes, err := json.Marshal(&users)
	if err != nil {
		log.Println("failed to prepare json describe list of users ", err)
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
	}
}
func (d *Handler) EditUserinfo(w http.ResponseWriter, r *http.Request) {
	var userinfo Userinfo
	var oldUserinfo Userinfo
	vars := mux.Vars(r)
	uid := vars["uid"]
	UidInt, err := strconv.Atoi(uid)
	if err != nil {
		panic(err)
	}

	oldUserinfo, err = selectRowById(d.db, UidInt)
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &userinfo); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	//If user edited username
	if userinfo.Username != "" && userinfo.Username != oldUserinfo.Username {
		updateName(d.db, int64(UidInt), userinfo.Username)
		// if userTableIndex, err := searchUserinfoById(UidInt); err != nil {
		// 	//TODO
		// 	//panic(err)
		// } else {
		// 	userinfoTable[userTableIndex].Username = userinfo.Username
		// }
	}
	//If user edited departname
	if userinfo.Departname != "" && userinfo.Departname != oldUserinfo.Departname {
		updateDepartname(d.db, int64(UidInt), userinfo.Departname)
		// if userTableIndex, err := searchUserinfoById(UidInt); err != nil {
		// 	//TODO
		// 	//panic(err)
		// } else {
		// 	userinfoTable[userTableIndex].Departname = userinfo.Departname
		// }
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(userinfo); err != nil {
		panic(err)
	}
}
func (d *Handler) DeleteUserinfo(w http.ResponseWriter, r *http.Request) {
	// type Confirmation struct {
	// 	Confirmed string `json:"confirmed"`
	// }
	// var confirmation Confirmation
	vars := mux.Vars(r)
	uid := vars["uid"]
	UidInt, err := strconv.Atoi(uid)
	if err != nil {
		panic(err)
	}
	deleteRow(d.db, int64(UidInt))
	// if err := deleteElementFromTableById(UidInt); err == nil {
	// 	if err := json.NewEncoder(w).Encode("Deleted!"); err != nil {
	// 		panic(err)
	// 	}
	// } else {
	// 	if err := json.NewEncoder(w).Encode("Not deleted"); err != nil {
	// 		panic(err)
	// 	}
	// }
}
