package main

import (
	"errors"
	"net/http"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/rest"
	"github.com/PuerkitoBio/goquery"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

var sendgridApiUrl string = "https://api.sendgrid.com"

var sendgridApiPath string = "/v3/api_keys"

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
	request := sendgrid.GetRequest(SendgridApiKey, sendgridApiPath, sendgridApiUrl)
	request.Method = "POST"
	body := bson.M{
		"personalizations": bson.M{
			"to": bson.A{
				bson.M{
					"email": email,
				},
			},
			"subject": "Please Verify Email",
		},
		"from": bson.M{
			"email": "noreply@joshuaschmidt.tech",
		},
		"content": bson.A{
			bson.M{
				"type": "text/html",
				"value": template,
			},
		},
	}
	bodybytes, err := bson.Marshal(body)
	if (err != nil) {
		return nil, err
	}
	request.Body = bodybytes
	response, err := sendgrid.API(request)
	return response, err
}

func SendTestEmail(email string) (*rest.Response, error) {
	res, err := http.Get(WebsiteUrl + "/emailtemplates/test.html")
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
	doc.Find("#content").SetHtml("test content")
	template, err := doc.Html()
	request := sendgrid.GetRequest(SendgridApiKey, sendgridApiPath, sendgridApiUrl)
	request.Method = "POST"
	body := bson.M{
		"personalizations": bson.M{
			"to": bson.A{
				bson.M{
					"email": email,
				},
			},
			"subject": "Please Verify Email",
		},
		"from": bson.M{
			"email": "noreply@joshuaschmidt.tech",
		},
		"content": bson.A{
			bson.M{
				"type": "text/html",
				"value": template,
			},
		},
	}
	bodybytes, err := bson.Marshal(body)
	if (err != nil) {
		return nil, err
	}
	request.Body = bodybytes
	response, err := sendgrid.API(request)
	return response, err
}
