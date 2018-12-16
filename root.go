package harbourauth

import (
	"crypto/rsa"
	"database/sql"
	"net/http"

	"github.com/corneldamian/httpway"
	"github.com/rs/cors"
)

var signKey *rsa.PrivateKey
var server *httpway.Server
var db *sql.DB

const (
	privKeyPath = "keys/app.rsa" //openssl genrsa -out app.rsa 1024
)

//Start starts the Authentication Service
func Start() {
	signKey, _ = LoadAsPrivateRSAKey(privKeyPath)

	credentials := loadCredentials("../sqlAuth.json")
	if ldb, err := connectToDB(credentials.toString()); err == nil {
		db = ldb
		defer db.Close()
	} else {
		println("Cant connect to Database")
	}

	router := httpway.New()
	public := router.Middleware(incomingConnection)

	/*PUBLIC ROUTES*/
	public.POST("/login", login)
	public.POST("/decode", decode)
	public.POST("/register", register)

	handler := cors.Default().Handler(router) //enable access from all origins
	http.ListenAndServe(":5000", handler)

	server = httpway.NewServer(nil)
	server.Addr = ":5000"
	server.Handler = router

	server.Start()

}
