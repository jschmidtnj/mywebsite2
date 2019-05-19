package main

import (
	"github.com/graphql-go/graphql"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
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
					var user User
					err = UserCollection.FindOne(CTX, bson.M{"email": account.(User).Email}).Decode(&user)
					if (err != nil) {
						return &User{}, nil
					}
					return user, nil
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
					var blogs []Blog
					findOptions := options.FindOptions{}
					findOptions.Sort = bson.D{
						{"_id", -1},
					}
					perpage, success := params.Args["perpage"].(int64)
					if (!success) {
						return &[]Blog{}, nil
					}
					findOptions.Limit = &perpage
					page, success := params.Args["page"].(int64)
					if (!success) {
						return &[]Blog{}, nil
					}
					findOptions.Skip = &page
					cursor, err := BlogCollection.Find(CTX, nil, &findOptions)
					if (err != nil) {
						return &[]Blog{}, nil
					}
					for cursor.Next(CTX) {
						var blog Blog
						err = cursor.Decode(&blog)
						if (err != nil) {
							return &[]Blog{}, nil
						}
						blogs = append(blogs, blog)
					}
					return blogs, nil
				},
			},
		},
	})
}
