package main

import (
	"github.com/graphql-go/graphql"
	"go.mongodb.org/mongo-driver/bson"
	"errors"
)

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
						return nil, errors.New("no title found")
					}
					content, success := params.Args["content"].(string)
					if (!success) {
						return nil, errors.New("no content found")
					}
					author, success := params.Args["author"].(string)
					if (!success) {
						return nil, errors.New("no author found")
					}
					blogdata := bson.M{
						"title": title,
						"content": content,
						"author": author,
					}
					_, err := BlogCollection.InsertOne(CTX, blogdata)
					if (err != nil) {
						return nil, err
					}
					return blogdata, nil
				},
			},
		},
	})
}