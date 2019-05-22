package main

import (
	"github.com/graphql-go/graphql"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"errors"
)

func RootQuery() (*graphql.Object) {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"account": &graphql.Field{
				Type: AccountType,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					account, err := ValidateLoggedIn(params.Context.Value("token").(string))
					if err != nil {
						return nil, err
					}
					accountdata := account.(map[string]interface{})
					if (accountdata["email"] == nil) {
						return nil, errors.New("email not found in token")
					}
					cursor, err := UserCollection.Find(CTX, bson.M{"email": accountdata["email"]})
					if (err != nil) {
						return nil, err
					}
					defer cursor.Close(CTX)
					userDataPrimitive := &bson.D{}
					err = cursor.Decode(userDataPrimitive)
					if (err != nil) {
						return nil, err
					}
					return userDataPrimitive, nil
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
					findOptions.Sort = bson.D{
						{"_id", -1},
					}
					perpage, success := params.Args["perpage"].(int64)
					if (!success) {
						return nil, errors.New("no perpage argument found")
					}
					findOptions.Limit = &perpage
					page, success := params.Args["page"].(int64)
					if (!success) {
						return nil, errors.New("no page argument found")
					}
					findOptions.Skip = &page
					cursor, err := BlogCollection.Find(CTX, nil, &findOptions)
					if (err != nil) {
						return nil, err
					}
					var blogs []*bson.D
					for cursor.Next(CTX) {
						blogPrimitive := &bson.D{}
						err = cursor.Decode(&blogPrimitive)
						if (err != nil) {
							return nil, err
						}
						blogs = append(blogs, blogPrimitive)
					}
					return blogs, nil
				},
			},
		},
	})
}
