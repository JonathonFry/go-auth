package main

import jwt "github.com/dgrijalva/jwt-go"

type registerPayload struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type user struct {
	Username string `json:"username"`
	Email    string `json:"email_address"`
	Password string `json:"password"`
}

type AuthClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type AuthToken struct {
	Token string `json:token`
}
