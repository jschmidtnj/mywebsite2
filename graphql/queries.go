package main

import (
	"github.com/graphql-go/graphql"
)

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
