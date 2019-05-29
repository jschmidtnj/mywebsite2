package main

import (
	"github.com/graphql-go/graphql"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/olivere/elastic/v7"
	"errors"
	"encoding/json"
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
					var foundstuff = false
					for cursor.Next(CTX) {
						userDataPrimitive := &bson.D{}
						err = cursor.Decode(userDataPrimitive)
						if (err != nil) {
							return nil, err
						}
						userData = userDataPrimitive.Map()
						id := userData["_id"].(primitive.ObjectID)
						userData["date"] = objectidtimestamp(id).Format(DateFormat)
						userData["id"] = id.Hex()
						delete(userData, "_id")
						foundstuff = true
						break
					}
					if (!foundstuff) {
						return nil, errors.New("account data not found")
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
					"searchterm": &graphql.ArgumentConfig{
            Type: graphql.String,
					},
        },
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					if (params.Args["perpage"] == nil) {
						return nil, errors.New("no perpage argument found")
					}
					perpage64, ok := params.Args["perpage"].(int64)
					if (!ok) {
						return nil, errors.New("perpage could not be cast to int")
					}
					perpage := int(perpage64)
					if (params.Args["page"] == nil) {
						return nil, errors.New("no page argument found")
					}
					page64, ok := params.Args["page"].(int64)
					if (!ok) {
						return nil, errors.New("page could not be cast to int")
					}
					page := int(page64)
					searchterm, ok := params.Args["searchterm"].(string)
					if (!ok) {
						return nil, errors.New("searchterm could not be cast to string")
					}
					queryString := elastic.NewQueryStringQuery(searchterm)
					searchResult, err := Elastic.Search().
						Index(BlogElasticIndex).
						Query(queryString).
						Sort("date", true). // ascending = true
						From(page).Size(perpage).
						Pretty(false).
						Do(CTX)
					if (err != nil) {
						return nil, err
					}
					var blogs []map[string]interface{}
					for _, hit := range searchResult.Hits.Hits {
						// Deserialize hit.Source into a Tweet (could also be just a map[string]interface{}).
						var blogData map[string]interface{}
						err := json.Unmarshal(hit.Source, &blogData)
						if (err != nil) {
							return nil, err
						}
						id, err := primitive.ObjectIDFromHex(hit.Id)
						if (err != nil) {
							return nil, err
						}
						blogData["date"] = objectidtimestamp(id).Format(DateFormat)
						blogData["id"] = id.Hex()
						delete(blogData, "_id")
						blogs = append(blogs, blogData)
					}
					return blogs, nil
				},
			},
			"blog": &graphql.Field{
				Type: BlogType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					if (params.Args["id"] == nil) {
						return nil, errors.New("no id argument found")
					}
					id, err := primitive.ObjectIDFromHex(params.Args["id"].(string))
					if (err != nil) {
						return nil, err
					}
					_, err = BlogCollection.UpdateOne(CTX, bson.M{
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
					var foundstuff = false
					for cursor.Next(CTX) {
						blogPrimitive := &bson.D{}
						err = cursor.Decode(blogPrimitive)
						if (err != nil) {
							return nil, err
						}
						blogData = blogPrimitive.Map()
						id := blogData["_id"].(primitive.ObjectID)
						blogData["date"] = objectidtimestamp(id).Format(DateFormat)
						blogData["views"] = int(blogData["views"].(int32))
						blogData["id"] = id.Hex()
						delete(blogData, "_id")
						foundstuff = true
						break
					}
					if (!foundstuff) {
						return nil, errors.New("blog not found with given id")
					}
					return blogData, nil
				},
			},
		},
	})
}
