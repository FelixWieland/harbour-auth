package harbourauth

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"

	"golang.org/x/crypto/bcrypt"

	jwt "github.com/dgrijalva/jwt-go"
)

func register(w http.ResponseWriter, r *http.Request) {
	var user userCredentials

	//Get Username and Password
	user.Username = r.FormValue("username")
	user.Password = r.FormValue("password")

	if len(user.Username) == 0 || len(user.Password) == 0 {
		//errors
	}

	userUUID, _ := uuid.NewUUID()

	//Insert a New User to userauth
	_, err := db.Exec(sqlQuery("INSERT INTO harbour_userauth VALUES(?, ?, ?, CURRENT_TIMESTAMP)").prep(userUUID.String(), user.Username, hashAndSalt(user.Password)))
	if err != nil {
		log.Printf(err.Error())
		return
	}
	//then insert the userid into userdata
	_, err = db.Exec(sqlQuery("INSERT INTO harbour_userdata(userid) VALUES(?)").prep(userUUID.String()))
	if err != nil {
		log.Printf(err.Error())
		return
	}
	//then insert the userid into settings
	_, err = db.Exec(sqlQuery("INSERT INTO harbour_usersettings(userid) VALUES(?)").prep(userUUID.String()))
	if err != nil {
		log.Printf(err.Error())
		return
	}

	bv, err := json.Marshal(struct {
		Type     string `json:"type"`
		Username string `json:"username"`
	}{
		"Data",
		user.Username,
	})
	if err == nil {
		w.Write(bv)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}

}

func login(w http.ResponseWriter, r *http.Request) {
	var user userCredentials

	//Get Username and Password
	user.Username = r.FormValue("username")
	user.Password = r.FormValue("password")

	if len(user.Username) == 0 || len(user.Password) == 0 {
		//errors
	}

	rows, err := db.Query(sqlQuery("SELECT * FROM harbour_userauth WHERE username=?").prep(user.Username))
	if err != nil {
		log.Printf("Error in Select")
	}

	//Create Harbour Claims
	claims := HarbourClaims{}
	foundResults := false

	for rows.Next() {
		var userid string
		var username string
		var password string
		var created string

		err = rows.Scan(&userid, &username, &password, &created)

		foundResults = true

		if bcrypt.CompareHashAndPassword([]byte(password), []byte(user.Password)) == nil {
			//LOGIN SUCCESS
			//Set Harbour Claims
			claims = HarbourClaims{
				userid,
				username,
				jwt.StandardClaims{
					//ExpiresAt: 15000,
					Issuer: "login",
				},
			}
			break
		}
		w.WriteHeader(http.StatusForbidden)
		apiError(w, r, newErrResponseLoginFailed())
		return
	}

	if !foundResults {
		//no such user
		w.WriteHeader(http.StatusForbidden)
		apiError(w, r, newErrResponseLoginFailed())
		return
	}

	/*
		//validate user credentials
		if strings.ToLower(user.Username) != "alexcons" {
			if user.Password != "kappa123" {

				fmt.Println("Error logging in")
				fmt.Fprint(w, "Invalid credentials")
				return
			}
		}*/

	tokenString, _ := claims.Encode(signKey)

	//log.Printf("%v", tokenString)

	bv, err := json.Marshal(struct {
		Type string `json:"type"`
		JWT  string `json:"jwt"`
	}{
		"Data",
		tokenString,
	})
	w.Write(bv)
}

func decode(w http.ResponseWriter, r *http.Request) {

	jwt := r.FormValue("jwt")

	if len(jwt) == 0 {
		w.Write([]byte("Not all Keys are satisfied"))
	}

	claims, err := HarbourJWT(jwt).Decode(signKey)
	if err != nil {
		log.Printf("%v", err)
		return
	}
	//w.Write([]byte(claims))
	log.Printf("%v", claims)
}
