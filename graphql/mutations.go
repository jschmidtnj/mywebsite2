package main

import (
	"github.com/graphql-go/graphql"
	"go.mongodb.org/mongo-driver/bson"
	"errors"
)

func handleErrorBlog(err error) (interface{}, error) {
	return &Blog{}, nil
}

func RootMutation() (*graphql.Object) {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"addBlog": &graphql.Field{
				Type: BlogType,
				Description: "Create a Blog Post",
				Args: graphql.FieldConfigArgument{
					"title": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"content": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"author": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					title, success := params.Args["title"].(string)
					if (!success) {
						return handleErrorBlog(errors.New("no title found"))
					}
					content, success := params.Args["content"].(string)
					if (!success) {
						return handleErrorBlog(errors.New("no content found"))
					}
					author, success := params.Args["author"].(string)
					if (!success) {
						return handleErrorBlog(errors.New("no author found"))
					}
					insertRes, err := BlogCollection.InsertOne(CTX, bson.M{
						"title": title,
						"content": content,
						"author": author,
					})
					if (err != nil) {
						return handleErrorBlog(err)
					}
					blog := Blog{
						Title: title,
						Content: content,
						Author: author,
						Id: insertRes.InsertedID.(string),
					}
					return &blog, nil
				},
			},
		},
	})
}