package main

import (
	"errors"
	jwt "github.com/dgrijalva/jwt-go"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"golang.org/x/crypto/bcrypt"
	"go.mongodb.org/mongo-driver/bson"
	"encoding/json"
	"net/http"
)

func handleErrorAuth(err error, response http.ResponseWriter) {
	response.Header().Set("content-type", "application/json")
	response.Write([]byte(`{ "error": "` + err.Error() + `" }`))
	return
}

func Register(response http.ResponseWriter, request *http.Request) {
	var user User
	err := json.NewDecoder(request.Body).Decode(&user)
	if (err != nil) {
		handleErrorAuth(err, response)
	}
	passwordbytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if (err != nil) {
		handleErrorAuth(err, response)
	}
	res, err := UserCollection.InsertOne(CTX, bson.M{
		"email": user.Email,
		"password": string(passwordbytes),
	})
	id := res.InsertedID
	fmt.Println(id)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
	})
	tokenString, err := token.SignedString(JwtSecret)
	if (err != nil) {
		handleErrorAuth(err, response)
	}
	response.Header().Set("content-type", "application/json")
	response.Write([]byte(`{ "token": "` + tokenString + `" }`))
}

func Login(response http.ResponseWriter, request *http.Request) {
	var userInput User
	err := json.NewDecoder(request.Body).Decode(&userInput)
	if (err != nil) {
		handleErrorAuth(err, response)
	}
	var userDb User
	err = UserCollection.FindOne(CTX, bson.M{"email": userInput.Email}).Decode(&userDb)
	if (err != nil) {
		handleErrorAuth(err, response)
	}
	err = bcrypt.CompareHashAndPassword([]byte(userDb.Password), []byte(userInput.Password))
	if (err != nil) {
		handleErrorAuth(err, response)
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": userDb.Email,
	})
	tokenString, err := token.SignedString(JwtSecret)
	if (err != nil) {
		handleErrorAuth(err, response)
	}
	response.Header().Set("content-type", "application/json")
	response.Write([]byte(`{ "token": "` + tokenString + `" }`))
}

// ValidateLoggedIn validates JWT token to confirm login
func ValidateLoggedIn(t string) (interface{}, error) {
	if t == "" {
		return nil, errors.New("Authorization token must be present")
	}
	token, _ := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error")
		}
		return JwtSecret, nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		var decodedToken interface{}
		mapstructure.Decode(claims, &decodedToken)
		return decodedToken, nil
	} else {
		return nil, errors.New("Invalid authorization token")
	}
}
