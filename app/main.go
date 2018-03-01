package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

var templates = template.Must(template.ParseFiles("tmpl/index.html.tmpl", "tmpl/login.html.tmpl", "tmpl/register.html.tmpl"))

var db *sql.DB

func main() {
	db = awaitDb()

	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/users", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		users, err := getUsers(db)
		if err != nil {
			fmt.Fprintf(w, "Error retrieving users: %s", err)
			return
		}
		fmt.Fprintf(w, "Users: %s", users)
	}))
	r.HandleFunc("/login", loginGetHandler).Methods("GET")
	r.HandleFunc("/login", loginPostHandler).Methods("POST")
	r.HandleFunc("/register", registerGetHandler).Methods("GET")
	r.HandleFunc("/register", registerPostHandler).Methods("POST")

	loggedRouter := handlers.LoggingHandler(os.Stdout, r)

	srv := &http.Server{
		Handler:      loggedRouter,
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "index.html.tmpl", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func registerGetHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "register.html.tmpl", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func registerPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	register := new(registerPayload)
	decoder := schema.NewDecoder()
	err := decoder.Decode(register, r.PostForm)
	if err != nil {
		log.Fatal("Error parsing registration data")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(register.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Error hashing password")
	}

	userid, err := insertUser(register, hash, db)
	if err != nil {
		fmt.Fprintf(w, "Error registering user %s", err)
	} else {
		fmt.Fprintf(w, "User registered succesfully %d", userid)
	}
}

func loginGetHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "login.html.tmpl", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func loginPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	login := new(loginPayload)
	decoder := schema.NewDecoder()
	err := decoder.Decode(login, r.PostForm)
	if err != nil {
		log.Fatal("Error parsing login data")
	}

	user, err := getUser(login, db)

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password)); err != nil {
		fmt.Fprintf(w, "Invalid password")
		return
	}

	fmt.Fprintf(w, "User logged in succesfully %s", user.Username)
}

/*
Users DB

- user_id = auto generated USER ID
- email = user entered email address
- password = bcrypt password hash

Register
API endpoint = POST /register
{
	"email" : "jane@doe.com",
	"password" : "alligator"
}
	-- register user with given credentials


Login
API endpoint = POST /login
{
	"email" : "jane@doe.com",
	"password" : "alligator"
}
	-- authenticates user with given credentials

Logout
API endpoint = POST /logout
	-- Clear session cookie

Sessions
	Implement cookie based session storage
	use redis DB to persist sessions
*/
