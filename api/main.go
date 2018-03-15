package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

var secret []byte
var db *sql.DB

func main() {
	db = awaitDb()

	secret = []byte("secret")
	r := mux.NewRouter()
	r.Handle("/user", authMiddleware(userHandler)).Methods("GET")
	r.Handle("/users", authMiddleware(usersHandler)).Methods("GET")
	r.HandleFunc("/login", loginHandler).Methods("POST")
	r.HandleFunc("/register", registerHandler).Methods("POST")

	headersOk := handlers.AllowedHeaders([]string{"Authorization", "Content-Type"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	loggedRouter := handlers.CORS(headersOk, originsOk, methodsOk)(handlers.LoggingHandler(os.Stdout, r))

	srv := &http.Server{
		Handler:      loggedRouter,
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		validToken := validToken(r)
		if !validToken {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

var usersHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	users, err := getUsers(db)
	if err != nil {
		fmt.Fprintf(w, "Error retrieving users: %s", err)
		return
	}

	returnJSON(w, users)
})

var userHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	user, err := loggedInUser(r)

	if err != nil {
		fmt.Fprintf(w, "Error retrieving user: %s", err)
		return
	}

	returnJSON(w, user)
})

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

	_, err = insertUser(register, hash, db)
	if err != nil {
		fmt.Fprintf(w, "Error registering user %s", err)
	} else {
		user := &user{Username: register.Username, Password: string(hash), Email: register.Email}
		token, err := createToken(user)

		if err != nil {
			fmt.Fprintf(w, "Error generating token %s", err)
			return
		}

		returnJSON(w, userResponse{user: user, authToken: authToken{Token: token}})
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

	token, err := createToken(user)

	if err != nil {
		fmt.Fprintf(w, "Error generating token %s", err)
		return
	}

	returnJSON(w, userResponse{user: user, authToken: authToken{Token: token}})
}

func createToken(user *user) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
	})

	tokenString, err := token.SignedString(secret)
	return tokenString, err
}

func token(tokenString string) (*jwt.Token, error) {
	tokenString = strings.Replace(tokenString, "Bearer ", "", -1)
	token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err == nil && token.Valid {
		return token, err
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			fmt.Println("That's not even a token")
			return nil, err
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			// Token is either expired or not active yet
			fmt.Println("Timing is everything")
			return nil, err
		} else {
			fmt.Println("Couldn't handle this token:", err)
			return nil, err
		}
	} else {
		fmt.Println("Couldn't handle this token:", err)
		return nil, err
	}
}

func validToken(r *http.Request) bool {
	tokenString := r.Header.Get("Authorization")

	if len(tokenString) > 0 {
		_, err := token(tokenString)
		if err != nil {
			return false
		}
		return true
	}

	return false
}

func loggedInUser(r *http.Request) (*user, error) {
	tokenString := r.Header.Get("Authorization")

	if len(tokenString) > 0 {
		token, err := token(tokenString)
		if err != nil {
			return nil, err
		}
		c := token.Claims.(*jwt.MapClaims)
		username := (*c)["username"].(string)

		log.Println(username)

		user, err := getUser(username, db)

		if err != nil {
			return nil, err
		}

		return user, nil
	}

	// No token sent
	return nil, errors.New("no user found")
}

func returnJSON(w http.ResponseWriter, data interface{}) {
	json, marshalErr := json.Marshal(data)

	if marshalErr != nil {
		http.Error(w, marshalErr.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}
