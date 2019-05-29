package main

import (
	"net/http"
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"encoding/binary"
)

var DateFormat = "Mon Jan _2 15:04:05 2006"

func handleError(message string, statuscode int, response http.ResponseWriter) {
	// Logger.Error(message)
	response.Header().Set("content-type", "application/json")
	response.WriteHeader(statuscode)
	response.Write([]byte(`{"message":"` + message + `"}`))
}

func objectidtimestamp(id primitive.ObjectID) time.Time {
	unixSecs := binary.BigEndian.Uint32(id[0:4])
	return time.Unix(int64(unixSecs), 0).UTC()
}
