package harbourauth

import (
	"encoding/json"
	"net/http"
)

//response for an unauthenticated user
type errorResponse struct {
	Type        string `json:"type"`
	Code        int    `json:"code"`
	Message     string `json:"message"`
	Description string `json:"description"`
}

type infoResponse struct {
	Type        string `json:"type"`
	GeneratedAt string `json:"generatedAt"`
	Code        int    `json:"code"`
	Message     string `json:"message"`
	Description string `json:"description"`
}

//stdResponseLoginFailed
func newErrResponseLoginFailed() errorResponse {
	return errorResponse{
		"error",
		1,
		"invalidCredentials",
		"Your Username or Password was incorrect. Please try again",
	}
}

//stdResponseInvalitJWT
func newErrResponseInvalidJWT() errorResponse {
	return errorResponse{
		"error",
		2,
		"invalidJWT",
		"Your JWT is invalid.",
	}
}

/*
//stdResponseCreatedJWT
func newInfoResponseCreatedJWT() infoResponse {
	return infoResponse{
		"info",
		time.Now().String(),
		1,
		"createdJWT",
		"You are alredy logged in",
	}
}*/

//API Error Response
func apiError(w http.ResponseWriter, r *http.Request, pl errorResponse) {

	js, err := json.Marshal(pl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(js)
}

//API Info Response
func apiInfo(w http.ResponseWriter, r *http.Request, pl infoResponse) {

	js, err := json.Marshal(pl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(js)
}
