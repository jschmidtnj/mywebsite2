package main

import (
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func createBlogPicture(response http.ResponseWriter, request *http.Request) {
	if !manageCors(response, request) {
		return
	}
	if request.Method != http.MethodPost {
		handleError("create blog picture http method not POST", http.StatusBadRequest, response)
		return
	}
	blogid := request.URL.Query().Get("blogid")
	if blogid == "" {
		handleError("error getting blog id from query", http.StatusBadRequest, response)
		return
	}
	herostr := request.URL.Query().Get("hero")
	var hero bool
	if herostr == "" {
		handleError("error getting hero from query", http.StatusBadRequest, response)
		return
	} else if herostr == "true" {
		hero = true
	} else if herostr == "false" {
		hero = false
	} else {
		handleError("hero is not a boolean", http.StatusBadRequest, response)
		return
	}
	if hero {
		logger.Info("hero is true")
	} else {
		logger.Info("hero is false")
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
	pictureid := uuid.New().String()
	fileobj := blogImageBucket.Object(pictureid)
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
	response.Write([]byte(`{"message":"file uploaded"}`))
}

func updateBlogPicture(response http.ResponseWriter, request *http.Request) {
	if !manageCors(response, request) {
		return
	}
	if request.Method != http.MethodPut {
		handleError("update blog picture http method not PUT", http.StatusBadRequest, response)
		return
	}
	blogid := request.URL.Query().Get("blogid")
	if blogid == "" {
		handleError("error getting blog id from query", http.StatusBadRequest, response)
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
	pictureid := request.URL.Query().Get("pictureid")
	if pictureid == "" && !hero {
		handleError("no hero and no picture id", http.StatusBadRequest, response)
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
	fileobj := blogImageBucket.Object(pictureid)
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

func deleteBlogPictures(response http.ResponseWriter, request *http.Request) {
	if !manageCors(response, request) {
		return
	}
	if request.Method != http.MethodDelete {
		handleError("delete blog picture http method not Delete", http.StatusBadRequest, response)
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
	if !(picturedata["pictureids"] != nil && picturedata["blogid"] != nil) {
		handleError("no pictureids or blogid provided", http.StatusBadRequest, response)
		return
	}
	pictureids, ok := picturedata["pictureids"].([]string)
	if !ok {
		handleError("pictureids cannot be cast to string array", http.StatusBadRequest, response)
		return
	}
	blogid, ok := picturedata["blogid"].(string)
	if !ok {
		handleError("blogid cannot be cast to string", http.StatusBadRequest, response)
		return
	}
	for _, pictureid := range pictureids {
		logger.Info("pictureid: " + pictureid + ", blogid: " + blogid)
		if err := blogImageBucket.Object(pictureid).Delete(ctxStorage); err != nil {
			handleError("error deleting file: "+err.Error(), http.StatusBadRequest, response)
			return
		}
	}
	response.Header().Set("content-type", "application/json")
	response.Write([]byte(`{"message":"files deleted"}`))
}
