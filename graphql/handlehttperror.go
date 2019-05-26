package main

import (
	"net/http"
)

func handleError(message string, statuscode int, response http.ResponseWriter) {
	// Logger.Error(message)
	response.Header().Set("content-type", "application/json")
	response.WriteHeader(statuscode)
	response.Write([]byte(`{"message":"` + message + `"}`))
}