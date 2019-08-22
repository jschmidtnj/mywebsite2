package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"net/http"
	"time"

	"github.com/go-redis/redis"
	"github.com/ikeikeikeike/go-sitemap-generator/v2/stm"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var sitemapIndexTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<sitemapindex xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
  <sitemap>
    <loc>%s</loc>
    <lastmod>%s</lastmod>
  </sitemap>
  <sitemap>
    <loc>%s</loc>
    <lastmod>%s</lastmod>
	</sitemap>
	<sitemap>
    <loc>%s</loc>
    <lastmod>%s</lastmod>
  </sitemap>
</sitemapindex>`

func getSitemapIndex() ([]byte, error) {
	t := time.Now()
	currentTime := t.Format(sitemapTimeFormat)
	sitemap := fmt.Sprintf(sitemapIndexTemplate, websiteURL+"/sitemap-main.xml", lastSitemapUpdate,
		apiURL+"/sitemap-blogs.xml.gz", currentTime, apiURL+"/sitemap-projects.xml.gz", currentTime)
	return []byte(sitemap), nil
}

func checkCache(cachepath string) ([]byte, error) {
	if mode != "debug" {
		cachedresStr, err := redisClient.Get(cachepath).Result()
		if err != nil {
			if err != redis.Nil {
				return nil, err
			}
		} else {
			if len(cachedresStr) > 0 {
				return []byte(cachedresStr), nil
			}
		}
	}
	return nil, nil
}

func sitemapIndex(response http.ResponseWriter, request *http.Request) {
	var cachepath = "sitemap-index"
	cachedSitemap, err := checkCache(cachepath)
	if err != nil {
		handleError("error getting sitemap from cache: "+err.Error(), http.StatusBadRequest, response)
		return
	}
	response.Header().Set("Content-Type", "application/xml")
	if cachedSitemap != nil {
		response.Write(cachedSitemap)
		return
	}
	sitemap, err := getSitemapIndex()
	if err != nil {
		handleError("could not get sitemap: "+err.Error(), http.StatusBadRequest, response)
		return
	}
	err = redisClient.Set(cachepath, string(sitemap), cacheTime).Err()
	if err != nil {
		handleError("cannot update cache: "+err.Error(), http.StatusBadRequest, response)
		return
	}
	response.Write(sitemap)
}

func gZipData(data []byte) ([]byte, error) {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	_, err := gz.Write(data)
	if err != nil {
		return nil, err
	}
	if err = gz.Flush(); err != nil {
		return nil, err
	}
	if err = gz.Close(); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func sitemapIndexGZip(response http.ResponseWriter, request *http.Request) {
	var cachepath = "sitemap-index-gzip"
	cachedSitemap, err := checkCache(cachepath)
	if err != nil {
		handleError("error getting sitemap from cache: "+err.Error(), http.StatusBadRequest, response)
		return
	}
	response.Header().Set("Content-Type", "application/xml")
	response.Header().Set("Content-Encoding", "gzip")
	if cachedSitemap != nil {
		response.Write(cachedSitemap)
		return
	}
	sitemap, err := getSitemapIndex()
	if err != nil {
		handleError("could not get sitemap: "+err.Error(), http.StatusBadRequest, response)
		return
	}
	gzipdata, err := gZipData(sitemap)
	if err != nil {
		handleError("could not get gzip for index sitemap: "+err.Error(), http.StatusBadRequest, response)
		return
	}
	err = redisClient.Set(cachepath, string(gzipdata), cacheTime).Err()
	if err != nil {
		handleError("cannot update cache: "+err.Error(), http.StatusBadRequest, response)
		return
	}
	response.Write(gzipdata)
}

func setupSitemap(sitemap *stm.Sitemap) {
	sitemap.SetDefaultHost(websiteURL)
	sitemap.SetSitemapsHost(apiURL)
	sitemap.SetCompress(true)
	sitemap.SetVerbose(false)
	sitemap.Create()
}

func getSitemapBlogs() ([]byte, error) {
	sitemap := stm.NewSitemap(1)
	setupSitemap(sitemap)
	findOptions := options.Find()
	findOptions.SetProjection(bson.M{
		"_id": 1,
	})
	cursor, err := blogCollection.Find(ctxMongo, bson.D{{}})
	defer cursor.Close(ctxMongo)
	if err != nil {
		return nil, err
	}
	var foundstuff = false
	for cursor.Next(ctxMongo) {
		blogDataPrimitive := &bson.D{}
		err = cursor.Decode(blogDataPrimitive)
		if err != nil {
			return nil, err
		}
		blogData := blogDataPrimitive.Map()
		id := blogData["_id"].(primitive.ObjectID).Hex()
		sitemap.Add(stm.URL{{"loc", "/blog/" + id}})
		foundstuff = true
	}
	if !foundstuff {
		sitemap.Add(stm.URL{{"loc", "/blogs"}})
	}
	return sitemap.XMLContent(), nil
}

func sitemapBlogs(response http.ResponseWriter, request *http.Request) {
	var cachepath = "sitemap-blogs"
	cachedSitemap, err := checkCache(cachepath)
	if err != nil {
		handleError("error getting sitemap from cache: "+err.Error(), http.StatusBadRequest, response)
		return
	}
	response.Header().Set("Content-Type", "application/xml")
	if cachedSitemap != nil {
		response.Write(cachedSitemap)
		return
	}
	sitemap, err := getSitemapBlogs()
	if err != nil {
		handleError("error getting sitemap data: "+err.Error(), http.StatusBadRequest, response)
		return
	}
	err = redisClient.Set(cachepath, string(sitemap), cacheTime).Err()
	if err != nil {
		handleError("cannot update cache: "+err.Error(), http.StatusBadRequest, response)
		return
	}
	response.Write(sitemap)
}

func sitemapBlogsGZip(response http.ResponseWriter, request *http.Request) {
	var cachepath = "sitemap-blogs-gzip"
	cachedSitemap, err := checkCache(cachepath)
	if err != nil {
		handleError("error getting sitemap from cache: "+err.Error(), http.StatusBadRequest, response)
		return
	}
	response.Header().Set("Content-Type", "application/xml")
	response.Header().Set("Content-Encoding", "gzip")
	if cachedSitemap != nil {
		response.Write(cachedSitemap)
		return
	}
	sitemap, err := getSitemapBlogs()
	if err != nil {
		handleError("error getting sitemap data: "+err.Error(), http.StatusBadRequest, response)
		return
	}
	gzipdata, err := gZipData(sitemap)
	if err != nil {
		handleError("could not get gzip for blogs sitemap: "+err.Error(), http.StatusBadRequest, response)
		return
	}
	err = redisClient.Set(cachepath, string(gzipdata), cacheTime).Err()
	if err != nil {
		handleError("cannot update cache: "+err.Error(), http.StatusBadRequest, response)
		return
	}
	response.Write(gzipdata)
}

func getSitemapProjects() ([]byte, error) {
	sitemap := stm.NewSitemap(1)
	setupSitemap(sitemap)
	findOptions := options.Find()
	findOptions.SetProjection(bson.M{
		"_id": 1,
	})
	cursor, err := projectCollection.Find(ctxMongo, bson.D{{}})
	defer cursor.Close(ctxMongo)
	if err != nil {
		return nil, err
	}
	var foundstuff = false
	for cursor.Next(ctxMongo) {
		blogDataPrimitive := &bson.D{}
		err = cursor.Decode(blogDataPrimitive)
		if err != nil {
			return nil, err
		}
		blogData := blogDataPrimitive.Map()
		id := blogData["_id"].(primitive.ObjectID).Hex()
		sitemap.Add(stm.URL{{"loc", "/project/" + id}})
		foundstuff = true
	}
	if !foundstuff {
		sitemap.Add(stm.URL{{"loc", "/projects"}})
	}
	return sitemap.XMLContent(), nil
}

func sitemapProjects(response http.ResponseWriter, request *http.Request) {
	var cachepath = "sitemap-projects"
	cachedSitemap, err := checkCache(cachepath)
	if err != nil {
		handleError("error getting sitemap from cache: "+err.Error(), http.StatusBadRequest, response)
		return
	}
	response.Header().Set("Content-Type", "application/xml")
	if cachedSitemap != nil {
		response.Write(cachedSitemap)
		return
	}
	sitemap, err := getSitemapProjects()
	if err != nil {
		handleError("error getting sitemap data: "+err.Error(), http.StatusBadRequest, response)
		return
	}
	err = redisClient.Set(cachepath, string(sitemap), cacheTime).Err()
	if err != nil {
		handleError("cannot update cache: "+err.Error(), http.StatusBadRequest, response)
		return
	}
	response.Write(sitemap)
}

func sitemapProjectsGZip(response http.ResponseWriter, request *http.Request) {
	var cachepath = "sitemap-projects-gzip"
	cachedSitemap, err := checkCache(cachepath)
	if err != nil {
		handleError("error getting sitemap from cache: "+err.Error(), http.StatusBadRequest, response)
		return
	}
	response.Header().Set("Content-Type", "application/xml")
	response.Header().Set("Content-Encoding", "gzip")
	if cachedSitemap != nil {
		response.Write(cachedSitemap)
		return
	}
	sitemap, err := getSitemapProjects()
	if err != nil {
		handleError("error getting sitemap data: "+err.Error(), http.StatusBadRequest, response)
		return
	}
	gzipdata, err := gZipData(sitemap)
	if err != nil {
		handleError("could not get gzip for projects sitemap: "+err.Error(), http.StatusBadRequest, response)
		return
	}
	err = redisClient.Set(cachepath, string(gzipdata), cacheTime).Err()
	if err != nil {
		handleError("cannot update cache: "+err.Error(), http.StatusBadRequest, response)
		return
	}
	response.Write(gzipdata)
}
