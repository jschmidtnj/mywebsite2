package main

import (
	"github.com/graphql-go/graphql"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
					_, err := ValidateAdmin(params.Context.Value("token").(string))
					if (err != nil) {
						return nil, err
					}
					if (params.Args["title"] == nil || params.Args["content"] == nil || params.Args["author"] == nil) {
						return nil, errors.New("title or content or author not provided")
					}
					blogdata := bson.M{
						"title": params.Args["title"].(string),
						"content": params.Args["content"].(string),
						"author": params.Args["author"].(string),
						"views": 0,
					}
					res, err := BlogCollection.InsertOne(CTX, blogdata)
					if (err != nil) {
						return nil, err
					}
					blogdata["id"] = res.InsertedID.(primitive.ObjectID).Hex()
					return blogdata, nil
				},
			},
			"deleteUser": &graphql.Field{
				Type: BlogType,
				Description: "Create a Blog Post",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					_, err := ValidateAdmin(params.Context.Value("token").(string))
					if (err != nil) {
						return nil, err
					}
					if (params.Args["id"] == nil) {
						return nil, errors.New("user id not provided")
					}
					res, err := UserCollection.DeleteOne(CTX, bson.M{
						"_id": params.Args["id"].(string),
					})
					if (err != nil) {
						return nil, err
					}
					return res, nil
				},
			},
		},
	})
}