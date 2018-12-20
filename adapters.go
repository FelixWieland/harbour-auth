package harbourauth

import (
	"log"
	"net/http"

	"github.com/corneldamian/httpway"
)

func incomingConnection(w http.ResponseWriter, r *http.Request) {
	ctx := httpway.GetContext(r)
	log.Printf("Incoming connection from %v", r.RemoteAddr)

	w.Header().Set("Content-Type", "application/json")

	ctx.Next(w, r)
}

func checkAuth(w http.ResponseWriter, r *http.Request) {
	ctx := httpway.GetContext(r)
	w.Header().Set("Content-Type", "application/json")
	jwt := r.FormValue("jwt")
	if len(jwt) == 0 {
		//no jwt provided
	}
	claims, err := HarbourJWT(jwt).Decode(signKey, secret)
	if err != nil {
		//notloggedin
		w.Write([]byte("Invalid or no JWT"))
	} else {
		ctx.Set("userid", claims.UserID)
		ctx.Set("username", claims.Username)
		ctx.Set("issue", claims.Issuer)
		ctx.Next(w, r)
	}
}
