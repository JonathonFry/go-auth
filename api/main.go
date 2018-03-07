package main

import (
	"database/sql"
	"encoding/json"
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

var templates = template.Must(template.ParseFiles("tmpl/index.html", "tmpl/login.html", "tmpl/register.html"))

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

	s := http.StripPrefix("/static/", http.FileServer(http.Dir("./static/")))
	r.PathPrefix("/static/").Handler(s)

	// amw := authenticationMiddleware{}
	// r.Use(amw.Middleware)
	corsObj := handlers.AllowedOrigins([]string{"*"})

	loggedRouter := handlers.CORS(corsObj)(handlers.LoggingHandler(os.Stdout, r))

	srv := &http.Server{
		Handler:      loggedRouter,
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	_, user := isLoggedIn(r)

	err := templates.ExecuteTemplate(w, "index.html", user)
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

	json, marshalErr := json.Marshal(users)

	if marshalErr != nil {
		http.Error(w, marshalErr.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

func registerGetHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "register.html", nil)
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
		setCookie(w, register.Username)
		fmt.Fprintf(w, "User registered succesfully %d", userid)
	}
}

func loginGetHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "login.html", nil)
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

	setCookie(w, user.Username)

	fmt.Fprintf(w, "User logged in succesfully %s", user.Username)
}

func setCookie(w http.ResponseWriter, value string) {
	cookie := &http.Cookie{
		Name:  "session",
		Value: uuid.NewV4().String(),
	}
	storageMutex.Lock()
	sessionStore[cookie.Value] = value
	storageMutex.Unlock()

	http.SetCookie(w, cookie)
}

func isLoggedIn(r *http.Request) (bool, *user) {
	cookie, err := r.Cookie("session")
	if err != nil {
		return false, nil
	}
	if cookie != nil {
		storageMutex.RLock()
		username, exists := sessionStore[cookie.Value]
		storageMutex.RUnlock()
		user, err := getUser(username, db)
		if err != nil {
			return false, nil
		}
		return exists, user
	}
	return false, nil
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

		if present || strings.Contains(r.URL.Path, "login") || strings.Contains(r.URL.Path, "register") || strings.Contains(r.URL.Path, "static") {
			next.ServeHTTP(w, r)
		} else {
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}
	})
}
