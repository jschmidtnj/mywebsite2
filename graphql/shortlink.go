package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/graphql-go/graphql"
	"github.com/rs/xid"
	"go.mongodb.org/mongo-driver/bson"
)

// ShortLinkType account type object for user accounts graphql
var ShortLinkType *graphql.Object = graphql.NewObject(graphql.ObjectConfig{
	Name: "ShortLink",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
		},
		"link": &graphql.Field{
			Type: graphql.String,
		},
	},
})

func createShortLink(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		handleError("short link http method not POST", http.StatusBadRequest, response)
		return
	}
	var shortlinkdata map[string]interface{}
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		handleError("error getting request body: "+err.Error(), http.StatusBadRequest, response)
		return
	}
	err = json.Unmarshal(body, &shortlinkdata)
	if err != nil {
		handleError("error parsing request body: "+err.Error(), http.StatusBadRequest, response)
		return
	}
	if !(shortlinkdata["link"] != nil && shortlinkdata["recaptcha"] != nil) {
		handleError("no link or password or recaptcha token provided", http.StatusBadRequest, response)
		return
	}
	link, ok := shortlinkdata["link"].(string)
	if !ok {
		handleError("link cannot be cast to string", http.StatusBadRequest, response)
		return
	}
	decodedLink, err := url.QueryUnescape(link)
	if err != nil {
		handleError("cannot decode url from query: "+err.Error(), http.StatusBadRequest, response)
		return
	}
	_, err = url.ParseRequestURI(decodedLink)
	if err != nil {
		handleError("invalid url: "+err.Error(), http.StatusBadRequest, response)
		return
	}
	recaptchatoken, ok := shortlinkdata["recaptcha"].(string)
	if !ok {
		handleError("recaptcha token cannot be cast to string", http.StatusBadRequest, response)
		return
	}
	err = verifyRecaptcha(recaptchatoken, shortlinkRecaptchaSecret)
	if err != nil {
		handleError("recaptcha error: "+err.Error(), http.StatusUnauthorized, response)
		return
	}
	linkid, err := generateShortLink(link)
	if err != nil {
		handleError("short link generate error: "+err.Error(), http.StatusUnauthorized, response)
		return
	}
	response.Header().Set("content-type", "application/json")
	response.Write([]byte(`{"id":"` + linkid + `"}`))
}

func generateShortLink(link string) (string, error) {
	guid := xid.New()
	id := guid.String()
	decodedLink, err := url.QueryUnescape(link)
	if err != nil {
		return "", err
	}
	shortLinkData := bson.M{
		"_id":  id,
		"link": decodedLink,
	}
	_, err = shortLinkCollection.InsertOne(ctxMongo, shortLinkData)
	if err != nil {
		return "", err
	}
	return id, nil
}

func deleteShortLink(id string) error {
	_, err := shortLinkCollection.DeleteOne(ctxMongo, bson.M{
		"_id": id,
	})
	if err != nil {
		return err
	}
	return nil
}

func getShortLink(idstring string) (string, error) {
	var idstringarr = []string{
		idstring,
	}
	idstrings, err := getShortLinks(idstringarr)
	if err != nil {
		return "", err
	}
	if len(idstrings) == 0 {
		return "", errors.New("no link found")
	}
	return idstrings[0], nil
}

func getShortLinks(ids []string) ([]string, error) {
	cursor, err := shortLinkCollection.Find(ctxMongo, bson.M{
		"_id": bson.M{
			"$in": ids,
		},
	})
	defer cursor.Close(ctxMongo)
	if err != nil {
		return nil, errors.New("error: " + err.Error())
	}
	fullLinks := make([]string, 0)
	var foundstuff = false
	for cursor.Next(ctxMongo) {
		foundstuff = true
		shortLinkDataPrimitive := &bson.D{}
		err = cursor.Decode(shortLinkDataPrimitive)
		if err != nil {
			return nil, errors.New("problem decoding shortlink data: " + err.Error())
		}
		shortLinkData := shortLinkDataPrimitive.Map()
		fullLink, ok := shortLinkData["link"].(string)
		if !ok {
			return nil, errors.New("cannot cast link to string")
		}
		fullLinks = append(fullLinks, fullLink)
	}
	if !foundstuff {
		return fullLinks, nil
	}
	return fullLinks, nil
}

func shortLinkRedirect(response http.ResponseWriter, request *http.Request) {
	id := request.URL.Query().Get("id")
	if len(id) != 20 {
		http.Redirect(response, request, shortlinkURL+"/404", 301)
		return
	}
	fullLink, err := getShortLink(id)
	if err != nil {
		http.Redirect(response, request, shortlinkURL+"/404", 301)
		return
	}
	logger.Info("got link " + fullLink)
	http.Redirect(response, request, fullLink, 301)
}
