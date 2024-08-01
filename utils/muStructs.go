package utils

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Status  string `json:"status"`
	Reason  string `json:"reason"`
	Context struct {
		SessionToken string `json:"session_token"`
	} `json:"context"`
}

type LogoutResponse struct {
	Status  string   `json:"status"`
	Reason  string   `json:"reason"`
	Context struct{} `json:"context"`
}
