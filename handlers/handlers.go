package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Daniorocket/RestApi-SQLite/sqldb"
	"github.com/Daniorocket/RestApi-SQLite/userinfo"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"

	"golang.org/x/crypto/bcrypt"
)

var Cache redis.Conn

type Handler struct {
	Db *sql.DB
}

type Credentials struct {
	Password string `json:"password", db:"password"`
	Username string `json:"username", db:"username"`
}

func UserLoggedIn(w http.ResponseWriter, r *http.Request) error {
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			w.WriteHeader(http.StatusUnauthorized)
			return err
		}
		// For any other type of error, return a bad request status
		w.WriteHeader(http.StatusBadRequest)
		return err
	}
	sessionToken := c.Value

	// We then get the name of the user from our cache, where we set the session token
	response, err := Cache.Do("GET", sessionToken)
	if err != nil {
		// If there is an error fetching from cache, return an internal server error status
		log.Println("Error fetching from cache:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}
	if response == nil {
		// If the session token is not present in cache, return an unauthorized error
		err = errors.New("Session token is not present in cache: ")
		w.WriteHeader(http.StatusUnauthorized)
		return err
	}
	// Finally, return the welcome message to the user
	fmt.Printf("User %s is logged in, operation possible.!\n", response)
	//w.Write([]byte(fmt.Sprintf("Welcome %s!", response)))
	return nil
}

const MaxBufferSize = 1048576

func (d *Handler) Index(w http.ResponseWriter, r *http.Request) {
	if err := UserLoggedIn(w, r); err != nil {
		log.Println("Error authentication, user is not logged in:", err)
		return
	}
	fmt.Println("Index works - user is logged in!")
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
	w.Write(bytes)
}

func (d *Handler) UserinfoShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	todoId := vars["todoId"]
	fmt.Fprintln(w, "Todo show:", todoId)
}

func (d *Handler) UserinfoCreate(w http.ResponseWriter, r *http.Request) { //Post
	if err := UserLoggedIn(w, r); err != nil {
		log.Println("Error authentication, user is not logged in:", err)
		return
	}
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
	userinfo, err = sqldb.InsertRowIntoUserinfo(d.Db, index, userinfo.Username, userinfo.Departname)
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
	if err := UserLoggedIn(w, r); err != nil {
		log.Println("Error authentication, user is not logged in:", err)
		return
	}
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

	oldUserinfo, err = sqldb.SelectRowFromUserinfoById(d.Db, UidInt)
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
	if err := UserLoggedIn(w, r); err != nil {
		log.Println("Error authentication, user is not logged in:", err)
		return
	}
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
}
func (d *Handler) Register(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	creds := &Credentials{}
	err := json.NewDecoder(r.Body).Decode(creds)
	if err != nil {
		// If there is something wrong with the request body, return a 400 status
		log.Println("Failed to decode body", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Salt and hash the password using the bcrypt algorithm
	// The second argument is the cost of hashing, which we arbitrarily set as 8 (this value can be more or less, depending on the computing power you wish to utilize)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), 8)
	if err != nil {
		log.Println("Failed to hash password using bcrypt:", err)
		w.WriteHeader(http.StatusBadRequest)
	}

	// Next, insert the username, along with the hashed password into the database
	if err := sqldb.InsertRowIntoUsers(d.Db, creds.Username, string(hashedPassword)); err != nil {
		log.Println("Failed to insert row into Users:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
func (d *Handler) Login(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r)
	defer r.Body.Close()
	creds := &Credentials{}
	err := json.NewDecoder(r.Body).Decode(creds)
	if err != nil {
		// If there is something wrong with the request body, return a 400 status
		log.Println("Failed to decode body", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Get the existing entry present in the database for the given username
	password, err := sqldb.SelectPasswordFromUserByName(d.Db, creds.Username)
	if err != nil {
		// If there is an issue with the database, return a 500 error
		log.Println("Failed to select row from database: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// We create another instance of `Credentials` to store the credentials we get from the database
	storedCreds := &Credentials{}
	// Store the obtained password in `storedCreds`
	storedCreds.Password = password
	if err != nil {
		// If an entry with the username does not exist, send an "Unauthorized"(401) status
		if err == sql.ErrNoRows {
			log.Println("Failed to receive entry with username:", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// If the error is of any other type, send a 500 status
		log.Println("Error of any other type", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Compare the stored hashed password, with the hashed version of the password that was received
	if err = bcrypt.CompareHashAndPassword([]byte(storedCreds.Password), []byte(creds.Password)); err != nil {
		// If the two passwords don't match, return a 401 status
		log.Println("Failed to authorize user:", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	fmt.Println("User has been logged in.")
	//	Create a new random session token
	sessionToken := uuid.NewV4().String()
	// Set the token in the cache, along with the user whom it represents
	// The token has an expiry time of 30 seconds
	_, err = Cache.Do("SETEX", sessionToken, "30", creds.Username)
	if err != nil {
		// If there is an error in setting the cache, return an internal server error
		log.Println("Error in setting the cache:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Finally, we set the client cookie for "session_token" as the session token we just generated
	// we also set an expiry time of 3600 seconds, the same as the cache
	cookie := http.Cookie{Name: "session_token", Value: sessionToken, Expires: time.Now().Add(30 * time.Second), HttpOnly: true, MaxAge: 50000}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
