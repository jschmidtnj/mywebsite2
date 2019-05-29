package main

import (
	"github.com/graphql-go/graphql"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"errors"
)

func RootQuery() (*graphql.Object) {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"hello": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return "Hello World!", nil
				},
			},
			"account": &graphql.Field{
				Type: AccountType,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					accountdata, err := ValidateLoggedIn(params.Context.Value("token").(string))
					if err != nil {
						return nil, err
					}
					if (accountdata["email"] == nil) {
						return nil, errors.New("email not found in token")
					}
					cursor, err := UserCollection.Find(CTX, bson.M{"email": accountdata["email"].(string)})
					defer cursor.Close(CTX)
					if (err != nil) {
						return nil, err
					}
					var userData map[string]interface{}
					for cursor.Next(CTX) {
						userDataPrimitive := &bson.D{}
						err = cursor.Decode(userDataPrimitive)
						if (err != nil) {
							return nil, err
						}
						userData = userDataPrimitive.Map()
						userData["id"] = userData["_id"].(primitive.ObjectID).Hex()
						delete(userData, "_id")
					}
					return userData, nil
				},
			},
			"blogs": &graphql.Field{
				Type: graphql.NewList(BlogType),
				Args: graphql.FieldConfigArgument{
          "perpage": &graphql.ArgumentConfig{
            Type: graphql.Int,
					},
					"page": &graphql.ArgumentConfig{
            Type: graphql.Int,
					},
        },
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					findOptions := options.FindOptions{}
					findOptions.Sort = bson.M{"_id": -1}
					if (params.Args["perpage"] == nil) {
						return nil, errors.New("no perpage argument found")
					}
					var perpage int64 = params.Args["perpage"].(int64)
					findOptions.Limit = &perpage
					if (params.Args["page"] == nil) {
						return nil, errors.New("no page argument found")
					}
					var page int64 = params.Args["page"].(int64)
					findOptions.Skip = &page
					cursor, err := BlogCollection.Find(CTX, nil, &findOptions)
					if (err != nil) {
						return nil, err
					}
					var blogs []map[string]interface{}
					for cursor.Next(CTX) {
						blogPrimitive := &bson.D{}
						err = cursor.Decode(blogPrimitive)
						if (err != nil) {
							return nil, err
						}
						blogData := blogPrimitive.Map()
						blogData["id"] = blogData["_id"].(primitive.ObjectID).Hex()
						delete(blogData, "_id")
						blogs = append(blogs, blogData)
					}
					return blogs, nil
				},
			},
			"blog": &graphql.Field{
				Type: graphql.NewList(BlogType),
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					if (params.Args["id"] == nil) {
						return nil, errors.New("no id argument found")
					}
					var id string = params.Args["id"].(string)
					_, err := BlogCollection.UpdateOne(CTX, bson.M{
						"_id": id,
					}, bson.M{
						"$inc": bson.M{
							"views": 1,
						},
					})
					if (err != nil) {
						return nil, err
					}
					cursor, err := BlogCollection.Find(CTX, bson.M{
						"_id": id,
					})
					defer cursor.Close(CTX)
					if (err != nil) {
						return nil, err
					}
					var blogData map[string]interface{}
					for cursor.Next(CTX) {
						blogPrimitive := &bson.D{}
						err = cursor.Decode(blogPrimitive)
						if (err != nil) {
							return nil, err
						}
						blogData = blogPrimitive.Map()
						blogData["id"] = blogData["_id"].(primitive.ObjectID).Hex()
						delete(blogData, "_id")
						break
					}
					return blogData, nil
				},
			},
		},
	})
}
