package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"image"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"cloud.google.com/go/storage"
	"github.com/disintegration/imaging"
)

func validateContentType(thetype string) error {
	for _, validtype := range validContentTypes {
		if validtype == thetype {
			return nil
		}
	}
	return errors.New("invalid content type provided")
}

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

func writeGenericFile(file io.Reader, filetype string, posttype string, fileidDecoded string, postid string) error {
	var filebuffer bytes.Buffer
	io.Copy(&filebuffer, file)
	defer filebuffer.Reset()
	var fileobj *storage.ObjectHandle
	if posttype == blogType {
		fileobj = storageBucket.Object(blogFileIndex + "/" + postid + "/" + fileidDecoded + originalPath)
	} else {
		fileobj = storageBucket.Object(projectFileIndex + "/" + postid + "/" + fileidDecoded + originalPath)
	}
	filewriter := fileobj.NewWriter(ctxStorage)
	filewriter.ContentType = filetype
	filewriter.Metadata = map[string]string{}
	errmessage := uploadFile(&filebuffer, filewriter)
	if len(errmessage) > 0 {
		return errors.New(errmessage)
	}
	return nil
}

func writeJpeg(file io.Reader, filetype string, posttype string, fileidDecoded string, postid string) error {
	originalImage, _, err := image.Decode(file)
	if err != nil {
		return err
	}
	originalImageBuffer := new(bytes.Buffer)
	defer originalImageBuffer.Reset()
	jpegOptionsOriginal := jpeg.Options{Quality: 90}
	err = jpeg.Encode(originalImageBuffer, originalImage, &jpegOptionsOriginal)
	if err != nil {
		return err
	}
	blurredImage := imaging.Blur(originalImage, progressiveImageBlurAmount)
	blurredImageBuffer := new(bytes.Buffer)
	defer blurredImageBuffer.Reset()
	jpegOptionsBlurred := jpeg.Options{Quality: 60}
	err = jpeg.Encode(blurredImageBuffer, blurredImage, &jpegOptionsBlurred)
	if err != nil {
		return err
	}
	var originalImageObj *storage.ObjectHandle
	var blurredImageObj *storage.ObjectHandle
	if posttype == blogType {
		originalImageObj = storageBucket.Object(blogFileIndex + "/" + postid + "/" + fileidDecoded + originalPath)
		blurredImageObj = storageBucket.Object(blogFileIndex + "/" + postid + "/" + fileidDecoded + blurPath)
	} else {
		originalImageObj = storageBucket.Object(projectFileIndex + "/" + postid + "/" + fileidDecoded + originalPath)
		blurredImageObj = storageBucket.Object(projectFileIndex + "/" + postid + "/" + fileidDecoded + blurPath)
	}
	originalImageWriter := originalImageObj.NewWriter(ctxStorage)
	originalImageWriter.ContentType = filetype
	originalImageWriter.Metadata = map[string]string{}
	errmessage := uploadFile(originalImageBuffer, originalImageWriter)
	if len(errmessage) > 0 {
		return errors.New(errmessage)
	}
	blurredImageWriter := blurredImageObj.NewWriter(ctxStorage)
	blurredImageWriter.ContentType = filetype
	blurredImageWriter.Metadata = map[string]string{}
	errmessage = uploadFile(blurredImageBuffer, blurredImageWriter)
	if len(errmessage) > 0 {
		return errors.New(errmessage)
	}
	return nil
}

func writePng(file io.Reader, filetype string, posttype string, fileidDecoded string, postid string) error {
	originalImage, _, err := image.Decode(file)
	if err != nil {
		return err
	}
	originalImageBuffer := new(bytes.Buffer)
	defer originalImageBuffer.Reset()
	err = png.Encode(originalImageBuffer, originalImage)
	if err != nil {
		return err
	}
	blurredImage := imaging.Blur(originalImage, progressiveImageBlurAmount)
	blurredImageBuffer := new(bytes.Buffer)
	defer blurredImageBuffer.Reset()
	err = png.Encode(blurredImageBuffer, blurredImage)
	if err != nil {
		return err
	}
	var originalImageObj *storage.ObjectHandle
	var blurredImageObj *storage.ObjectHandle
	if posttype == blogType {
		originalImageObj = storageBucket.Object(blogFileIndex + "/" + postid + "/" + fileidDecoded + originalPath)
		blurredImageObj = storageBucket.Object(blogFileIndex + "/" + postid + "/" + fileidDecoded + blurPath)
	} else {
		originalImageObj = storageBucket.Object(projectFileIndex + "/" + postid + "/" + fileidDecoded + originalPath)
		blurredImageObj = storageBucket.Object(projectFileIndex + "/" + postid + "/" + fileidDecoded + blurPath)
	}
	originalImageWriter := originalImageObj.NewWriter(ctxStorage)
	originalImageWriter.ContentType = filetype
	originalImageWriter.Metadata = map[string]string{}
	errmessage := uploadFile(originalImageBuffer, originalImageWriter)
	if len(errmessage) > 0 {
		return errors.New(errmessage)
	}
	blurredImageWriter := blurredImageObj.NewWriter(ctxStorage)
	blurredImageWriter.ContentType = filetype
	blurredImageWriter.Metadata = map[string]string{}
	errmessage = uploadFile(blurredImageBuffer, blurredImageWriter)
	if len(errmessage) > 0 {
		return errors.New(errmessage)
	}
	return nil
}

