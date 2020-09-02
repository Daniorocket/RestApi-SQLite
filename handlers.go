package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome!")
}

func UserinfoIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(userinfoTable); err != nil {
		panic(err)
	}
}

func UserinfoShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	todoId := vars["todoId"]
	fmt.Fprintln(w, "Todo show:", todoId)
}

func UserinfoCreate(w http.ResponseWriter, r *http.Request) { //Post
	var userinfo Userinfo
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
	userinfo = insertRow(Db, dbCountOfUserinfo()+1, userinfo.Username, userinfo.Departname)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(userinfo); err != nil {
		panic(err)
	}
}
func EditUserinfo(w http.ResponseWriter, r *http.Request) {
	var userinfo Userinfo
	var oldUserinfo Userinfo
	vars := mux.Vars(r)
	uid := vars["uid"]
	UidInt, err := strconv.Atoi(uid)
	if err != nil {
		panic(err)
	}

	oldUserinfo = selectRowById(Db, UidInt)
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
		updateName(Db, int64(UidInt), userinfo.Username)
		userTableIndex := searchUserinfoById(UidInt)
		if userTableIndex != -1 {
			userinfoTable[userTableIndex].Username = userinfo.Username
		}
	}
	//If user edited departname
	if userinfo.Departname != "" && userinfo.Departname != oldUserinfo.Departname {
		updateDepartname(Db, int64(UidInt), userinfo.Departname)
		userTableIndex := searchUserinfoById(UidInt)
		if userTableIndex != -1 {
			userinfoTable[userTableIndex].Departname = userinfo.Departname
		}
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(userinfo); err != nil {
		panic(err)
	}
}
func DeleteUserinfo(w http.ResponseWriter, r *http.Request) {
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
	deleteRow(Db, int64(UidInt))
	if sum := deleteElementFromTableById(UidInt); sum == 0 {
		if err := json.NewEncoder(w).Encode("Deleted!"); err != nil {
			panic(err)
		}
	} else {
		if err := json.NewEncoder(w).Encode("Not deleted"); err != nil {
			panic(err)
		}
	}
	// body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	// if err != nil {
	// 	panic(err)
	// }
	// if err := r.Body.Close(); err != nil {
	// 	panic(err)
	// }
	// if err := json.Unmarshal(body, &confirmation); err != nil {
	// 	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	// 	w.WriteHeader(http.StatusOK)
	// 	fmt.Println("Aha", confirmation.Confirmed)
	// }
}
