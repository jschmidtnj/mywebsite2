package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/mitchellh/mapstructure"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

/**
 * authentication for api
 */

/**
 * @apiDefine authToken
 * @apiHeader {String} Authorization "Bearer `token`" - login auth token
 */

var numHashes = 12

type loginClaims struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Type  string `json:"type"`
	jwt.StandardClaims
}

func verifyRecaptcha(recaptchaToken string, recaptchaSecret string) error {
	request, err := http.NewRequest("POST", "https://www.google.com/recaptcha/api/siteverify", nil)
	if err != nil {
		return err
	}
	query := request.URL.Query()
	query.Add("secret", recaptchaSecret)
	query.Add("response", recaptchaToken)
	request.URL.RawQuery = query.Encode()
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	var responsedata map[string]interface{}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &responsedata)
	if err != nil {
		return err
	}
	if !responsedata["success"].(bool) {
		errorcodes := responsedata["error-codes"].([]interface{})
		codesStrArr := make([]string, len(errorcodes))
		for i, code := range errorcodes {
			codeStr, ok := code.(string)
			if !ok {
				return errors.New("cannot cast error code to string in recaptcha result")
			}
			codesStrArr[i] = codeStr
		}
		return errors.New("invalid recaptcha token: " + strings.Join(codesStrArr, ", "))
	}
	return nil
}

/**
 * @api {post} /register User registration
 * @apiVersion 0.0.1
 * @apiParam {String} email Registration email
 * @apiParam {String} password User password
 * @apiSuccess {String} message Response message
 * @apiGroup authentication
 */
func register(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		handleError("register http method not POST", http.StatusBadRequest, response)
		return
	}
	var registerdata map[string]interface{}
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		handleError("error getting request body: "+err.Error(), http.StatusBadRequest, response)
		return
	}
	err = json.Unmarshal(body, &registerdata)
	if err != nil {
		handleError("error parsing request body: "+err.Error(), http.StatusBadRequest, response)
		return
	}
	if !(registerdata["password"] != nil && registerdata["email"] != nil && registerdata["recaptcha"] != nil) {
		handleError("no email or password or recaptcha token provided", http.StatusBadRequest, response)
		return
	}
	password, ok := registerdata["password"].(string)
	if !ok {
		handleError("password cannot be cast to string", http.StatusBadRequest, response)
		return
	}
	email, ok := registerdata["email"].(string)
	if !ok {
		handleError("email cannot be cast to string", http.StatusBadRequest, response)
		return
	}
	recaptchatoken, ok := registerdata["recaptcha"].(string)
	if !ok {
		handleError("recaptcha token cannot be cast to string", http.StatusBadRequest, response)
		return
	}
	err = verifyRecaptcha(recaptchatoken, mainRecaptchaSecret)
	if err != nil {
		handleError("recaptcha error: "+err.Error(), http.StatusUnauthorized, response)
		return
	}
	countemail, err := userCollection.CountDocuments(ctxMongo, bson.M{"email": email})
	if err != nil {
		handleError("error counting users with same email: "+err.Error(), http.StatusBadRequest, response)
		return
	}
	if countemail > 0 {
		handleError("email is already taken", http.StatusBadRequest, response)
		return
	}
	passwordhashed, err := bcrypt.GenerateFromPassword([]byte(password), numHashes)
	if err != nil {
		handleError("error hashing password: "+err.Error(), http.StatusBadRequest, response)
		return
	}
	emailres, err := sendEmailVerification(email)
	if err != nil {
		handleError("error sending email verification: "+err.Error(), http.StatusBadRequest, response)
		return
	} else if emailres.StatusCode != 202 {
		handleError("error sending email verification: got status code "+strconv.Itoa(emailres.StatusCode), http.StatusBadRequest, response)
		return
	}
	res, err := userCollection.InsertOne(ctxMongo, bson.M{
		"email":         email,
		"password":      string(passwordhashed),
		"emailverified": false,
		"type":          "user",
		"shortlinks":    []string{},
	})
	if err != nil {
		handleError("error inserting user to database: "+err.Error(), http.StatusBadRequest, response)
		return
	}
	id := res.InsertedID.(primitive.ObjectID).Hex()
	logger.Info("User register",
		zap.String("id", id),
		zap.String("email", registerdata["email"].(string)),
	)
	response.Header().Set("Content-Type", "application/json")
	response.Write([]byte(`{"message":"please check email for verification"}`))
}