func getGifDimensions(gif *gif.GIF) (x, y int) {
	var lowestX int
	var lowestY int
	var highestX int
	var highestY int
	for _, img := range gif.Image {
		if img.Rect.Min.X < lowestX {
			lowestX = img.Rect.Min.X
		}
		if img.Rect.Min.Y < lowestY {
			lowestY = img.Rect.Min.Y
		}
		if img.Rect.Max.X > highestX {
			highestX = img.Rect.Max.X
		}
		if img.Rect.Max.Y > highestY {
			highestY = img.Rect.Max.Y
		}
	}
	return highestX - lowestX, highestY - lowestY
}

func writeGif(file io.Reader, filetype string, posttype string, fileidDecoded string, postid string) error {
	originalGif, err := gif.DecodeAll(file)
	if err != nil {
		return err
	}
	originalGifBuffer := new(bytes.Buffer)
	defer originalGifBuffer.Reset()
	err = gif.EncodeAll(originalGifBuffer, originalGif)
	if err != nil {
		return err
	}
	imgWidth, imgHeight := getGifDimensions(originalGif)
	originalImage := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))
	draw.Draw(originalImage, originalImage.Bounds(), originalGif.Image[0], image.ZP, draw.Src)
	originalImageBuffer := new(bytes.Buffer)
	defer originalImageBuffer.Reset()
	jpegOptionsOriginal := jpeg.Options{Quality: 90}
	err = jpeg.Encode(originalImageBuffer, originalImage, &jpegOptionsOriginal)
	if err != nil {
		return err
	}
	blurredImage := imaging.Blur(originalImage, progressiveImageBlurAmount)
	blurredImageBuffer := new(bytes.Buffer)
	defer blurredImageBuffer.Reset()
	jpegOptionsBlurred := jpeg.Options{Quality: 60}
	err = jpeg.Encode(blurredImageBuffer, blurredImage, &jpegOptionsBlurred)
	if err != nil {
		return err
	}
	var originalGifObj *storage.ObjectHandle
	var originalImageObj *storage.ObjectHandle
	var blurredImageObj *storage.ObjectHandle
	if posttype == blogType {
		originalGifObj = storageBucket.Object(blogFileIndex + "/" + postid + "/" + fileidDecoded + originalPath)
		originalImageObj = storageBucket.Object(blogFileIndex + "/" + postid + "/" + fileidDecoded + placeholderPath + originalPath)
		blurredImageObj = storageBucket.Object(blogFileIndex + "/" + postid + "/" + fileidDecoded + placeholderPath + blurPath)
	} else {
		originalGifObj = storageBucket.Object(blogFileIndex + "/" + postid + "/" + fileidDecoded + originalPath)
		originalImageObj = storageBucket.Object(projectFileIndex + "/" + postid + "/" + fileidDecoded + placeholderPath + originalPath)
		blurredImageObj = storageBucket.Object(projectFileIndex + "/" + postid + "/" + fileidDecoded + placeholderPath + blurPath)
	}
	originalGifWriter := originalGifObj.NewWriter(ctxStorage)
	originalGifWriter.ContentType = filetype
	originalGifWriter.Metadata = map[string]string{}
	errmessage := uploadFile(originalGifBuffer, originalGifWriter)
	if len(errmessage) > 0 {
		return errors.New(errmessage)
	}
	var placeholderFileType = "image/jpeg"
	originalImageWriter := originalImageObj.NewWriter(ctxStorage)
	originalImageWriter.ContentType = placeholderFileType
	originalImageWriter.Metadata = map[string]string{}
	errmessage = uploadFile(originalImageBuffer, originalImageWriter)
	if len(errmessage) > 0 {
		return errors.New(errmessage)
	}
	blurredImageWriter := blurredImageObj.NewWriter(ctxStorage)
	blurredImageWriter.ContentType = placeholderFileType
	blurredImageWriter.Metadata = map[string]string{}
	errmessage = uploadFile(blurredImageBuffer, blurredImageWriter)
	if len(errmessage) > 0 {
		return errors.New(errmessage)
	}
	return nil
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
	filetype := request.URL.Query().Get("filetype")
	if filetype == "" {
		handleError("error getting filetype from query", http.StatusBadRequest, response)
		return
	}
	if err := validateContentType(filetype); err != nil {
		handleError(err.Error(), http.StatusBadRequest, response)
		return
	}
	posttype := request.URL.Query().Get("posttype")
	if posttype == "" {
		handleError("error getting posttype from query", http.StatusBadRequest, response)
		return
	}
	if !validType(posttype) {
		handleError("invalid posttype in query", http.StatusBadRequest, response)
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
	file, _, err := request.FormFile("file")
	if err != nil {
		handleError(err.Error(), http.StatusBadRequest, response)
		return
	}
	defer file.Close()
	switch filetype {
	case "image/jpeg":
		if err = writeJpeg(file, filetype, posttype, fileidDecoded, postid); err != nil {
			handleError(err.Error(), http.StatusBadRequest, response)
			return
		}
		break
	case "image/png":
		if err = writePng(file, filetype, posttype, fileidDecoded, postid); err != nil {
			handleError(err.Error(), http.StatusBadRequest, response)
			return
		}
		break
	case "image/gif":
		if err = writeGif(file, filetype, posttype, fileidDecoded, postid); err != nil {
			handleError(err.Error(), http.StatusBadRequest, response)
			return
		}
		break
	default:
		if err = writeGenericFile(file, filetype, posttype, fileidDecoded, postid); err != nil {
			handleError(err.Error(), http.StatusBadRequest, response)
			return
		}
		break
	}
	response.Header().Set("Content-Type", "application/json")
	response.Write([]byte(`{"message":"file updated"}`))
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
	if !(filedata["fileids"] != nil && filedata["postid"] != nil && filedata["posttype"] != nil) {
		handleError("no fileids or postid or type provided", http.StatusBadRequest, response)
		return
	}
	postid, ok := filedata["postid"].(string)
	if !ok {
		handleError("unable to cast post id to string", http.StatusBadRequest, response)
		return
	}
	posttype, ok := filedata["posttype"].(string)
	if !ok {
		handleError("unable to cast posttype to string", http.StatusBadRequest, response)
		return
	}
	if !validType(posttype) {
		handleError("invalid posttype in body", http.StatusBadRequest, response)
		return
	}
	fileids, ok := filedata["fileids"].([]interface{})
	if !ok {
		handleError("file ids cannot be cast to interface array", http.StatusBadRequest, response)
		return
	}
	for _, fileidinterface := range fileids {
		fileid, ok := fileidinterface.(string)
		if !ok {
			handleError("file id cannot be cast to string", http.StatusBadRequest, response)
			return
		}
		var fileobj *storage.ObjectHandle
		if posttype == blogType {
			fileobj = storageBucket.Object(blogFileIndex + "/" + postid + "/" + fileid + originalPath)
		} else {
			fileobj = storageBucket.Object(projectFileIndex + "/" + postid + "/" + fileid + originalPath)
		}
		fileobjattributes, err := fileobj.Attrs(ctxStorage)
		if err != nil {
			handleError(err.Error(), http.StatusBadRequest, response)
			return
		}
		var filetype = fileobjattributes.ContentType
		var hasblur = false
		for _, blurtype := range haveblur {
			if blurtype == filetype {
				hasblur = true
				break
			}
		}
		if err := fileobj.Delete(ctxStorage); err != nil {
			handleError("error deleting original file: "+err.Error(), http.StatusBadRequest, response)
			return
		}
		if hasblur {
			if posttype == blogType {
				fileobj = storageBucket.Object(blogFileIndex + "/" + postid + "/" + fileid + blurPath)
			} else {
				fileobj = storageBucket.Object(projectFileIndex + "/" + postid + "/" + fileid + blurPath)
			}
			if err := fileobj.Delete(ctxStorage); err != nil {
				handleError("error deleting blur file: "+err.Error(), http.StatusBadRequest, response)
				return
			}
		}
	}
	response.Header().Set("Content-Type", "application/json")
	response.Write([]byte(`{"message":"files deleted"}`))
}

func getPostFile(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		handleError("get post file http method not GET", http.StatusBadRequest, response)
		return
	}
	posttype := request.URL.Query().Get("posttype")
	if posttype == "" {
		handleError("error getting posttype from query", http.StatusBadRequest, response)
		return
	}
	if !validType(posttype) {
		handleError("invalid posttype in query", http.StatusBadRequest, response)
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
	if posttype == blogType {
		fileobj = storageBucket.Object(blogFileIndex + "/" + postid + "/" + fileid + originalPath)
	} else {
		fileobj = storageBucket.Object(projectFileIndex + "/" + postid + "/" + fileid + originalPath)
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
