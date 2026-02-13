package types

type ErrorResponse struct {
	Error string `json:"error"`
}

type LoginResponse struct {
	Token string `json:"token"`
}
