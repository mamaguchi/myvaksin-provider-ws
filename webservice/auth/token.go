package auth

import (
	"log"
	"fmt"
	"time"
	"github.com/dgrijalva/jwt-go"
)

// For HMAC signing method, the secret key can be any []byte. It is recommended to generate
// a key using crypto/rand or something equivalent. You need the same secret key for signing
// and validating.
var hmacSecret = []byte(`patricksecretkey`)

func NewTokenHMAC(userId string) (tokenString string, err error) {
	now := time.Now()
	expiredAt := now.Add(time.Hour * 1)

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId,
		"nbf": now.Unix(),
		"exp": expiredAt.Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err = token.SignedString(hmacSecret)

	return
}

func VerifyTokenHMAC(tokenString string) (bool) {
	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		
		// hmacSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return hmacSecret, nil
	})
	if err != nil {
		log.Print(err)
		return false
	}	

	return token.Valid
}










