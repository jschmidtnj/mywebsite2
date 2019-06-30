package main

import (
	"github.com/olivere/elastic"
	"net/http"
	"strconv"
)

func countPosts(response http.ResponseWriter, request *http.Request) {
	if !manageCors(&response, request) {
		return
	}
	if request.Method != http.MethodGet {
		handleError("register http method not Get", http.StatusBadRequest, response)
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
	searchterm := request.URL.Query().Get("imageid")
	var count int64
	var err error
	if len(searchterm) > 0 {
		queryString := elastic.NewQueryStringQuery(searchterm)
		count, err = elasticClient.Count().
			Type(thetype).
			Query(queryString).
			Pretty(false).
			Do(ctxElastic)
	} else {
		count, err = elasticClient.Count().
			Type(thetype).
			Query(nil).
			Pretty(false).
			Do(ctxElastic)
	}
	if err != nil {
		handleError(err.Error(), http.StatusBadRequest, response)
		return
	}
	response.Header().Set("content-type", "application/json")
	response.Write([]byte(`{"count":` + strconv.FormatInt(count, 10) + `}`))
}
