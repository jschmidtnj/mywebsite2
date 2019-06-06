package main

import (
	"github.com/graphql-go/graphql"
)

// BlogType graphql blog type is a blog object
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
		"views": &graphql.Field{
			Type: graphql.Int,
		},
		"date": &graphql.Field{
			Type: graphql.String,
		},
		"heroimage": &graphql.Field{
			Type: graphql.String,
		},
		"images": &graphql.Field{
			Type: graphql.NewList(graphql.String),
		},
	},
})
