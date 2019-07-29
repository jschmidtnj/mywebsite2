package main

import (
	"bytes"
	"encoding/json"
	"image"
	"image/jpeg"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"cloud.google.com/go/storage"
	"github.com/disintegration/imaging"
)

func uploadFile(fileBuffer *bytes.Buffer, filewriter *storage.Writer) string {
	byteswritten, err := fileBuffer.WriteTo(filewriter)
	if err != nil {
		return "error writing to filewriter: num bytes: " + strconv.FormatInt(byteswritten, 10) + ", " + err.Error()
	}
	err = filewriter.Close()
	if err != nil {
		return "error closing writer: " + err.Error()
	}
	return ""
}

func writePostFile(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPut {
		handleError("upload file http method not PUT", http.StatusBadRequest, response)
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
	fileid := request.URL.Query().Get("fileid")
	if fileid == "" {
		handleError("error getting file id from query", http.StatusBadRequest, response)
		return
	}
	var filebuffer bytes.Buffer
	file, _, err := request.FormFile("file")
	if err != nil {
		handleError(err.Error(), http.StatusBadRequest, response)
	}
	defer file.Close()
	io.Copy(&filebuffer, file)
	defer filebuffer.Reset()
	var fileobj *storage.ObjectHandle
	if thetype == "blog" {
		fileobj = storageBucket.Object(blogFileIndex + "/" + fileid)
	} else {
		fileobj = storageBucket.Object(projectFileIndex + "/" + fileid)
	}
	filewriter := fileobj.NewWriter(ctxStorage)
	errmessage := uploadFile(&filebuffer, filewriter)
	if len(errmessage) > 0 {
		handleError(errmessage, http.StatusBadRequest, response)
		return
	}
	response.Header().Set("content-type", "application/json")
	response.Write([]byte(`{"message":"file updated"}`))
}

func writePostPicture(response http.ResponseWriter, request *http.Request) {

	if request.Method != http.MethodPut {
		handleError("post picture http method not PUT", http.StatusBadRequest, response)
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
	imageid := request.URL.Query().Get("imageid")
	if imageid == "" {
		handleError("error getting image id from query", http.StatusBadRequest, response)
		return
	}
	file, _, err := request.FormFile("file")
	if err != nil {
		handleError(err.Error(), http.StatusBadRequest, response)
	}
	defer file.Close()
	originalImage, _, err := image.Decode(file)
	if err != nil {
		handleError(err.Error(), http.StatusBadRequest, response)
	}
	originalImageBuffer := new(bytes.Buffer)
	defer originalImageBuffer.Reset()
	jpegOptionsOriginal := jpeg.Options{Quality: 90}
	err = jpeg.Encode(originalImageBuffer, originalImage, &jpegOptionsOriginal)
	if err != nil {
		handleError(err.Error(), http.StatusBadRequest, response)
	}
	blurredImage := imaging.Blur(originalImage, progressiveImageBlurAmount)
	blurredImageBuffer := new(bytes.Buffer)
	defer blurredImageBuffer.Reset()
	jpegOptionsBlurred := jpeg.Options{Quality: 60}
	err = jpeg.Encode(blurredImageBuffer, blurredImage, &jpegOptionsBlurred)
	if err != nil {
		handleError(err.Error(), http.StatusBadRequest, response)
	}
	var originalImageObj *storage.ObjectHandle
	var blurredImageObj *storage.ObjectHandle
	if thetype == "blog" {
		originalImageObj = storageBucket.Object(blogImageIndex + "/" + imageid + "/original")
		blurredImageObj = storageBucket.Object(blogImageIndex + "/" + imageid + "/blur")
	} else {
		originalImageObj = storageBucket.Object(projectImageIndex + "/" + imageid + "/original")
		blurredImageObj = storageBucket.Object(blogImageIndex + "/" + imageid + "/blur")
	}
	originalImageWriter := originalImageObj.NewWriter(ctxStorage)
	errmessage := uploadFile(originalImageBuffer, originalImageWriter)
	if len(errmessage) > 0 {
		handleError(errmessage, http.StatusBadRequest, response)
		return
	}
	blurredImageWriter := blurredImageObj.NewWriter(ctxStorage)
	errmessage = uploadFile(blurredImageBuffer, blurredImageWriter)
	if len(errmessage) > 0 {
		handleError(errmessage, http.StatusBadRequest, response)
		return
	}
	response.Header().Set("content-type", "application/json")
	response.Write([]byte(`{"message":"image written"}`))
}

func deletePostPictures(response http.ResponseWriter, request *http.Request) {

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
	imageids, ok := picturedata["imageids"].([]string)
	if !ok {
		handleError("imageids cannot be cast to string array", http.StatusBadRequest, response)
		return
	}
	for _, imageid := range imageids {
		var imageobj *storage.ObjectHandle
		if thetype == "blog" {
			imageobj = storageBucket.Object(blogImageIndex + "/" + imageid)
		} else {
			imageobj = storageBucket.Object(projectImageIndex + "/" + imageid)
		}
		if err := imageobj.Delete(ctxStorage); err != nil {
			handleError("error deleting image: "+err.Error(), http.StatusBadRequest, response)
			return
		}
	}
	response.Header().Set("content-type", "application/json")
	response.Write([]byte(`{"message":"files deleted"}`))
}

func deletePostFiles(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodDelete {
		handleError("delete post files http method not Delete", http.StatusBadRequest, response)
		return
	}
	if _, err := validateAdmin(getAuthToken(request)); err != nil {
		handleError("auth error: "+err.Error(), http.StatusBadRequest, response)
		return
	}
	var filedata map[string]interface{}
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		handleError("error getting request body: "+err.Error(), http.StatusBadRequest, response)
		return
	}
	err = json.Unmarshal(body, &filedata)
	if err != nil {
		handleError("error parsing request body: "+err.Error(), http.StatusBadRequest, response)
		return
	}
	if !(filedata["fileids"] != nil && filedata["postid"] != nil && filedata["type"] != nil) {
		handleError("no fileids or postid or type provided", http.StatusBadRequest, response)
		return
	}
	thetype, ok := filedata["type"].(string)
	if !ok {
		handleError("unable to cast type to string", http.StatusBadRequest, response)
		return
	}
	if !validType(thetype) {
		handleError("invalid type in body", http.StatusBadRequest, response)
		return
	}
	fileids, ok := filedata["fileids"].([]string)
	if !ok {
		handleError("fileids cannot be cast to string array", http.StatusBadRequest, response)
		return
	}
	for _, fileid := range fileids {
		var fileobj *storage.ObjectHandle
		if thetype == "blog" {
			fileobj = storageBucket.Object(blogFileIndex + "/" + fileid)
		} else {
			fileobj = storageBucket.Object(projectFileIndex + "/" + fileid)
		}
		if err := fileobj.Delete(ctxStorage); err != nil {
			handleError("error deleting file: "+err.Error(), http.StatusBadRequest, response)
			return
		}
	}
	response.Header().Set("content-type", "application/json")
	response.Write([]byte(`{"message":"files deleted"}`))
}

func getPostPicture(response http.ResponseWriter, request *http.Request) {
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
	imageid := request.URL.Query().Get("imageid")
	if imageid == "" {
		handleError("no picture id", http.StatusBadRequest, response)
		return
	}
	var imageobj *storage.ObjectHandle
	if thetype == "blog" {
		imageobj = storageBucket.Object(blogImageIndex + "/" + imageid)
	} else {
		imageobj = storageBucket.Object(projectImageIndex + "/" + imageid)
	}
	imagereader, err := imageobj.NewReader(ctxStorage)
	if err != nil {
		handleError("error reading image: "+err.Error(), http.StatusBadRequest, response)
		return
	}
	defer imagereader.Close()
	imagebuffer := new(bytes.Buffer)
	if bytesread, err := imagebuffer.ReadFrom(imagereader); err != nil {
		handleError("error reading to buffer: num bytes: "+strconv.FormatInt(bytesread, 10)+", "+err.Error(), http.StatusBadRequest, response)
		return
	}
	contentType := imagereader.Attrs.ContentType
	response.Header().Set("Content-Type", contentType)
	response.Write(imagebuffer.Bytes())
}

func getPostFile(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		handleError("get post file http method not GET", http.StatusBadRequest, response)
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
	fileid := request.URL.Query().Get("fileid")
	if fileid == "" {
		handleError("no file id", http.StatusBadRequest, response)
		return
	}
	var fileobj *storage.ObjectHandle
	if thetype == "blog" {
		fileobj = storageBucket.Object(blogImageIndex + "/" + fileid)
	} else {
		fileobj = storageBucket.Object(projectImageIndex + "/" + fileid)
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
	response.Header().Set("Content-Type", contentType)
	response.Write(filebuffer.Bytes())
}
