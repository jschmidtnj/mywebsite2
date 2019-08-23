package main

import (
	"bytes"
	"encoding/json"
	"image"
	"image/jpeg"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
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
	postid := request.URL.Query().Get("postid")
	if postid == "" {
		handleError("error getting post id from query", http.StatusBadRequest, response)
		return
	}
	fileid := request.URL.Query().Get("fileid")
	if fileid == "" {
		handleError("error getting file id from query", http.StatusBadRequest, response)
		return
	}
	fileidDecoded, err := url.QueryUnescape(fileid)
	if err != nil {
		handleError(err.Error(), http.StatusBadRequest, response)
		return
	}
	var filebuffer bytes.Buffer
	file, _, err := request.FormFile("file")
	if err != nil {
		handleError(err.Error(), http.StatusBadRequest, response)
		return
	}
	defer file.Close()
	io.Copy(&filebuffer, file)
	defer filebuffer.Reset()
	var fileobj *storage.ObjectHandle
	if thetype == "blog" {
		fileobj = storageBucket.Object(blogFileIndex + "/" + postid + "/" + fileidDecoded)
	} else {
		fileobj = storageBucket.Object(projectFileIndex + "/" + postid + "/" + fileidDecoded)
	}
	filewriter := fileobj.NewWriter(ctxStorage)
	errmessage := uploadFile(&filebuffer, filewriter)
	if len(errmessage) > 0 {
		handleError(errmessage, http.StatusBadRequest, response)
		return
	}
	response.Header().Set("Content-Type", "application/json")
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
	postid := request.URL.Query().Get("postid")
	if postid == "" {
		handleError("error getting post id from query", http.StatusBadRequest, response)
		return
	}
	imageid := request.URL.Query().Get("imageid")
	if imageid == "" {
		handleError("error getting image id from query", http.StatusBadRequest, response)
		return
	}
	imageidDecoded, err := url.QueryUnescape(imageid)
	if err != nil {
		handleError(err.Error(), http.StatusBadRequest, response)
		return
	}
	file, _, err := request.FormFile("file")
	if err != nil {
		handleError(err.Error(), http.StatusBadRequest, response)
		return
	}
	defer file.Close()
	originalImage, _, err := image.Decode(file)
	if err != nil {
		handleError(err.Error(), http.StatusBadRequest, response)
		return
	}
	originalImageBuffer := new(bytes.Buffer)
	defer originalImageBuffer.Reset()
	jpegOptionsOriginal := jpeg.Options{Quality: 90}
	err = jpeg.Encode(originalImageBuffer, originalImage, &jpegOptionsOriginal)
	if err != nil {
		handleError(err.Error(), http.StatusBadRequest, response)
		return
	}
	blurredImage := imaging.Blur(originalImage, progressiveImageBlurAmount)
	blurredImageBuffer := new(bytes.Buffer)
	defer blurredImageBuffer.Reset()
	jpegOptionsBlurred := jpeg.Options{Quality: 60}
	err = jpeg.Encode(blurredImageBuffer, blurredImage, &jpegOptionsBlurred)
	if err != nil {
		handleError(err.Error(), http.StatusBadRequest, response)
		return
	}
	var originalImageObj *storage.ObjectHandle
	var blurredImageObj *storage.ObjectHandle
	if thetype == "blog" {
		originalImageObj = storageBucket.Object(blogImageIndex + "/" + postid + "/" + imageidDecoded + "/original")
		blurredImageObj = storageBucket.Object(blogImageIndex + "/" + postid + "/" + imageidDecoded + "/blur")
	} else {
		originalImageObj = storageBucket.Object(projectImageIndex + "/" + postid + "/" + imageidDecoded + "/original")
		blurredImageObj = storageBucket.Object(projectImageIndex + "/" + postid + "/" + imageidDecoded + "/blur")
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
	response.Header().Set("Content-Type", "application/json")
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
	postid, ok := picturedata["postid"].(string)
	if !ok {
		handleError("unable to post id to string", http.StatusBadRequest, response)
		return
	}
	if !validType(thetype) {
		handleError("invalid type in body", http.StatusBadRequest, response)
		return
	}
	imageids, ok := picturedata["imageids"].([]interface{})
	if !ok {
		handleError("imageids cannot be cast to interface array", http.StatusBadRequest, response)
		return
	}
	for _, imageidinterface := range imageids {
		imageid, ok := imageidinterface.(string)
		if !ok {
			handleError("imageid cannot be cast to string", http.StatusBadRequest, response)
			return
		}
		var imageobjblur *storage.ObjectHandle
		var imageobjoriginal *storage.ObjectHandle
		if thetype == "blog" {
			imageobjblur = storageBucket.Object(blogImageIndex + "/" + postid + "/" + imageid + "/blur")
			imageobjoriginal = storageBucket.Object(blogImageIndex + "/" + postid + "/" + imageid + "/original")
		} else {
			imageobjblur = storageBucket.Object(projectImageIndex + "/" + postid + "/" + imageid + "/blur")
			imageobjoriginal = storageBucket.Object(projectImageIndex + "/" + postid + "/" + imageid + "/original")
		}
		if err := imageobjblur.Delete(ctxStorage); err != nil {
			handleError("error deleting image: "+err.Error(), http.StatusBadRequest, response)
			return
		}
		if err := imageobjoriginal.Delete(ctxStorage); err != nil {
			handleError("error deleting image: "+err.Error(), http.StatusBadRequest, response)
			return
		}
	}
	response.Header().Set("Content-Type", "application/json")
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
	postid, ok := filedata["postid"].(string)
	if !ok {
		handleError("unable to cast post id to string", http.StatusBadRequest, response)
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
			fileobj = storageBucket.Object(blogFileIndex + "/" + postid + "/" + fileid)
		} else {
			fileobj = storageBucket.Object(projectFileIndex + "/" + postid + "/" + fileid)
		}
		if err := fileobj.Delete(ctxStorage); err != nil {
			handleError("error deleting file: "+err.Error(), http.StatusBadRequest, response)
			return
		}
	}
	response.Header().Set("Content-Type", "application/json")
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
	postid := request.URL.Query().Get("postid")
	if postid == "" {
		handleError("no post id", http.StatusBadRequest, response)
		return
	}
	var imageobj *storage.ObjectHandle
	if thetype == "blog" {
		imageobj = storageBucket.Object(blogImageIndex + "/" + postid + "/" + imageid)
	} else {
		imageobj = storageBucket.Object(projectImageIndex + "/" + postid + "/" + imageid)
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
	postid := request.URL.Query().Get("postid")
	if postid == "" {
		handleError("error getting post id from query", http.StatusBadRequest, response)
		return
	}
	fileid := request.URL.Query().Get("fileid")
	if fileid == "" {
		handleError("no file id", http.StatusBadRequest, response)
		return
	}
	var fileobj *storage.ObjectHandle
	if thetype == "blog" {
		fileobj = storageBucket.Object(blogImageIndex + "/" + postid + "/" + fileid)
	} else {
		fileobj = storageBucket.Object(projectImageIndex + "/" + postid + "/" + fileid)
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
