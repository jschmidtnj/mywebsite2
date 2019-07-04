package main

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"github.com/olivere/elastic"
	"net/http"
)

/**
 * @api {get} /countPosts Count posts for search term
 * @apiVersion 0.0.1
 * @apiParam {String} searchterm Search term to count results
 * @apiParam {string="blog","project"} type Post type
 * @apiSuccess {String} count Result count
 * @apiGroup misc
 */
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
	response.Header().Set("content-type", "application/json")
	searchterm := request.URL.Query().Get("searchterm")
	pathMap := map[string]string{
		"path":       "count",
		"type":       thetype,
		"searchterm": searchterm,
	}
	cachepathBytes, err := json.Marshal(pathMap)
	if err != nil {
		handleError(err.Error(), http.StatusBadRequest, response)
		return
	}
	cachepath := string(cachepathBytes)
	cachedres, err := redisClient.Get(cachepath).Result()
	if err != nil {
		if err != redis.Nil {
			handleError(err.Error(), http.StatusBadRequest, response)
			return
		}
	} else {
		response.Write([]byte(cachedres))
		return
	}
	var count int64
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
	countMap := map[string]int64{
		"count": count,
	}
	countResBytes, err := json.Marshal(countMap)
	if err != nil {
		handleError(err.Error(), http.StatusBadRequest, response)
		return
	}
	err = redisClient.Set(cachepath, string(countResBytes), cacheTime).Err()
	if err != nil {
		handleError(err.Error(), http.StatusBadRequest, response)
		return
	}
	response.Write(countResBytes)
}
