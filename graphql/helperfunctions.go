package main

import (
	"encoding/binary"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

var dateFormat = "Mon Jan _2 15:04:05 2006"

func handleError(message string, statuscode int, response http.ResponseWriter) {
	// logger.Error(message)
	response.Header().Set("content-type", "application/json")
	response.WriteHeader(statuscode)
	response.Write([]byte(`{"message":"` + message + `"}`))
}

func objectidtimestamp(id primitive.ObjectID) time.Time {
	unixSecs := binary.BigEndian.Uint32(id[0:4])
	return time.Unix(int64(unixSecs), 0).UTC()
}
