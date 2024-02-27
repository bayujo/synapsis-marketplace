package user

type User struct {
    ID       int    `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
    Password string `json:"password"`
	
}

type TokenResponse struct {
    Token string `json:"token"`
}