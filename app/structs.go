package main

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
