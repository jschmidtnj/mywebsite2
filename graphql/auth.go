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
  "time"
  "strconv"
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
	passwordhashed, err := bcrypt.GenerateFromPassword([]byte(registerdata["password"].(string)), 14)
	if (err != nil) {
		handleErrorAuth("error hashing password: " + err.Error(), http.StatusBadRequest, response)
		return
  }
  email := registerdata["email"].(string)
  emailres, err := SendEmailVerification(email)
  if (err != nil) {
    handleErrorAuth("error sending email verification: " + err.Error(), http.StatusBadRequest, response)
		return
  } else if (emailres.StatusCode != 202) {
    handleErrorAuth("error sending email verification: got status code " + strconv.Itoa(emailres.StatusCode), http.StatusBadRequest, response)
		return
  }
	res, err := UserCollection.InsertOne(CTX, bson.M{
		"email": email,
		"password": string(passwordhashed),
		"emailverified": false,
		"type": "user",
	})
	if (err != nil) {
		handleErrorAuth("error inserting user to database: " + err.Error(), http.StatusBadRequest, response)
		return
	}
	id := res.InsertedID.(primitive.ObjectID).Hex()
	Logger.Info("User register",
		zap.String("id", id),
		zap.String("email", registerdata["email"].(string)),
	)
	response.Header().Set("content-type", "application/json")
	response.Write([]byte(`{"message":"please check email for verification"}`))
}

func Login(response http.ResponseWriter, request *http.Request) {
	if (request.Method != http.MethodPut) {
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
		handleErrorAuth("error finding user: " + err.Error(), http.StatusUnauthorized, response)
		return
	}
  defer cursor.Close(CTX)
  for cursor.Next(CTX) {
    userDataPrimitive := &bson.D{}
    err = cursor.Decode(userDataPrimitive)
    if (err != nil) {
      handleErrorAuth("error decoding user data: " + err.Error(), http.StatusBadRequest, response)
      return
    }
    userData := userDataPrimitive.Map()
    if (!(userData["emailverified"] != nil && userData["emailverified"].(bool))) {
      handleErrorAuth("email not verified", http.StatusUnauthorized, response)
      return
    }
    err = bcrypt.CompareHashAndPassword([]byte(userData["password"].(string)), []byte(logindata["password"].(string)))
    if (err != nil) {
      handleErrorAuth("invalid password: " + err.Error(), http.StatusUnauthorized, response)
      return
    }
    id := userData["_id"].(primitive.ObjectID).Hex()
    expirationTime := time.Now().Add(time.Duration(TokenExpiration) * time.Hour)
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
      "id": id,
      "email": userData["email"].(string),
      "type": userData["type"].(string),
      "StandardClaims": jwt.StandardClaims{
        ExpiresAt: expirationTime.Unix(),
        Issuer: "Joshua Schmidt",
      },
    })
    tokenString, err := token.SignedString(JwtSecret)
    if (err != nil) {
      handleErrorAuth("error creating token: " + err.Error(), http.StatusBadRequest, response)
      return
    }
    Logger.Info("User login",
      zap.String("id", id),
    )
    response.Header().Set("content-type", "application/json")
    response.Write([]byte(`{ "token": "` + tokenString + `" }`))
    break
  }
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
