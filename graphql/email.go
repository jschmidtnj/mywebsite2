package main

import (
	"errors"
	"net/http"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/rest"
	"github.com/PuerkitoBio/goquery"
	"time"
	"io/ioutil"
	"encoding/json"
	"strconv"
	"fmt"
)

var sendEmail string = "noreply@joshuaschmidt.tech"

var sendgridApiUrl string = "https://api.sendgrid.com"

var sendgridApiPath string = "/v3"

func SendEmailVerification(email string) (*rest.Response, error) {
	expirationTime := time.Now().Add(time.Duration(TokenExpiration) * time.Hour)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	"email": email,
	"verify": true,
		"StandardClaims": jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			Issuer: "Joshua Schmidt",
		},
	})
	tokenString, err := token.SignedString(JwtSecret)
	if (err != nil) {
		return nil, err
	}
	res, err := http.Get(WebsiteUrl + "/emailtemplates/verifyemail.html")
	if (err != nil) {
		return nil, err
	}
	defer res.Body.Close()
	if (res.StatusCode != 200) {
		return nil, errors.New("could not get email template")
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if (err != nil) {
		return nil, err
	}
	doc.Find("#verify").SetAttr("href", WebsiteUrl + "/signin?token=" + tokenString)
	template, err := doc.Html()
	request := sendgrid.GetRequest(SendgridApiKey, sendgridApiPath + "/mail/send", sendgridApiUrl)
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
	body = fmt.Sprintf(body, email, sendEmail)
	if (err != nil) {
		return nil, err
	}
	var bodyJson map[string]interface{}
	err = json.Unmarshal([]byte(body), &bodyJson)
	if (err != nil) {
		return nil, err
	}
	bodyJson["content"].([]interface{})[0].(map[string]interface{})["value"] = template
	bodyBytes, err := json.Marshal(bodyJson)
	if (err != nil) {
		return nil, err
	}
	request.Body = bodyBytes
	response, err := sendgrid.API(request)
	return response, err
}

func SendTestEmail(response http.ResponseWriter, request *http.Request) {
	if (request.Method != http.MethodPost) {
		handleError("register http method not POST", http.StatusBadRequest, response)
		return
	}
	var emaildata map[string]interface{}
	emaildatabody, err := ioutil.ReadAll(request.Body)
	if (err != nil) {
		handleError("error getting request body: " + err.Error(), http.StatusBadRequest, response)
		return
	}
  err = json.Unmarshal(emaildatabody, &emaildata)
	if (err != nil) {
		handleError("error parsing request body: " + err.Error(), http.StatusBadRequest, response)
		return
	}
	if (!(emaildata["to"] != nil && emaildata["from"] != nil && emaildata["content"] != nil && emaildata["subject"] != nil)) {
		handleError("no to or from or content or subject provided", http.StatusBadRequest, response)
		return
	}
	res, err := http.Get(WebsiteUrl + "/emailtemplates/test.html")
	if (err != nil) {
		handleError(err.Error(), http.StatusBadRequest, response)
		return
	}
	defer res.Body.Close()
	if (res.StatusCode != 200) {
		handleError("could not get email template", http.StatusBadRequest, response)
		return
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if (err != nil) {
		handleError(err.Error(), http.StatusBadRequest, response)
		return
	}
	doc.Find("#content").SetHtml(emaildata["content"].(string))
	template, err := doc.Html()
	req := sendgrid.GetRequest(SendgridApiKey, sendgridApiPath + "/mail/send", sendgridApiUrl)
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
	var bodyJson map[string]interface{}
	err = json.Unmarshal([]byte(body), &bodyJson)
	if (err != nil) {
		handleError(err.Error(), http.StatusBadRequest, response)
		return
	}
	bodyJson["content"].([]interface{})[0].(map[string]interface{})["value"] = template
	bodyBytes, err := json.Marshal(bodyJson)
	if (err != nil) {
		handleError(err.Error(), http.StatusBadRequest, response)
		return
	}
	req.Body = bodyBytes
	res1, err := sendgrid.API(req)
	if (err != nil) {
		handleError(err.Error(), http.StatusBadRequest, response)
		return
	}
	if (res1.StatusCode != 202) {
		handleError("invalid response code from email: " + strconv.Itoa(res1.StatusCode) + ", body: " + res1.Body, http.StatusBadRequest, response)
		return
	}
	response.Header().Set("content-type", "application/json")
	response.Write([]byte(`{"message":"email successfully sent"}`))
}
