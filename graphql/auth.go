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

func Register(response http.ResponseWriter, request *http.Request) {
	if !ManageCors(response, request) {
		return
	}
	if (request.Method != http.MethodPost) {
		handleError("register http method not POST", http.StatusBadRequest, response)
		return
	}
	var registerdata map[string]interface{}
	body, err := ioutil.ReadAll(request.Body)
	if (err != nil) {
		handleError("error getting request body: " + err.Error(), http.StatusBadRequest, response)
		return
	}
  err = json.Unmarshal(body, &registerdata)
	if (err != nil) {
		handleError("error parsing request body: " + err.Error(), http.StatusBadRequest, response)
		return
	}
	if (!(registerdata["password"] != nil && registerdata["email"] != nil)) {
		handleError("no email or password provided", http.StatusBadRequest, response)
		return
	}
	countemail, err := UserCollection.CountDocuments(CTX, bson.M{"email": registerdata["email"]})
	if (err != nil) {
		handleError("error counting users with same email: " + err.Error(), http.StatusBadRequest, response)
		return
	}
	if (countemail > 0) {
		handleError("email is already taken", http.StatusBadRequest, response)
		return
  }
	passwordhashed, err := bcrypt.GenerateFromPassword([]byte(registerdata["password"].(string)), 14)
	if (err != nil) {
		handleError("error hashing password: " + err.Error(), http.StatusBadRequest, response)
		return
  }
  email := registerdata["email"].(string)
  emailres, err := SendEmailVerification(email)
  if (err != nil) {
    handleError("error sending email verification: " + err.Error(), http.StatusBadRequest, response)
		return
  } else if (emailres.StatusCode != 202) {
    handleError("error sending email verification: got status code " + strconv.Itoa(emailres.StatusCode), http.StatusBadRequest, response)
		return
  }
	res, err := UserCollection.InsertOne(CTX, bson.M{
		"email": email,
		"password": string(passwordhashed),
		"emailverified": false,
		"type": "user",
	})
	if (err != nil) {
		handleError("error inserting user to database: " + err.Error(), http.StatusBadRequest, response)
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

func LoginEmailPassword(response http.ResponseWriter, request *http.Request) {
	if !ManageCors(response, request) {
		return
	}
	if (request.Method != http.MethodPut) {
		handleError("login http method not PUT", http.StatusBadRequest, response)
		return
	}
	var logindata map[string]interface{}
	body, err := ioutil.ReadAll(request.Body)
	if (err != nil) {
		handleError("error getting request body: " + err.Error(), http.StatusBadRequest, response)
		return
	}
  err = json.Unmarshal(body, &logindata)
	if (err != nil) {
		handleError("error parsing request body: " + err.Error(), http.StatusBadRequest, response)
		return
	}
	if (logindata["email"] == nil || logindata["password"] == nil) {
		handleError("no email or password provided", http.StatusBadRequest, response)
		return
	}
	cursor, err := UserCollection.Find(CTX, bson.M{"email": logindata["email"]})
	defer cursor.Close(CTX)
	if (err != nil) {
		handleError("error finding user: " + err.Error(), http.StatusUnauthorized, response)
		return
	}
  for cursor.Next(CTX) {
    userDataPrimitive := &bson.D{}
    err = cursor.Decode(userDataPrimitive)
    if (err != nil) {
      handleError("error decoding user data: " + err.Error(), http.StatusBadRequest, response)
      return
		}
		userData := userDataPrimitive.Map()
    if (!(userData["emailverified"] != nil && userData["emailverified"].(bool))) {
      handleError("email not verified", http.StatusUnauthorized, response)
      return
    }
    err = bcrypt.CompareHashAndPassword([]byte(userData["password"].(string)), []byte(logindata["password"].(string)))
    if (err != nil) {
      handleError("invalid password: " + err.Error(), http.StatusUnauthorized, response)
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
      handleError("error creating token: " + err.Error(), http.StatusBadRequest, response)
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

func LogoutEmailPassword(response http.ResponseWriter, request *http.Request) {
	if !ManageCors(response, request) {
		return
	}
	if (request.Method != http.MethodPut) {
		handleError("logout http method not PUT", http.StatusBadRequest, response)
		return
	}
	_, err := ValidateLoggedIn(GetAuthToken(request))
	if (err != nil) {
		handleError(err.Error(), http.StatusBadRequest, response)
		return
	}
	response.Header().Set("content-type", "application/json")
	response.Write([]byte(`{ "message": "successfully signed out" }`))
}

func VerifyEmail(response http.ResponseWriter, request *http.Request) {
	if !ManageCors(response, request) {
		return
	}
	if (request.Method != http.MethodPost) {
		handleError("verify http method not POST", http.StatusBadRequest, response)
		return
	}
	var verifydata map[string]interface{}
	body, err := ioutil.ReadAll(request.Body)
	if (err != nil) {
		handleError("error getting request body: " + err.Error(), http.StatusBadRequest, response)
		return
	}
  err = json.Unmarshal(body, &verifydata)
	if (err != nil) {
		handleError("error parsing request body: " + err.Error(), http.StatusBadRequest, response)
		return
	}
	if (verifydata["token"] == nil) {
		handleError("no token provided", http.StatusBadRequest, response)
		return
	}
	token, err := jwt.Parse(verifydata["token"].(string), func(token *jwt.Token) (interface{}, error) {
		if _, success := token.Method.(*jwt.SigningMethodHMAC); !success {
			return nil, errors.New("There was an error")
		}
		return JwtSecret, nil
	})
	if (err != nil) {
		handleError(err.Error(), http.StatusBadRequest, response)
		return
	}
	var decodedToken map[string]interface{}
	if claims, success := token.Claims.(jwt.MapClaims); success && token.Valid {
		mapstructure.Decode(claims, &decodedToken)
	} else {
		handleError("invalid token", http.StatusBadRequest, response)
		return
	}
	if (!(decodedToken["email"] != nil && decodedToken["verify"] != nil && decodedToken["verify"].(bool))) {
		handleError("token does not contian email or verify boolean", http.StatusBadRequest, response)
		return
	}
	cursor, err := UserCollection.Find(CTX, bson.M{"email": decodedToken["email"]})
	if (err != nil) {
		handleError("error finding user: " + err.Error(), http.StatusBadRequest, response)
		return
	}
  defer cursor.Close(CTX)
  for cursor.Next(CTX) {
    userDataPrimitive := &bson.D{}
    err = cursor.Decode(userDataPrimitive)
    if (err != nil) {
      handleError("error decoding user data: " + err.Error(), http.StatusBadRequest, response)
      return
    }
    userData := userDataPrimitive.Map()
    if (userData["emailverified"] != nil && !userData["emailverified"].(bool)) {
      handleError("email already verified", http.StatusBadRequest, response)
      return
		}
		var id string = userData["_id"].(string)
		_, err := UserCollection.UpdateOne(CTX, bson.M{
			"_id": id,
		}, bson.M{
			"emailverified": true,
		})
		if (err != nil) {
			handleError("error updating user in database: " + err.Error(), http.StatusBadRequest, response)
			return
		}
		Logger.Info("User email verify",
			zap.String("id", id),
			zap.String("email", userData["email"].(string)),
		)
    response.Header().Set("content-type", "application/json")
    response.Write([]byte(`{"message":"email successfully verified"}`))
    break
  }
}

func ResetPassword(response http.ResponseWriter, request *http.Request) {
	if !ManageCors(response, request) {
		return
	}
	if (request.Method != http.MethodPost) {
		handleError("reset http method not POST", http.StatusBadRequest, response)
		return
	}
	var resetdata map[string]interface{}
	body, err := ioutil.ReadAll(request.Body)
	if (err != nil) {
		handleError("error getting request body: " + err.Error(), http.StatusBadRequest, response)
		return
	}
  err = json.Unmarshal(body, &resetdata)
	if (err != nil) {
		handleError("error parsing request body: " + err.Error(), http.StatusBadRequest, response)
		return
	}
	if (resetdata["token"] == nil || resetdata["password"] == nil) {
		handleError("no token or new password provided", http.StatusBadRequest, response)
		return
	}
	token, err := jwt.Parse(resetdata["token"].(string), func(token *jwt.Token) (interface{}, error) {
		if _, success := token.Method.(*jwt.SigningMethodHMAC); !success {
			return nil, errors.New("There was an error")
		}
		return JwtSecret, nil
	})
	if (err != nil) {
		handleError(err.Error(), http.StatusBadRequest, response)
		return
	}
	var decodedToken map[string]interface{}
	if claims, success := token.Claims.(jwt.MapClaims); success && token.Valid {
		mapstructure.Decode(claims, &decodedToken)
	} else {
		handleError("invalid token", http.StatusBadRequest, response)
		return
	}
	if (!(decodedToken["email"] != nil && decodedToken["reset"] != nil && decodedToken["reset"].(bool))) {
		handleError("token does not contian email or reset boolean", http.StatusBadRequest, response)
		return
	}
	cursor, err := UserCollection.Find(CTX, bson.M{"email": decodedToken["email"]})
	defer cursor.Close(CTX)
	if (err != nil) {
		handleError("error finding user: " + err.Error(), http.StatusBadRequest, response)
		return
	}
  for cursor.Next(CTX) {
    userDataPrimitive := &bson.D{}
    err = cursor.Decode(userDataPrimitive)
    if (err != nil) {
      handleError("error decoding user data: " + err.Error(), http.StatusBadRequest, response)
      return
    }
    userData := userDataPrimitive.Map()
    if (!(userData["_id"] != nil && userData["emailverified"] != nil && userData["emailverified"].(bool))) {
      handleError("user id invalid or email not verified", http.StatusBadRequest, response)
      return
		}
		var id string = userData["_id"].(string)
		passwordhashed, err := bcrypt.GenerateFromPassword([]byte(resetdata["password"].(string)), 14)
		if (err != nil) {
			handleError("error hashing password: " + err.Error(), http.StatusBadRequest, response)
			return
		}
		_, err = UserCollection.UpdateOne(CTX, bson.M{
			"_id": id,
		}, bson.M{
			"password": passwordhashed,
		})
		if (err != nil) {
			handleError("error updating user in database: " + err.Error(), http.StatusBadRequest, response)
			return
		}
		Logger.Info("User password reset",
			zap.String("id", id),
			zap.String("email", userData["email"].(string)),
		)
    response.Header().Set("content-type", "application/json")
    response.Write([]byte(`{"message":"password reset successfully"}`))
    break
  }
}

// ValidateLoggedIn validates JWT token to confirm login
func ValidateLoggedIn(t string) (jwt.MapClaims, error) {
	if (t == "") {
		return nil, errors.New("Authorization token must be present")
	}
	token, _ := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		if _, success := token.Method.(*jwt.SigningMethodHMAC); !success {
			return nil, errors.New("There was an error")
		}
		return JwtSecret, nil
	})
	if claims, success := token.Claims.(jwt.MapClaims); success && token.Valid {
		var decodedToken jwt.MapClaims
		mapstructure.Decode(claims, &decodedToken)
		return decodedToken, nil
	}
	return nil, errors.New("Invalid authorization token")
}

func ValidateAdmin(t string) (jwt.MapClaims, error) {
	accountdata, err := ValidateLoggedIn(t)
	if (err != nil) {
		return nil, err
	}
	if (accountdata["emailverified"] != nil && accountdata["emailverified"].(bool)) {
		return nil, errors.New("email not found or verified")
	}
	if (accountdata["type"] != "admin") {
		return nil, errors.New("user not admin")
	}
	return accountdata, nil
}
