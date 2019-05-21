package main

import (
	"github.com/graphql-go/graphql"
)

var BlogType *graphql.Object = graphql.NewObject(graphql.ObjectConfig{
	Name: "Blog",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
		},
		"title": &graphql.Field{
			Type: graphql.String,
		},
		"content": &graphql.Field{
			Type: graphql.String,
		},
		"author": &graphql.Field{
			Type: graphql.String,
		},
		"pageviews": &graphql.Field{
			Type: graphql.Int,
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				return 0, nil
			},
		},
	},
})