/**
 * @api {put} /loginEmailPassword User login
 * @apiVersion 0.0.1
 * @apiParam {String} email User email
 * @apiParam {String} password User password
 * @apiSuccess {String} token User token for authenticated requests
 * @apiGroup authentication
 */
func loginEmailPassword(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPut {
		handleError("login http method not PUT", http.StatusBadRequest, response)
		return
	}
	var logindata map[string]interface{}
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		handleError("error getting request body: "+err.Error(), http.StatusBadRequest, response)
		return
	}
	err = json.Unmarshal(body, &logindata)
	if err != nil {
		handleError("error parsing request body: "+err.Error(), http.StatusBadRequest, response)
		return
	}
	if logindata["email"] == nil || logindata["password"] == nil || logindata["recaptcha"] == nil {
		handleError("no email or password or recaptcha token provided", http.StatusBadRequest, response)
		return
	}
	email, ok := logindata["email"].(string)
	if !ok {
		handleError("email cannot be cast to string", http.StatusBadRequest, response)
		return
	}
	password, ok := logindata["password"].(string)
	if !ok {
		handleError("password cannot be cast to string", http.StatusBadRequest, response)
		return
	}
	recaptchatoken, ok := logindata["recaptcha"].(string)
	if !ok {
		handleError("recaptcha token cannot be cast to string", http.StatusBadRequest, response)
		return
	}
	err = verifyRecaptcha(recaptchatoken, mainRecaptchaSecret)
	if err != nil {
		handleError("recaptcha error: "+err.Error(), http.StatusUnauthorized, response)
		return
	}
	cursor, err := userCollection.Find(ctxMongo, bson.M{"email": email})
	defer cursor.Close(ctxMongo)
	if err != nil {
		handleError("error finding user: "+err.Error(), http.StatusUnauthorized, response)
		return
	}
	var foundstuff = false
	for cursor.Next(ctxMongo) {
		userDataPrimitive := &bson.D{}
		err = cursor.Decode(userDataPrimitive)
		if err != nil {
			handleError("error decoding user data: "+err.Error(), http.StatusBadRequest, response)
			return
		}
		userData := userDataPrimitive.Map()
		if !userData["emailverified"].(bool) {
			handleError("email not verified", http.StatusUnauthorized, response)
			return
		}
		err = bcrypt.CompareHashAndPassword([]byte(userData["password"].(string)), []byte(password))
		if err != nil {
			handleError("invalid password: "+err.Error(), http.StatusUnauthorized, response)
			return
		}
		id := userData["_id"].(primitive.ObjectID).Hex()
		expirationTime := time.Now().Add(time.Duration(tokenExpiration) * time.Hour)
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, loginClaims{
			id,
			userData["email"].(string),
			userData["type"].(string),
			jwt.StandardClaims{
				ExpiresAt: expirationTime.Unix(),
				Issuer:    jwtIssuer,
			},
		})
		tokenString, err := token.SignedString(jwtSecret)
		if err != nil {
			handleError("error creating token: "+err.Error(), http.StatusBadRequest, response)
			return
		}
		logger.Info("User login",
			zap.String("id", id),
		)
		response.Header().Set("Content-Type", "application/json")
		response.Write([]byte(`{ "token": "` + tokenString + `" }`))
		foundstuff = true
		break
	}
	if !foundstuff {
		handleError("no user data found", http.StatusBadRequest, response)
	}
}

/**
 * @api {put} /logoutEmailPassword User logout
 * @apiVersion 0.0.1
 * @apiUse authToken
 * @apiSuccess {String} message Response message
 * @apiGroup authentication
 */
