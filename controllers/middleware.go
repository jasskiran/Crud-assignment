package controllers

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const jwtsecret = "JWT_SECRET"



// generates jwt token
func GenerateJWT(userId int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":   userId,
		"exp":       time.Now().Add(time.Hour * 24 * 7).Unix(),
	})
	tokenString, err := token.SignedString([]byte(jwtsecret))
	return tokenString, err
}

type Token struct {
	//AuthUUID string `json:"auth_uuid"`
	UserId   int    `json:"user_id"`
	jwt.StandardClaims
}

func AuthRequired(handlerFunc func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")

		if auth == "" { //Token is missing, returns with error code 403 Unauthorized
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		splitted := strings.Split(auth, " ") //The token normally comes in format `Bearer {token-body}`, we check if the retrieved token matched this requirement
		if len(splitted) != 2 {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		tokenPart := splitted[1] //Grab the token part, what we are truly interested in

		var tk Token

		token, err := jwt.ParseWithClaims(tokenPart, &tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtsecret), nil
		})


		if err != nil { //Malformed token, returns with http code 403 as usual
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		if !token.Valid { //Token is invalid, maybe not signed on this server
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		//create a new request context containing the authenticated user
		ctxWithUser := context.WithValue(r.Context(),"userId" , tk)
		//create a new request using that new context
		rWithUser := r.WithContext(ctxWithUser)
		//call the real handler, passing the new request
		handlerFunc(w, rWithUser)
	}
}


