package main

import (
	"github.com/graphql-go/graphql"
)

// AccountType account type object for user accounts graphql
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
		"emailverified": &graphql.Field{
			Type: graphql.Boolean,
		},
		"type": &graphql.Field{
			Type: graphql.String,
		},
	},
})
