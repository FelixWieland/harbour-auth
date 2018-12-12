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
