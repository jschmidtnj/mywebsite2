package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go"
)

var sendgridAPIUrl = "https://api.sendgrid.com"

var sendgridAPIPath = "/v3"

func sendEmailVerification(email string) (*rest.Response, error) {
	expirationTime := time.Now().Add(time.Duration(tokenExpiration) * time.Hour)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":  email,
		"verify": true,
		"StandardClaims": jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			Issuer:    jwtIssuer,
		},
	})
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return nil, err
	}
	res, err := http.Get(websiteURL + "/emailtemplates/verifyemail.html")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, errors.New("could not get email template")
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}
	doc.Find("#verify").SetAttr("href", websiteURL+"/login?verify=true&token="+tokenString)
	template, err := doc.Html()
	request := sendgrid.GetRequest(sendgridAPIKey, sendgridAPIPath+"/mail/send", sendgridAPIUrl)
	request.Method = "POST"
	body := `
	{
		"personalizations": [
			{
				"to": [
					{
						"email": "%s"
					}
				],
				"subject": "Please Verify Email"
			}
		],
		"from": {
			"email": "%s"
		},
		"content": [
			{
				"type": "text/html",
				"value": ""
			}
		]
	}
	`
	body = fmt.Sprintf(body, email, serviceEmail)
	if err != nil {
		return nil, err
	}
	var bodyJSON map[string]interface{}
	err = json.Unmarshal([]byte(body), &bodyJSON)
	if err != nil {
		return nil, err
	}
	bodyJSON["content"].([]interface{})[0].(map[string]interface{})["value"] = template
	bodyBytes, err := json.Marshal(bodyJSON)
	if err != nil {
		return nil, err
	}
	request.Body = bodyBytes
	response, err := sendgrid.API(request)
	return response, err
}

/**
 * @api {put} /sendResetEmail Send reset email
 * @apiVersion 0.0.1
 * @apiParam {String} email User email
 * @apiSuccess {String} message Success message for email sent
 * @apiGroup emails
 */
func sendPasswordResetEmail(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPut {
		handleError("reset http method not PUT", http.StatusBadRequest, response)
		return
	}
	var resetdata map[string]interface{}
	emaildatabody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		handleError("error getting request body: "+err.Error(), http.StatusBadRequest, response)
		return
	}
	err = json.Unmarshal(emaildatabody, &resetdata)
	if err != nil {
		handleError("error parsing request body: "+err.Error(), http.StatusBadRequest, response)
		return
	}
	if resetdata["email"] == nil {
		handleError("no email provided", http.StatusBadRequest, response)
		return
	}
	email, ok := resetdata["email"].(string)
	if !ok {
		handleError("email cannot be cast to string", http.StatusBadRequest, response)
		return
	}
	recaptchatoken, ok := resetdata["recaptcha"].(string)
	if !ok {
		handleError("recaptcha token cannot be cast to string", http.StatusBadRequest, response)
		return
	}
	err = verifyRecaptcha(recaptchatoken, mainRecaptchaSecret)
	if err != nil {
		handleError("recaptcha error: "+err.Error(), http.StatusUnauthorized, response)
		return
	}
	expirationTime := time.Now().Add(time.Duration(tokenExpiration) * time.Hour)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"reset": true,
		"StandardClaims": jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			Issuer:    jwtIssuer,
		},
	})
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		handleError(err.Error(), http.StatusBadRequest, response)
		return
	}
	res, err := http.Get(websiteURL + "/emailtemplates/passwordreset.html")
	if err != nil {
		handleError(err.Error(), http.StatusBadRequest, response)
		return
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		handleError("could not get email template", http.StatusBadRequest, response)
		return
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		handleError(err.Error(), http.StatusBadRequest, response)
		return
	}
	doc.Find("#reset").SetAttr("href", websiteURL+"/reset?reset=true&token="+tokenString)
	template, err := doc.Html()
	req := sendgrid.GetRequest(sendgridAPIKey, sendgridAPIPath+"/mail/send", sendgridAPIUrl)
	req.Method = "POST"
	body := `
	{
		"personalizations": [
			{
				"to": [
					{
						"email": "%s"
					}
				],
				"subject": "Reset Password"
			}
		],
		"from": {
			"email": "%s"
		},
		"content": [
			{
				"type": "text/html",
				"value": ""
			}
		]
	}
	`
	body = fmt.Sprintf(body, email, serviceEmail)
	if err != nil {
		handleError(err.Error(), http.StatusBadRequest, response)
		return
	}
	var bodyJSON map[string]interface{}
	err = json.Unmarshal([]byte(body), &bodyJSON)
	if err != nil {
		handleError(err.Error(), http.StatusBadRequest, response)
		return
	}
	bodyJSON["content"].([]interface{})[0].(map[string]interface{})["value"] = template
	bodyBytes, err := json.Marshal(bodyJSON)
	if err != nil {
		handleError(err.Error(), http.StatusBadRequest, response)
		return
	}
	req.Body = bodyBytes
	res1, err := sendgrid.API(req)
	if err != nil {
		handleError(err.Error(), http.StatusBadRequest, response)
		return
	}
	if res1.StatusCode != 202 {
		handleError("invalid response code from email: "+strconv.Itoa(res1.StatusCode)+", body: "+res1.Body, http.StatusBadRequest, response)
		return
	}
	response.Header().Set("content-type", "application/json")
	response.Write([]byte(`{"message":"reset email sent"}`))
}

