package main

import (
	"bytes"
	"cloud.google.com/go/storage"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func createPostPicture(response http.ResponseWriter, request *http.Request) {
	if !manageCors(response, request) {
		return
	}
	if request.Method != http.MethodPut {
		handleError("create post picture http method not PUT", http.StatusBadRequest, response)
		return
	}
	if _, err := validateAdmin(getAuthToken(request)); err != nil {
		handleError("auth error: "+err.Error(), http.StatusBadRequest, response)
		return
	}
	thetype := request.URL.Query().Get("type")
	if thetype == "" {
		handleError("error getting type from query", http.StatusBadRequest, response)
		return
	}
	if !validType(thetype) {
		handleError("invalid type in query", http.StatusBadRequest, response)
		return
	}
	var mongoCollection *mongo.Collection
	if thetype == "blog" {
		mongoCollection = blogCollection
	} else {
		mongoCollection = projectCollection
	}
	postid := request.URL.Query().Get("postid")
	if postid == "" {
		handleError("error getting post id from query", http.StatusBadRequest, response)
		return
	}
	id, err := primitive.ObjectIDFromHex(postid)
	if err != nil {
		handleError("error creating objectid from postid: "+err.Error(), http.StatusBadRequest, response)
		return
	}
	imageid := request.URL.Query().Get("imageid")
	if imageid == "" {
		handleError("error getting image id from query", http.StatusBadRequest, response)
		return
	}
	var filebuffer bytes.Buffer
	file, header, err := request.FormFile("file")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	name := strings.Split(header.Filename, ".")
	logger.Info("File name: " + name[0])
	io.Copy(&filebuffer, file)
	var fileobj *storage.ObjectHandle
	if thetype == "blog" {
		fileobj = imageBucket.Object(blogImageIndex + "/" + postid + "/" + imageid)
	} else {
		fileobj = imageBucket.Object(projectImageIndex + "/" + postid + "/" + imageid)
	}
	filewriter := fileobj.NewWriter(ctxStorage)
	if byteswritten, err := filebuffer.WriteTo(filewriter); err != nil {
		handleError("error writing to filewriter: num bytes: "+strconv.FormatInt(byteswritten, 10)+", "+err.Error(), http.StatusBadRequest, response)
		return
	}
	if err := filewriter.Close(); err != nil {
		handleError("error closing writer: "+err.Error(), http.StatusBadRequest, response)
		return
	}
	contents := filebuffer.String()
	logger.Info(contents)
	filebuffer.Reset()
	_, err = mongoCollection.UpdateOne(ctxMongo, bson.M{
		"_id": id,
	}, bson.M{
		"$push": bson.M{
			"images": imageid,
		},
	})
	if err != nil {
		handleError("error updating mongodb: "+err.Error(), http.StatusBadRequest, response)
		return
	}
	response.Header().Set("content-type", "application/json")
	response.Write([]byte(`{"message":"file uploaded"}`))
}

func updatePostPicture(response http.ResponseWriter, request *http.Request) {
	if !manageCors(response, request) {
		return
	}
	if request.Method != http.MethodPut {
		handleError("update post picture http method not PUT", http.StatusBadRequest, response)
		return
	}
	if _, err := validateAdmin(getAuthToken(request)); err != nil {
		handleError("auth error: "+err.Error(), http.StatusBadRequest, response)
		return
	}
	thetype := request.URL.Query().Get("type")
	if thetype == "" {
		handleError("error getting type from query", http.StatusBadRequest, response)
		return
	}
	if !validType(thetype) {
		handleError("invalid type in query", http.StatusBadRequest, response)
		return
	}
	postid := request.URL.Query().Get("postid")
	if postid == "" {
		handleError("error getting post id from query", http.StatusBadRequest, response)
		return
	}
	_, err := primitive.ObjectIDFromHex(postid)
	if err != nil {
		handleError("error creating objectid from postid: "+err.Error(), http.StatusBadRequest, response)
		return
	}
	imageid := request.URL.Query().Get("imageid")
	if imageid == "" {
		handleError("error getting image id from query", http.StatusBadRequest, response)
		return
	}
	var filebuffer bytes.Buffer
	file, header, err := request.FormFile("file")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	name := strings.Split(header.Filename, ".")
	logger.Info("File name: " + name[0])
	io.Copy(&filebuffer, file)
	var fileobj *storage.ObjectHandle
	if thetype == "blog" {
		fileobj = imageBucket.Object(blogImageIndex + "/" + postid + "/" + imageid)
	} else {
		fileobj = imageBucket.Object(projectImageIndex + "/" + postid + "/" + imageid)
	}
	filewriter := fileobj.NewWriter(ctxStorage)
	if byteswritten, err := filebuffer.WriteTo(filewriter); err != nil {
		handleError("error writing to filewriter: num bytes: "+strconv.FormatInt(byteswritten, 10)+", "+err.Error(), http.StatusBadRequest, response)
		return
	}
	if err := filewriter.Close(); err != nil {
		handleError("error closing writer: "+err.Error(), http.StatusBadRequest, response)
		return
	}
	contents := filebuffer.String()
	logger.Info(contents)
	filebuffer.Reset()
	response.Header().Set("content-type", "application/json")
	response.Write([]byte(`{"message":"file updated"}`))
}

func deletePostPictures(response http.ResponseWriter, request *http.Request) {
	if !manageCors(response, request) {
		return
	}
	if request.Method != http.MethodDelete {
		handleError("delete post picture http method not Delete", http.StatusBadRequest, response)
		return
	}
	if _, err := validateAdmin(getAuthToken(request)); err != nil {
		handleError("auth error: "+err.Error(), http.StatusBadRequest, response)
		return
	}
	var picturedata map[string]interface{}
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		handleError("error getting request body: "+err.Error(), http.StatusBadRequest, response)
		return
	}
	err = json.Unmarshal(body, &picturedata)
	if err != nil {
		handleError("error parsing request body: "+err.Error(), http.StatusBadRequest, response)
		return
	}
	if !(picturedata["imageids"] != nil && picturedata["postid"] != nil && picturedata["type"] != nil) {
		handleError("no imageids or postid or type provided", http.StatusBadRequest, response)
		return
	}
	thetype, ok := picturedata["type"].(string)
	if !ok {
		handleError("unable to cast type to string", http.StatusBadRequest, response)
		return
	}
	if !validType(thetype) {
		handleError("invalid type in body", http.StatusBadRequest, response)
		return
	}
	var mongoCollection *mongo.Collection
	if thetype == "blog" {
		mongoCollection = blogCollection
	} else {
		mongoCollection = projectCollection
	}
	imageids, ok := picturedata["imageids"].([]string)
	if !ok {
		handleError("imageids cannot be cast to string array", http.StatusBadRequest, response)
		return
	}
	postid, ok := picturedata["postid"].(string)
	if !ok {
		handleError("postid cannot be cast to string", http.StatusBadRequest, response)
		return
	}
	id, err := primitive.ObjectIDFromHex(postid)
	if err != nil {
		handleError("error creating objectid from postid: "+err.Error(), http.StatusBadRequest, response)
		return
	}
	for _, imageid := range imageids {
		logger.Info("imageid: " + imageid + ", postid: " + postid)
		var fileobj *storage.ObjectHandle
		if thetype == "blog" {
			fileobj = imageBucket.Object(blogImageIndex + "/" + postid + "/" + imageid)
		} else {
			fileobj = imageBucket.Object(projectImageIndex + "/" + postid + "/" + imageid)
		}
		if err := fileobj.Delete(ctxStorage); err != nil {
			handleError("error deleting file: "+err.Error(), http.StatusBadRequest, response)
			return
		}
	}
	_, err = mongoCollection.UpdateOne(ctxMongo, bson.M{
		"_id": id,
	}, bson.M{
		"$push": bson.M{
			"images": bson.M{
				"$each": imageids,
			},
		},
	})
	if err != nil {
		handleError("error updating mongodb: "+err.Error(), http.StatusBadRequest, response)
		return
	}
	response.Header().Set("content-type", "application/json")
	response.Write([]byte(`{"message":"files deleted"}`))
}