func logoutEmailPassword(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPut {
		handleError("logout http method not PUT", http.StatusBadRequest, response)
		return
	}
	_, err := validateLoggedIn(getAuthToken(request))
	if err != nil {
		handleError(err.Error(), http.StatusBadRequest, response)
		return
	}
	response.Header().Set("Content-Type", "application/json")
	response.Write([]byte(`{ "message": "successfully signed out" }`))
}

func verifyEmail(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		handleError("verify http method not POST", http.StatusBadRequest, response)
		return
	}
	var verifydata map[string]interface{}
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		handleError("error getting request body: "+err.Error(), http.StatusBadRequest, response)
		return
	}
	err = json.Unmarshal(body, &verifydata)
	if err != nil {
		handleError("error parsing request body: "+err.Error(), http.StatusBadRequest, response)
		return
	}
	if verifydata["token"] == nil {
		handleError("no token provided", http.StatusBadRequest, response)
		return
	}
	giventoken, ok := verifydata["token"].(string)
	if !ok {
		handleError("token cannot be cast to string", http.StatusBadRequest, response)
		return
	}
	token, err := jwt.Parse(giventoken, func(token *jwt.Token) (interface{}, error) {
		if _, success := token.Method.(*jwt.SigningMethodHMAC); !success {
			return nil, errors.New("There was an error")
		}
		return jwtSecret, nil
	})
	if err != nil {
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
	if !(decodedToken["email"] != nil && decodedToken["verify"] != nil) {
		handleError("token does not contian email or verify", http.StatusBadRequest, response)
		return
	}
	email, ok := decodedToken["email"].(string)
	if !ok {
		handleError("email in token cannot be cast to string", http.StatusBadRequest, response)
		return
	}
	verify, ok := decodedToken["verify"].(bool)
	if !ok {
		handleError("verify in token cannot be cast to boolean", http.StatusBadRequest, response)
		return
	}
	if !verify {
		handleError("verify in token is false", http.StatusBadRequest, response)
		return
	}
	cursor, err := userCollection.Find(ctxMongo, bson.M{"email": email})
	if err != nil {
		handleError("error finding user: "+err.Error(), http.StatusBadRequest, response)
		return
	}
	defer cursor.Close(ctxMongo)
	var foundstuff = false
	for cursor.Next(ctxMongo) {
		userDataPrimitive := &bson.D{}
		err = cursor.Decode(userDataPrimitive)
		if err != nil {
			handleError("error decoding user data: "+err.Error(), http.StatusBadRequest, response)
			return
		}
		userData := userDataPrimitive.Map()
		if userData["emailverified"] != nil && userData["emailverified"].(bool) {
			handleError("email already verified", http.StatusBadRequest, response)
			return
		}
		id := userData["_id"].(primitive.ObjectID)
		_, err := userCollection.UpdateOne(ctxMongo, bson.M{
			"_id": id,
		}, bson.M{
			"$set": bson.M{
				"emailverified": true,
			},
		})
		if err != nil {
			handleError("error updating user in database: "+err.Error(), http.StatusBadRequest, response)
			return
		}
		idstring := id.Hex()
		logger.Info("User email verify",
			zap.String("id", idstring),
			zap.String("email", userData["email"].(string)),
		)
		response.Header().Set("Content-Type", "application/json")
		response.Write([]byte(`{"message":"email successfully verified"}`))
		foundstuff = true
		break
	}
	if !foundstuff {
		handleError("no user data found", http.StatusBadRequest, response)
	}
}

/**
 * @api {put} /resetPassword Reset password
 * @apiVersion 0.0.1
 * @apiParam {String} token Password reset token
 * @apiParam {String} password New password
 * @apiSuccess {String} message Response message
 * @apiGroup authentication
 */
