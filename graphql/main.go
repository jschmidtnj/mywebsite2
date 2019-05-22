package main

import (
	"context"
	"encoding/json"
	"net/http"
	"go.uber.org/zap"
	"time"
	"github.com/joho/godotenv"
	"os"
	"github.com/graphql-go/graphql"
	"go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/mongo/options"
  "strconv"
)

var JwtSecret []byte

var TokenExpiration int

var SendgridApiKey string

var WebsiteUrl string

var Client *mongo.Client

var CTX context.Context

var Database string = "testing"

var UserCollection *mongo.Collection

var BlogCollection *mongo.Collection

var Logger *zap.Logger

func Hello(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	response.Write([]byte(`{"message":"Hello!"}`))
}

func main() {
	loggerconfig := []byte(`{
		"level": "debug",
		"encoding": "json",
		"outputPaths": ["stdout", "./logs"],
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
  Logger, err = zapconfig.Build()
  if err != nil {
    panic(err)
  }
  defer Logger.Sync()
  Logger.Info("logger created")
	err = godotenv.Load()
	if err != nil {
		Logger.Fatal("Error loading .env file")
	}
  JwtSecret = []byte(os.Getenv("SECRET"))
  TokenExpiration, err = strconv.Atoi(os.Getenv("TOKENEXPIRATION"))
  if (err != nil) {
    Logger.Fatal(err.Error())
  }
  SendgridApiKey = os.Getenv("SENDGRIDAPIKEY")
  WebsiteUrl = os.Getenv("WEBSITEURL")
	CTX, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	cancel()
	mongouri := os.Getenv("MONGOURI")
	Client, err = mongo.Connect(CTX, options.Client().ApplyURI(mongouri))
	if err != nil {
		Logger.Fatal(err.Error())
	}
	UserCollection = Client.Database(Database).Collection("users")
	BlogCollection = Client.Database(Database).Collection("blogs")
	port := ":" + os.Getenv("PORT")
	Logger.Info("Starting the application at " + port)
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: RootQuery(),
		Mutation: RootMutation(),
	})
	if err != nil {
		Logger.Fatal(err.Error())
	}
	http.HandleFunc("/graphql", func(response http.ResponseWriter, request *http.Request) {
		result := graphql.Do(graphql.Params{
			Schema:        schema,
			RequestString: request.URL.Query().Get("query"),
			Context:       context.WithValue(context.Background(), "token", request.URL.Query().Get("token")),
		})
		json.NewEncoder(response).Encode(result)
	})
	http.HandleFunc("/login", Login)
  http.HandleFunc("/register", Register)
  http.HandleFunc("/hello", Hello)
	http.ListenAndServe(port, nil)
}
