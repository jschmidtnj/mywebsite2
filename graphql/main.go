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
	"github.com/olivere/elastic"
	"strconv"
	"strings"
)

var JwtSecret []byte

var TokenExpiration int

var SendgridApiKey string

var WebsiteUrl string

var Client *mongo.Client

var CTX context.Context

var Database = "testing"

var UserCollection *mongo.Collection

var BlogCollection *mongo.Collection

var Elastic *elastic.Client

var CTXElastic context.Context

var BlogElasticIndex = "blogs"

var Logger *zap.Logger

func Hello(response http.ResponseWriter, request *http.Request) {
	if !ManageCors(response, request) {
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
  Logger, err = zapconfig.Build()
  if (err != nil) {
    panic(err)
  }
  defer Logger.Sync()
  Logger.Info("logger created")
	err = godotenv.Load()
	if (err != nil) {
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
	if (err != nil) {
		Logger.Fatal(err.Error())
	}
	UserCollection = Client.Database(Database).Collection("users")
	BlogCollection = Client.Database(Database).Collection("blogs")
	elasticuri := os.Getenv("ELASTICURI")
	Elastic, err = elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(elasticuri))
	if (err != nil) {
		Logger.Fatal(err.Error())
	}
	CTXElastic = context.Background()
	port := ":" + os.Getenv("PORT")
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: RootQuery(),
		Mutation: RootMutation(),
	})
	if (err != nil) {
		Logger.Fatal(err.Error())
	}
	http.HandleFunc("/graphql", func(response http.ResponseWriter, request *http.Request) {
		if !ManageCors(response, request) {
			return
		}
		result := graphql.Do(graphql.Params{
			Schema:        schema,
			RequestString: request.URL.Query().Get("query"),
			Context:       context.WithValue(context.Background(), "token", GetAuthToken(request)),
		})
		response.Header().Set("content-type", "application/json")
		json.NewEncoder(response).Encode(result)
	})
	http.HandleFunc("/countBlogs", CountBlogs)
	http.HandleFunc("/sendTestEmail", SendTestEmail)
	http.HandleFunc("/loginEmailPassword", LoginEmailPassword)
	http.HandleFunc("/logoutEmailPassword", LogoutEmailPassword)
	http.HandleFunc("/register", Register)
	http.HandleFunc("/verify", VerifyEmail)
	http.HandleFunc("/sendResetEmail", SendPasswordResetEmail)
	http.HandleFunc("/reset", ResetPassword)
  http.HandleFunc("/hello", Hello)
	http.ListenAndServe(port, nil)
	Logger.Info("Starting the application at " + port + " 🚀")
}

func GetAuthToken(request *http.Request) string {
	authToken := request.Header.Get("Authorization")
	splitToken := strings.Split(authToken, "Bearer ")
	if (splitToken != nil && len(splitToken) > 1) {
		authToken = splitToken[1]
	}
	return authToken
}

func ManageCors(w http.ResponseWriter, r *http.Request) bool {
  w.Header().Set("Access-Control-Allow-Origin", "*")
  w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
  w.Header().Set("Access-Control-Allow-Headers", "*")
  if r.Method == "OPTIONS" {
    w.Header().Set("Access-Control-Max-Age", "86400")
		w.WriteHeader(http.StatusOK)
		return false
	}
	return true
}
