package main

import (
	"github.com/graphql-go/graphql"
)

var accountsMock []User = []User{
	User{
		Id:       "1",
		Username: "nraboy",
		Password: "1234",
	},
	User{
		Id:       "2",
		Username: "mraboy",
		Password: "5678",
	},
}

var blogsMock []Blog = []Blog{
	Blog{
		Id:        "1",
		Author:    "nraboy",
		Title:     "Sample Article",
		Content:   "This is a sample article written by Nic Raboy",
		Pageviews: 1000,
	},
}

func RootQuery() (*graphql.Object) {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"account": &graphql.Field{
				Type: AccountType,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					account, err := ValidateJWT(params.Context.Value("token").(string))
					if err != nil {
						return nil, err
					}
					for _, accountMock := range accountsMock {
						if accountMock.Username == account.(User).Username {
							return accountMock, nil
						}
					}
					return &User{}, nil
				},
			},
			"blogs": &graphql.Field{
				Type: graphql.NewList(BlogType),
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return blogsMock, nil
				},
			},
		},
	})
}
