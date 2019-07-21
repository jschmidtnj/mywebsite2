package main

import (
	"encoding/json"
	"errors"
	"github.com/graphql-go/graphql"
	"github.com/rs/xid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io/ioutil"
	"net/http"
	"net/url"
)

// AccountType account type object for user accounts graphql
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
	if !manageCors(&response, request) {
		return
	}
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
	_, err = url.ParseRequestURI(link)
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
	shortLinkData := bson.M{
		"_id":  id,
		"link": link,
	}
	_, err := shortLinkCollection.InsertOne(ctxMongo, shortLinkData)
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
	return idstrings[0], nil
}

func getShortLinks(idstrings []string) ([]string, error) {
	ids := make([]primitive.ObjectID, len(idstrings))
	for i, idstring := range idstrings {
		id, err := primitive.ObjectIDFromHex(idstring)
		if err != nil {
			return nil, err
		}
		ids[i] = id
	}
	cursor, err := shortLinkCollection.Find(ctxMongo, bson.M{
		"_id": bson.M{
			"$in": ids,
		},
	})
	defer cursor.Close(ctxMongo)
	if err != nil {
		return nil, errors.New("no link data found")
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
		return nil, errors.New("no link data found")
	}
	return fullLinks, nil
}

func shortLinkRedirect(response http.ResponseWriter, request *http.Request) {
	id := request.URL.Query().Get("id")
	if len(id) != 20 {
		handleError("error getting valid id from query", http.StatusBadRequest, response)
		return
	}
	fullLink, err := getShortLink(id)
	if err != nil {
		handleError(err.Error(), http.StatusBadRequest, response)
		return
	}
	http.Redirect(response, request, fullLink, 301)
}
