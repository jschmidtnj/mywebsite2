package main

import (
	"cloud.google.com/go/storage"
	"context"
	"encoding/json"
	"github.com/graphql-go/graphql"
	"github.com/joho/godotenv"
	// medium "github.com/medium/medium-sdk-go"
	"github.com/olivere/elastic"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"google.golang.org/api/option"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var jwtSecret []byte

var tokenExpiration int

var sendgridAPIKey string

var websiteURL string

var mongoClient *mongo.Client

var ctxMongo context.Context

var database = "testing"

var userCollection *mongo.Collection

var userMongoName = "users"

var blogCollection *mongo.Collection

var blogMongoName = "blogs"

var projectCollection *mongo.Collection

var projectMongoName = "projects"

var elasticClient *elastic.Client

var ctxElastic context.Context

var blogElasticIndex = "blogs"

var blogElasticType = "blog"

var projectElasticIndex = "projects"

var projectElasticType = "project"

var validTypes = []string{
	"blog",
	"project",
}

var ctxStorage context.Context

var storageClient *storage.Client

var imageBucket *storage.BucketHandle

var blogImageIndex = "blogs"

var projectImageIndex = "projects"

var logger *zap.Logger

type tokenKeyType string

var tokenKey tokenKeyType

// var mediumClient *medium.Medium

// var mediumUser *medium.User

func hello(response http.ResponseWriter, request *http.Request) {
	if !manageCors(&response, request) {
		return
	}
	response.Header().Set("content-type", "application/json")
	response.Write([]byte(`{"message":"Hello!"}`))
}

func main() {
	// "./logs"
	loggerconfig := []byte(`{
		"level": "debug",
		"encoding": "json",
		"outputPaths": ["stdout"],
		"errorOutputPaths": ["stderr"],
		"initialFields": {},
		"encoderConfig": {
		  "messageKey": "message",
		  "levelKey": "level",
		  "levelEncoder": "lowercase"
		}
  }`)
	var zapconfig zap.Config
	if err := json.Unmarshal(loggerconfig, &zapconfig); err != nil {
		panic(err)
	}
	var err error
	logger, err = zapconfig.Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	logger.Info("logger created")
	err = godotenv.Load()
	if err != nil {
		logger.Fatal("Error loading .env file")
	}
	jwtSecret = []byte(os.Getenv("SECRET"))
	tokenExpiration, err = strconv.Atoi(os.Getenv("TOKENEXPIRATION"))
	if err != nil {
		logger.Fatal(err.Error())
	}
	sendgridAPIKey = os.Getenv("SENDGRIDAPIKEY")
	websiteURL = os.Getenv("WEBSITEURL")
	ctxMongo, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	cancel()
	mongouri := os.Getenv("MONGOURI")
	mongoClient, err = mongo.Connect(ctxMongo, options.Client().ApplyURI(mongouri))
	if err != nil {
		logger.Fatal(err.Error())
	}
	userCollection = mongoClient.Database(database).Collection(userMongoName)
	projectCollection = mongoClient.Database(database).Collection(projectMongoName)
	blogCollection = mongoClient.Database(database).Collection(blogMongoName)
	elasticuri := os.Getenv("ELASTICURI")
	elasticClient, err = elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(elasticuri))
	if err != nil {
		logger.Fatal(err.Error())
	}
	ctxElastic = context.Background()
	var storageconfigstr = os.Getenv("STORAGECONFIG")
	var storageconfigjson map[string]interface{}
	json.Unmarshal([]byte(storageconfigstr), &storageconfigjson)
	ctxStorage = context.Background()
	storageconfigjsonbytes, err := json.Marshal(storageconfigjson)
	if err != nil {
		logger.Fatal(err.Error())
	}
	storageClient, err = storage.NewClient(ctxStorage, option.WithCredentialsJSON([]byte(storageconfigjsonbytes)))
	if err != nil {
		logger.Fatal(err.Error())
	}
	bucketName := os.Getenv("STORAGEBUCKETNAME")
	imageBucket = storageClient.Bucket(bucketName)
	gcpprojectid, ok := storageconfigjson["project_id"].(string)
	if !ok {
		logger.Fatal("could not cast gcp project id to string")
	}
	if err := imageBucket.Create(ctxStorage, gcpprojectid, nil); err != nil {
		logger.Info(err.Error())
	}
	/*
		mediumAccessToken := os.Getenv("MEDIUMACCESSTOKEN")
		mediumClient = medium.NewClientWithAccessToken(mediumAccessToken)
		mediumUser, err := mediumClient.GetUser("")
		if err != nil {
			logger.Fatal(err.Error())
		}
		logger.Info("medium id " + mediumUser.ID)
	*/
	port := ":" + os.Getenv("PORT")
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:    rootQuery(),
		Mutation: rootMutation(),
	})
	if err != nil {
		logger.Fatal(err.Error())
	}
	http.HandleFunc("/graphql", func(response http.ResponseWriter, request *http.Request) {
		if !manageCors(&response, request) {
			return
		}
		tokenKey = tokenKeyType("token")
		result := graphql.Do(graphql.Params{
			Schema:        schema,
			RequestString: request.URL.Query().Get("query"),
			Context:       context.WithValue(context.Background(), tokenKey, getAuthToken(request)),
		})
		response.Header().Set("content-type", "application/json")
		json.NewEncoder(response).Encode(result)
	})
	http.HandleFunc("/countPosts", countPosts)
	http.HandleFunc("/sendTestEmail", sendTestEmail)
	http.HandleFunc("/loginEmailPassword", loginEmailPassword)
	http.HandleFunc("/logoutEmailPassword", logoutEmailPassword)
	http.HandleFunc("/register", register)
	http.HandleFunc("/verify", verifyEmail)
	http.HandleFunc("/sendResetEmail", sendPasswordResetEmail)
	http.HandleFunc("/reset", resetPassword)
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/createPostPicture", createPostPicture)
	http.HandleFunc("/updatePostPicture", updatePostPicture)
	http.HandleFunc("/deletePostPictures", deletePostPictures)
	http.HandleFunc("/getPostPicture", getPostPicture)
	http.ListenAndServe(port, nil)
	logger.Info("Starting the application at " + port + " 🚀")
}

func getAuthToken(request *http.Request) string {
	authToken := request.Header.Get("Authorization")
	splitToken := strings.Split(authToken, "Bearer ")
	if splitToken != nil && len(splitToken) > 1 {
		authToken = splitToken[1]
	}
	return authToken
}

func manageCors(w *http.ResponseWriter, r *http.Request) bool {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "*")
	if (*r).Method == "OPTIONS" {
		(*w).Header().Set("Access-Control-Max-Age", "86400")
		(*w).WriteHeader(http.StatusOK)
		return false
	}
	return true
}

func validType(thetype string) bool {
	for _, b := range validTypes {
		if b == thetype {
			return true
		}
	}
	return false
}
