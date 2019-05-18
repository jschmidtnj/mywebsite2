package main

import (
	"github.com/graphql-go/graphql"
)

type Blog struct {
	Id        string `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Author    string `json:"author"`
	Pageviews int32  `json:"pageviews"`
}

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
				_, err := ValidateJWT(params.Context.Value("token").(string))
				if err != nil {
					return nil, err
				}
				return params.Source.(Blog).Pageviews, nil
			},
		},
	},
})
