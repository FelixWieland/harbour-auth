package harbourauth

import (
	"crypto/rsa"
	"errors"
	"io/ioutil"
	"log"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

//HarbourJWT Custom Type to Decode and Encode a JWT
type HarbourJWT string

var verifyKey *rsa.PublicKey

type userCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

//HarbourClaims Claims
type HarbourClaims struct {
	UserID   string `json:"UserID"`
	Username string `json:"Username"`
	Secret   string `json:"secret"`
	jwt.StandardClaims
}

//LoadAsPrivateRSAKey returns a parsed RSA Key from path
func LoadAsPrivateRSAKey(path string) (*rsa.PrivateKey, error) {
	signKey := &rsa.PrivateKey{}
	lsignKey, err := ioutil.ReadFile(privKeyPath)
	if err != nil {
		return signKey, err
	}
	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(lsignKey)
	if err != nil {
		return signKey, err
	}
	return signKey, nil
}

//Decode HarbourJWT
func (HarbourJWT HarbourJWT) Decode(key *rsa.PrivateKey, pSecret string) (HarbourClaims, error) {
	tokenString := HarbourJWT
	encodedClaims := &HarbourClaims{}
	_, _, err := new(jwt.Parser).ParseUnverified(string(tokenString), encodedClaims)
	if err == nil {
		if encodedClaims.Valid() == nil && encodedClaims.Secret == pSecret {
			return *encodedClaims, nil
		}
		err = errors.New("provided jwt is not valid")
	}
	return HarbourClaims{}, err
}

//Encode to encode Claims to a JWT
func (HarbourClaims HarbourClaims) Encode(key *rsa.PrivateKey) (string, error) {
	signer := jwt.NewWithClaims(jwt.SigningMethodRS256, HarbourClaims)
	tokenString, err := signer.SignedString(signKey)

	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func hashAndSalt(pwd string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}
