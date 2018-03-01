package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	_ "github.com/lib/pq"
	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

var templates = template.Must(template.ParseFiles("tmpl/index.html.tmpl", "tmpl/login.html.tmpl", "tmpl/register.html.tmpl"))

var sessionStore map[string]string //TODO - Replace with redis DB
var storageMutex sync.RWMutex

var db *sql.DB

func main() {
	db = awaitDb()

	sessionStore = make(map[string]string)

	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler).Methods("GET")
	r.HandleFunc("/users", usersHandler).Methods("GET")
	r.HandleFunc("/login", loginGetHandler).Methods("GET")
	r.HandleFunc("/login", loginPostHandler).Methods("POST")
	r.HandleFunc("/register", registerGetHandler).Methods("GET")
	r.HandleFunc("/register", registerPostHandler).Methods("POST")

	amw := authenticationMiddleware{}
	r.Use(amw.Middleware)

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

func usersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := getUsers(db)
	if err != nil {
		fmt.Fprintf(w, "Error retrieving users: %s", err)
		return
	}
	fmt.Fprintf(w, "Users: %s", users)
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

	user, err := getUser(login.Username, db)

	if err != nil {
		fmt.Fprintf(w, "User doesn't exist")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password)); err != nil {
		fmt.Fprintf(w, "Invalid password")
		return
	}

	cookie := &http.Cookie{
		Name:  "session",
		Value: uuid.NewV4().String(),
	}
	storageMutex.Lock()
	sessionStore[cookie.Value] = user.Username
	storageMutex.Unlock()

	http.SetCookie(w, cookie)

	fmt.Fprintf(w, "User logged in succesfully %s", user.Username)
}

type authenticationMiddleware struct {
}

func (amw *authenticationMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var present bool

		cookie, err := r.Cookie("session")
		if err != nil {
			present = false
		}

		if cookie != nil {
			storageMutex.RLock()
			// Retrieve user from cache
			_, present = sessionStore[cookie.Value]
			storageMutex.RUnlock()
		} else {
			present = false
		}

		if present || strings.Contains(r.URL.Path, "login") || strings.Contains(r.URL.Path, "register") {
			next.ServeHTTP(w, r)
		} else {
			http.Redirect(w, r, "/login", http.StatusPermanentRedirect)
			return
		}
	})
}