func resetPassword(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		handleError("reset http method not POST", http.StatusBadRequest, response)
		return
	}
	var resetdata map[string]interface{}
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		handleError("error getting request body: "+err.Error(), http.StatusBadRequest, response)
		return
	}
	err = json.Unmarshal(body, &resetdata)
	if err != nil {
		handleError("error parsing request body: "+err.Error(), http.StatusBadRequest, response)
		return
	}
	if resetdata["token"] == nil || resetdata["password"] == nil {
		handleError("no token or new password provided", http.StatusBadRequest, response)
		return
	}
	giventoken, ok := resetdata["token"].(string)
	if !ok {
		handleError("token cannot be cast to string", http.StatusBadRequest, response)
		return
	}
	password, ok := resetdata["password"].(string)
	if !ok {
		handleError("password cannot be cast to string", http.StatusBadRequest, response)
		return
	}
	token, err := jwt.Parse(giventoken, func(token *jwt.Token) (interface{}, error) {
		if _, success := token.Method.(*jwt.SigningMethodHMAC); !success {
			return nil, errors.New("There was an error")
		}
		return jwtSecret, nil
	})
	if err != nil {
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
	if !(decodedToken["email"] != nil && decodedToken["reset"] != nil) {
		handleError("token does not contian email or reset", http.StatusBadRequest, response)
		return
	}
	email, ok := decodedToken["email"].(string)
	if !ok {
		handleError("email in token cannot be cast to string", http.StatusBadRequest, response)
		return
	}
	reset, ok := decodedToken["reset"].(bool)
	if !ok {
		handleError("reset in token cannot be cast to boolean", http.StatusBadRequest, response)
		return
	}
	if !reset {
		handleError("reset in token is false", http.StatusBadRequest, response)
		return
	}
	cursor, err := userCollection.Find(ctxMongo, bson.M{"email": email})
	defer cursor.Close(ctxMongo)
	if err != nil {
		handleError("error finding user: "+err.Error(), http.StatusBadRequest, response)
		return
	}
	var foundstuff = false
	for cursor.Next(ctxMongo) {
		userDataPrimitive := &bson.D{}
		err = cursor.Decode(userDataPrimitive)
		if err != nil {
			handleError("error decoding user data: "+err.Error(), http.StatusBadRequest, response)
			return
		}
		userData := userDataPrimitive.Map()
		if !(userData["_id"] != nil && userData["emailverified"] != nil && userData["emailverified"].(bool)) {
			handleError("user id invalid or email not verified", http.StatusBadRequest, response)
			return
		}
		id := userData["_id"].(primitive.ObjectID)
		passwordhashed, err := bcrypt.GenerateFromPassword([]byte(password), numHashes)
		if err != nil {
			handleError("error hashing password: "+err.Error(), http.StatusBadRequest, response)
			return
		}
		_, err = userCollection.UpdateOne(ctxMongo, bson.M{
			"_id": id,
		}, bson.M{
			"$set": bson.M{
				"password": string(passwordhashed),
			},
		})
		if err != nil {
			handleError("error updating user in database: "+err.Error(), http.StatusBadRequest, response)
			return
		}
		idStr := id.Hex()
		logger.Info("User password reset",
			zap.String("id", idStr),
			zap.String("email", userData["email"].(string)),
		)
		response.Header().Set("Content-Type", "application/json")
		response.Write([]byte(`{"message":"password reset successfully"}`))
		foundstuff = true
		break
	}
	if !foundstuff {
		handleError("no user data found", http.StatusBadRequest, response)
	}
}

// validateLoggedIn validates JWT token to confirm login
func validateLoggedIn(t string) (jwt.MapClaims, error) {
	if t == "" {
		return nil, errors.New("Authorization token must be present")
	}
	token, err := jwt.ParseWithClaims(t, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("unable to parse token claims")
	}
	return claims, nil
}

func validateAdmin(t string) (jwt.MapClaims, error) {
	accountdata, err := validateLoggedIn(t)
	if err != nil {
		return nil, err
	}
	if accountdata["emailverified"] != nil && accountdata["emailverified"].(bool) {
		return nil, errors.New("email not found or verified")
	}
	if accountdata["type"] != "admin" {
		return nil, errors.New("user not admin")
	}
	return accountdata, nil
}
