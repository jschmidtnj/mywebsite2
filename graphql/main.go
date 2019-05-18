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
	jwt "github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	
)

var JwtSecret []byte

var Client *mongo.Client

func CreateToken(response http.ResponseWriter, request *http.Request) {
	var user User
	_ = json.NewDecoder(request.Body).Decode(&user)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"password": user.Password,
	})
	tokenString, error := token.SignedString(JwtSecret)
	if error != nil {
		fmt.Println(error)
	}
	response.Header().Set("content-type", "application/json")
	response.Write([]byte(`{ "token": "` + tokenString + `" }`))
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	JwtSecret = []byte(os.Getenv("SECRET"))
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	// Client mongodb client
	Client, _ = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	port := ":" + os.Getenv("PORT")
	fmt.Println("Starting the application at " + port)
	schema, _ := graphql.NewSchema(graphql.SchemaConfig{
		Query: RootQuery(),
	})
	http.HandleFunc("/graphql", func(response http.ResponseWriter, request *http.Request) {
		result := graphql.Do(graphql.Params{
			Schema:        schema,
			RequestString: request.URL.Query().Get("query"),
			Context:       context.WithValue(context.Background(), "token", request.URL.Query().Get("token")),
		})
		json.NewEncoder(response).Encode(result)
	})
	http.HandleFunc("/login", CreateToken)
	http.ListenAndServe(port, nil)
}
