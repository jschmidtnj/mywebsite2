package main

import (
	"encoding/json"
	"errors"
	"github.com/graphql-go/graphql"
	"github.com/olivere/elastic"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func rootQuery() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"hello": &graphql.Field{
				Type:        graphql.String,
				Description: "Say Hi",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return "Hello World!", nil
				},
			},
			"account": &graphql.Field{
				Type:        AccountType,
				Description: "Get your account info",
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					accountdata, err := validateLoggedIn(params.Context.Value(tokenKey).(string))
					if err != nil {
						return nil, err
					}
					if accountdata["email"] == nil {
						return nil, errors.New("email not found in token")
					}
					cursor, err := userCollection.Find(ctxMongo, bson.M{"email": accountdata["email"].(string)})
					defer cursor.Close(ctxMongo)
					if err != nil {
						return nil, err
					}
					var userData map[string]interface{}
					var foundstuff = false
					for cursor.Next(ctxMongo) {
						userDataPrimitive := &bson.D{}
						err = cursor.Decode(userDataPrimitive)
						if err != nil {
							return nil, err
						}
						userData = userDataPrimitive.Map()
						id := userData["_id"].(primitive.ObjectID)
						userData["date"] = objectidtimestamp(id).Format(dateFormat)
						userData["id"] = id.Hex()
						delete(userData, "_id")
						foundstuff = true
						break
					}
					if !foundstuff {
						return nil, errors.New("account data not found")
					}
					return userData, nil
				},
			},
			"user": &graphql.Field{
				Type:        AccountType,
				Description: "Get a user by id as admin",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					accountdata, err := validateAdmin(params.Context.Value(tokenKey).(string))
					if err != nil {
						return nil, err
					}
					if accountdata["id"] == nil {
						return nil, errors.New("email not found in token")
					}
					idstring, ok := params.Args["id"].(string)
					if !ok {
						return nil, errors.New("cannot cast id to string")
					}
					id, err := primitive.ObjectIDFromHex(idstring)
					if err != nil {
						return nil, err
					}
					cursor, err := userCollection.Find(ctxMongo, bson.M{
						"_id": id,
					})
					defer cursor.Close(ctxMongo)
					if err != nil {
						return nil, err
					}
					var userData map[string]interface{}
					var foundstuff = false
					for cursor.Next(ctxMongo) {
						userDataPrimitive := &bson.D{}
						err = cursor.Decode(userDataPrimitive)
						if err != nil {
							return nil, err
						}
						userData = userDataPrimitive.Map()
						userData["date"] = objectidtimestamp(id).Format(dateFormat)
						userData["id"] = idstring
						delete(userData, "_id")
						foundstuff = true
						break
					}
					if !foundstuff {
						return nil, errors.New("account data not found")
					}
					return userData, nil
				},
			},
			"posts": &graphql.Field{
				Type:        graphql.NewList(PostType),
				Description: "Get list of posts",
				Args: graphql.FieldConfigArgument{
					"type": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"perpage": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
					"page": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
					"searchterm": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"sort": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"ascending": &graphql.ArgumentConfig{
						Type: graphql.Boolean,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					// see this: https://github.com/olivere/elastic/issues/483
					// for potential fix to source issue (tried gave null pointer error)
					if params.Args["perpage"] == nil {
						return nil, errors.New("no perpage argument found")
					}
					perpage, ok := params.Args["perpage"].(int)
					if !ok {
						return nil, errors.New("perpage could not be cast to int")
					}
					if params.Args["page"] == nil {
						return nil, errors.New("no page argument found")
					}
					page, ok := params.Args["page"].(int)
					if !ok {
						return nil, errors.New("page could not be cast to int")
					}
					var searchterm string
					if params.Args["searchterm"] != nil {
						searchterm, ok = params.Args["searchterm"].(string)
						if !ok {
							return nil, errors.New("searchterm could not be cast to string")
						}
					}
					if params.Args["sort"] == nil {
						return nil, errors.New("sort is undefined")
					}
					sort, ok := params.Args["sort"].(string)
					if !ok {
						return nil, errors.New("sort could not be cast to string")
					}
					if params.Args["ascending"] == nil {
						return nil, errors.New("ascending is undefined")
					}
					ascending, ok := params.Args["ascending"].(bool)
					if !ok {
						return nil, errors.New("ascending could not be cast to boolean")
					}
					if params.Args["type"] == nil {
						return nil, errors.New("type is undefined")
					}
					thetype, ok := params.Args["type"].(string)
					if !ok {
						return nil, errors.New("problem casting type to string")
					}
					if !validType(thetype) {
						return nil, errors.New("invalid type given")
					}
					var postElasticIndex string
					if thetype == "blog" {
						postElasticIndex = blogElasticIndex
					} else {
						postElasticIndex = projectElasticIndex
					}
					var searchResult *elastic.SearchResult
					var err error
					if len(searchterm) > 0 {
						queryString := elastic.NewQueryStringQuery(searchterm)
						searchResult, err = elasticClient.Search().
							Index(postElasticIndex).
							Query(queryString).
							Sort(sort, ascending).
							From(page * perpage).Size(perpage).
							Pretty(false).
							Do(ctxElastic)
					} else {
						searchResult, err = elasticClient.Search().
							Index(postElasticIndex).
							Query(nil).
							Sort(sort, ascending).
							From(page * perpage).Size(perpage).
							Pretty(false).
							Do(ctxElastic)
					}
					if err != nil {
						return nil, err
					}
					var posts []map[string]interface{}
					for _, hit := range searchResult.Hits.Hits {
						if *hit.Source == nil {
							return nil, errors.New("no hit source found")
						}
						var postData map[string]interface{}
						err := json.Unmarshal(*hit.Source, &postData)
						if err != nil {
							return nil, err
						}
						id, err := primitive.ObjectIDFromHex(hit.Id)
						if err != nil {
							return nil, err
						}
						postData["date"] = objectidtimestamp(id).Format(dateFormat)
						postData["id"] = id.Hex()
						delete(postData, "_id")
						posts = append(posts, postData)
					}
					return posts, nil
				},
			},
			"post": &graphql.Field{
				Type:        PostType,
				Description: "Get a Post Post",
				Args: graphql.FieldConfigArgument{
					"type": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"id": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					if params.Args["type"] == nil {
						return nil, errors.New("type is undefined")
					}
					thetype, ok := params.Args["type"].(string)
					if !ok {
						return nil, errors.New("problem casting type to string")
					}
					if !validType(thetype) {
						return nil, errors.New("invalid type given")
					}
					var mongoCollection *mongo.Collection
					var postElasticIndex string
					var postElasticType string
					if thetype == "blog" {
						mongoCollection = blogCollection
						postElasticIndex = blogElasticIndex
						postElasticType = blogElasticType
					} else {
						mongoCollection = projectCollection
						postElasticIndex = projectElasticIndex
						postElasticType = projectElasticType
					}
					if params.Args["id"] == nil {
						return nil, errors.New("no id argument found")
					}
					idstring, ok := params.Args["id"].(string)
					if !ok {
						return nil, errors.New("cannot cast id to string")
					}
					id, err := primitive.ObjectIDFromHex(idstring)
					if err != nil {
						return nil, err
					}
					_, err = mongoCollection.UpdateOne(ctxMongo, bson.M{
						"_id": id,
					}, bson.M{
						"$inc": bson.M{
							"views": 1,
						},
					})
					if err != nil {
						return nil, err
					}
					cursor, err := mongoCollection.Find(ctxMongo, bson.M{
						"_id": id,
					})
					defer cursor.Close(ctxMongo)
					if err != nil {
						return nil, err
					}
					var postData map[string]interface{}
					var foundstuff = false
					for cursor.Next(ctxMongo) {
						postPrimitive := &bson.D{}
						err = cursor.Decode(postPrimitive)
						if err != nil {
							return nil, err
						}
						postData = postPrimitive.Map()
						postData["date"] = objectidtimestamp(id).Format(dateFormat)
						postData["id"] = idstring
						delete(postData, "_id")
						foundstuff = true
						break
					}
					if !foundstuff {
						return nil, errors.New("post not found with given id")
					}
					_, err = elasticClient.Update().
						Index(postElasticIndex).
						Type(postElasticType).
						Id(idstring).
						Doc(bson.M{
							"views": int(postData["views"].(int32)),
						}).
						Do(ctxElastic)
					return postData, nil
				},
			},
		},
	})
}
