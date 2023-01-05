package auth

type User struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Hash     string `json:"hash"`
	Email    string `json:"email"`
}
