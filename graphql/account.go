package main

import (
	"github.com/graphql-go/graphql"
)

type User struct {
	Id       string `json:"id"`
	Email string `json:"email"`
	Password string `json:"password"`
}

var AccountType *graphql.Object = graphql.NewObject(graphql.ObjectConfig{
	Name: "Account",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
		},
		"email": &graphql.Field{
			Type: graphql.String,
		},
		"password": &graphql.Field{
			Type: graphql.String,
		},
	},
})
