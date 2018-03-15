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

type authToken struct {
	Token string `json:"token"`
}

type userResponse struct {
	*user
	authToken
}

type errorResponse struct {
	Error string `json:"error"`
}
