package auth

type AuthResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refresh"`
	User         User   `json:"user"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type RefreshRequest struct {
	RefreshToken         string `json:"refresh"`
	LastValidAccessToken string `json:"accessToken"`
}

type RefreshResponse struct {
	AccessToken string `json:"accessToken"`
}

type ValidateRequest struct {
	AccessToken string `json:"accessToken"`
}

type ValidateResponse struct {
	UserId int `json:"user_id"`
}
