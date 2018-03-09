package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

var sessionStore map[string]string //TODO - Replace with redis DB
var storageMutex sync.RWMutex

var db *sql.DB

func main() {
	db = awaitDb()

	sessionStore = make(map[string]string)

	r := mux.NewRouter()
	r.HandleFunc("/users", usersHandler).Methods("GET")
	r.HandleFunc("/login", loginHandler).Methods("POST")
	r.HandleFunc("/register", registerHandler).Methods("POST")

	amw := authenticationMiddleware{}
	r.Use(amw.Middleware)
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

func registerHandler(w http.ResponseWriter, r *http.Request) {
	register := new(registerPayload)
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(register)
	if err != nil {
		fmt.Fprintf(w, "Error parsing registration data")
		return
	}

	if len(register.Email) == 0 || len(register.Username) == 0 || len(register.Password) == 0 {
		fmt.Fprintf(w, "Error required fields are missing")
		return
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

func loginHandler(w http.ResponseWriter, r *http.Request) {
	login := new(loginPayload)
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(login)
	if err != nil {
		fmt.Fprintf(w, "Error parsing login data")
		return
	}

	log.Println(login)
	if len(login.Username) == 0 || len(login.Password) == 0 {
		fmt.Fprintf(w, "Error required fields are missing")
		return
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
			http.Error(w, "Unauthorised", http.StatusUnauthorized)
			return
		}
	})
}