/**
 * @api {put} /sendTestEmail Send test email
 * @apiVersion 0.0.1
 * @apiSuccess {String} message Success message for email sent
 * @apiGroup emails
 */
func sendTestEmail(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		handleError("register http method not POST", http.StatusBadRequest, response)
		return
	}
	if _, err := validateAdmin(getAuthToken(request)); err != nil {
		handleError("auth error: "+err.Error(), http.StatusBadRequest, response)
		return
	}
	var emaildata map[string]interface{}
	emaildatabody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		handleError("error getting request body: "+err.Error(), http.StatusBadRequest, response)
		return
	}
	err = json.Unmarshal(emaildatabody, &emaildata)
	if err != nil {
		handleError("error parsing request body: "+err.Error(), http.StatusBadRequest, response)
		return
	}
	if !(emaildata["to"] != nil && emaildata["from"] != nil && emaildata["content"] != nil && emaildata["subject"] != nil) {
		handleError("no to or from or content or subject provided", http.StatusBadRequest, response)
		return
	}
	res, err := http.Get(websiteURL + "/emailtemplates/test.html")
	if err != nil {
		handleError(err.Error(), http.StatusBadRequest, response)
		return
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		handleError("could not get email template", http.StatusBadRequest, response)
		return
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		handleError(err.Error(), http.StatusBadRequest, response)
		return
	}
	doc.Find("#content").SetHtml(emaildata["content"].(string))
	template, err := doc.Html()
	req := sendgrid.GetRequest(sendgridAPIKey, sendgridAPIPath+"/mail/send", sendgridAPIUrl)
	req.Method = "POST"
	body := `
	{
		"personalizations": [
			{
				"to": [
					{
						"email": "%s"
					}
				],
				"subject": "%s"
			}
		],
		"from": {
			"email": "%s"
		},
		"content": [
			{
				"type": "text/html",
				"value": ""
			}
		]
	}
	`
	body = fmt.Sprintf(body, emaildata["to"], emaildata["subject"], emaildata["from"])
	var bodyJSON map[string]interface{}
	err = json.Unmarshal([]byte(body), &bodyJSON)
	if err != nil {
		handleError(err.Error(), http.StatusBadRequest, response)
		return
	}
	bodyJSON["content"].([]interface{})[0].(map[string]interface{})["value"] = template
	bodyBytes, err := json.Marshal(bodyJSON)
	if err != nil {
		handleError(err.Error(), http.StatusBadRequest, response)
		return
	}
	req.Body = bodyBytes
	res1, err := sendgrid.API(req)
	if err != nil {
		handleError(err.Error(), http.StatusBadRequest, response)
		return
	}
	if res1.StatusCode != 202 {
		handleError("invalid response code from email: "+strconv.Itoa(res1.StatusCode)+", body: "+res1.Body, http.StatusBadRequest, response)
		return
	}
	response.Header().Set("content-type", "application/json")
	response.Write([]byte(`{"message":"email successfully sent"}`))
}
