package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"log"
	"time"
	"github.com/joho/godotenv"
	"os"
	"github.com/graphql-go/graphql"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var JwtSecret []byte

var Client *mongo.Client

var CTX context.Context

var Database string = "testing"

var UserCollection *mongo.Collection

var BlogCollection *mongo.Collection

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	JwtSecret = []byte(os.Getenv("SECRET"))
	CTX, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	cancel()
	mongouri := os.Getenv("MONGOURI")
	Client, err = mongo.Connect(CTX, options.Client().ApplyURI(mongouri))
	if err != nil {
		log.Fatal(err.Error())
	}
	UserCollection = Client.Database(Database).Collection("users")
	BlogCollection = Client.Database(Database).Collection("blogs")
	port := ":" + os.Getenv("PORT")
	fmt.Println("Starting the application at " + port)
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: RootQuery(),
		Mutation: RootMutation(),
	})
	if err != nil {
		log.Fatal(err.Error())
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
	http.ListenAndServe(port, nil)
}
