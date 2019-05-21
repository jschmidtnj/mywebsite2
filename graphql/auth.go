package main

import (
	"errors"
	jwt "github.com/dgrijalva/jwt-go"
	"go.uber.org/zap"
	"github.com/mitchellh/mapstructure"
	"golang.org/x/crypto/bcrypt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"encoding/json"
	"net/http"
	"io/ioutil"
)

func handleErrorAuth(message string, statuscode int, response http.ResponseWriter) {
	// Logger.Error(message)
	response.Header().Set("content-type", "application/json")
	response.WriteHeader(statuscode)
	response.Write([]byte(`{"message":"` + message + `"}`))
}

func Register(response http.ResponseWriter, request *http.Request) {
	if (request.Method != http.MethodPost) {
		handleErrorAuth("register http method not POST", http.StatusBadRequest, response)
		return
	}
	var registerdata map[string]interface{}
	body, err := ioutil.ReadAll(request.Body)
	if (err != nil) {
		handleErrorAuth("error getting request body: " + err.Error(), http.StatusBadRequest, response)
		return
	}
    err = json.Unmarshal(body, &registerdata)
	if (err != nil) {
		handleErrorAuth("error parsing request body: " + err.Error(), http.StatusBadRequest, response)
		return
	}
	if (!(registerdata["password"] != nil && registerdata["email"] != nil)) {
		handleErrorAuth("no email or password provided", http.StatusBadRequest, response)
		return
	}
	countemail, err := UserCollection.CountDocuments(CTX, bson.M{"email": registerdata["email"]})
	if (err != nil) {
		handleErrorAuth("error counting users with same email: " + err.Error(), http.StatusBadRequest, response)
		return
	}
	if (countemail > 0) {
		handleErrorAuth("email is already taken", http.StatusBadRequest, response)
		return
	}
	passwordbytes, err := bcrypt.GenerateFromPassword(registerdata["password"].([]byte), 14)
	if (err != nil) {
		handleErrorAuth("error hashing password: " + err.Error(), http.StatusBadRequest, response)
		return
	}
	res, err := UserCollection.InsertOne(CTX, bson.M{
		"email": registerdata["email"].(string),
		"password": string(passwordbytes),
		"emailverified": false,
		"type": "user",
	})
	if (err != nil) {
		handleErrorAuth("error inserting user to database: " + err.Error(), http.StatusBadRequest, response)
		return
	}
	id := res.InsertedID
	Logger.Info("User register",
		zap.Int("id", id.(int)),
		zap.String("email", registerdata["email"].(string)),
	)
	response.Header().Set("content-type", "application/json")
	response.Write([]byte(`{"message":"please check email for verification"}`))
}

func Login(response http.ResponseWriter, request *http.Request) {
	if (request.Method != http.MethodPost) {
		handleErrorAuth("login http method not PUT", http.StatusBadRequest, response)
		return
	}
	var logindata map[string]interface{}
	body, err := ioutil.ReadAll(request.Body)
	if (err != nil) {
		handleErrorAuth("error getting request body: " + err.Error(), http.StatusBadRequest, response)
		return
	}
    err = json.Unmarshal(body, &logindata)
	if (err != nil) {
		handleErrorAuth("error parsing request body: " + err.Error(), http.StatusBadRequest, response)
		return
	}
	if (logindata["email"] == nil || logindata["password"] == nil) {
		handleErrorAuth("no email or password provided", http.StatusBadRequest, response)
		return
	}
	cursor, err := UserCollection.Find(CTX, bson.M{"email": logindata["email"]})
	if (err != nil) {
		handleErrorAuth("error finding user: " + err.Error(), http.StatusBadRequest, response)
		return
	}
	defer cursor.Close(CTX)
	userDataPrimitive := &bson.D{}
	err = cursor.Decode(userDataPrimitive)
	if (err != nil) {
		handleErrorAuth("error decoding user data: " + err.Error(), http.StatusBadRequest, response)
		return
	}
	userData := userDataPrimitive.Map()
	if (!(userData["emailverified"] != nil && userData["emailverified"].(bool))) {
		handleErrorAuth("email not verified", http.StatusBadRequest, response)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(userData["password"].(string)), []byte(logindata["password"].(string)))
	if (err != nil) {
		handleErrorAuth("error comparing password: " + err.Error(), http.StatusBadRequest, response)
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id": userData["_id"].(primitive.ObjectID).Hex(),
		"email": userData["email"].(string),
		"type": userData["type"].(string),
	})
	tokenString, err := token.SignedString(JwtSecret)
	if (err != nil) {
		handleErrorAuth("error creating token: " + err.Error(), http.StatusBadRequest, response)
		return
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
		if _, success := token.Method.(*jwt.SigningMethodHMAC); !success {
			return nil, errors.New("There was an error")
		}
		return JwtSecret, nil
	})
	if claims, success := token.Claims.(jwt.MapClaims); success && token.Valid {
		var decodedToken interface{}
		mapstructure.Decode(claims, &decodedToken)
		return decodedToken, nil
	} else {
		return nil, errors.New("Invalid authorization token")
	}
}