func getPostPicture(response http.ResponseWriter, request *http.Request) {
	if !manageCors(response, request) {
		return
	}
	if request.Method != http.MethodGet {
		handleError("get post picture http method not GET", http.StatusBadRequest, response)
		return
	}
	thetype := request.URL.Query().Get("type")
	if thetype == "" {
		handleError("error getting type from query", http.StatusBadRequest, response)
		return
	}
	if !validType(thetype) {
		handleError("invalid type in query", http.StatusBadRequest, response)
		return
	}
	postid := request.URL.Query().Get("postid")
	if postid == "" {
		handleError("error getting post id from query", http.StatusBadRequest, response)
		return
	}
	herostr := request.URL.Query().Get("hero")
	var hero bool
	if herostr == "" {
		hero = false
	} else if herostr == "true" {
		hero = true
	} else if herostr == "false" {
		hero = false
	} else {
		handleError("hero is not a boolean", http.StatusBadRequest, response)
		return
	}
	imageid := request.URL.Query().Get("imageid")
	if imageid == "" {
		if hero {
			imageid = "hero"
		} else {
			handleError("no hero and no picture id", http.StatusBadRequest, response)
			return
		}
	}
	var fileobj *storage.ObjectHandle
	if thetype == "blog" {
		fileobj = imageBucket.Object(blogImageIndex + "/" + postid + "/" + imageid)
	} else {
		fileobj = imageBucket.Object(projectImageIndex + "/" + postid + "/" + imageid)
	}
	filereader, err := fileobj.NewReader(ctxStorage)
	if err != nil {
		handleError("error reading file: "+err.Error(), http.StatusBadRequest, response)
		return
	}
	defer filereader.Close()
	filebuffer := new(bytes.Buffer)
	if bytesread, err := filebuffer.ReadFrom(filereader); err != nil {
		handleError("error reading to buffer: num bytes: "+strconv.FormatInt(bytesread, 10)+", "+err.Error(), http.StatusBadRequest, response)
		return
	}
	contentType := filereader.Attrs.ContentType
	response.Header().Set("content-type", contentType)
	response.Write(filebuffer.Bytes())
}
