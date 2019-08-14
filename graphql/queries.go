package main

import (
	"encoding/json"
	"errors"

	"github.com/go-redis/redis"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/olivere/elastic/v7"
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
			"shortlinks": &graphql.Field{
				Type:        graphql.NewList(ShortLinkType),
				Description: "Get user short links",
				Args: graphql.FieldConfigArgument{
					"linkids": &graphql.ArgumentConfig{
						Type: graphql.NewList(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					_, err := validateLoggedIn(params.Context.Value(tokenKey).(string))
					if err != nil {
						return nil, err
					}
					if params.Args["linkids"] == nil {
						return nil, errors.New("no linkids argument found")
					}
					shortlinkidsInterface, ok := params.Args["linkids"].([]interface{})
					if !ok {
						return nil, errors.New("unable to cast shortlinks to array")
					}
					shortlinkids, err := interfaceListToStringList(shortlinkidsInterface)
					if err != nil {
						return nil, err
					}
					shortlinks, err := getShortLinks(shortlinkids)
					if err != nil {
						return nil, err
					}
					shortlinkObjects := make([]bson.M, len(shortlinks))
					for i, shortlink := range shortlinks {
						shortlinkObjects[i] = bson.M{
							"id":   shortlinkids[i],
							"link": shortlink,
						}
					}
					return shortlinkObjects, nil
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
					"categories": &graphql.ArgumentConfig{
						Type: graphql.NewList(graphql.String),
					},
					"tags": &graphql.ArgumentConfig{
						Type: graphql.NewList(graphql.String),
					},
					"cache": &graphql.ArgumentConfig{
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
					if params.Args["tags"] == nil {
						return nil, errors.New("no tags argument found")
					}
					tags, ok := params.Args["tags"].([]interface{})
					if !ok {
						return nil, errors.New("tags could not be cast to string array")
					}
					if params.Args["categories"] == nil {
						return nil, errors.New("no categories argument found")
					}
					categories, ok := params.Args["categories"].([]interface{})
					if !ok {
						return nil, errors.New("categories could not be cast to string array")
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
					if params.Args["cache"] == nil {
						return nil, errors.New("no cache argument found")
					}
					cache, ok := params.Args["cache"].(bool)
					if !ok {
						return nil, errors.New("cache could not be cast to bool")
					}
					fieldarray := params.Info.FieldASTs
					fieldselections := fieldarray[0].SelectionSet.Selections
					fields := make([]string, len(fieldselections))
					for i, field := range fieldselections {
						fieldast, ok := field.(*ast.Field)
						if !ok {
							return nil, errors.New("field cannot be converted to *ast.FIeld")
						}
						fields[i] = fieldast.Name.Value
					}
					params.Args["cache"] = true
					pathMap := map[string]interface{}{
						"path":   "posts",
						"args":   params.Args,
						"fields": fields,
					}
					cachepathBytes, err := json.Marshal(pathMap)
					if err != nil {
						return nil, err
					}
					cachepath := string(cachepathBytes)
					if cache {
						cachedresStr, err := redisClient.Get(cachepath).Result()
						if err != nil {
							if err != redis.Nil {
								return nil, err
							}
						} else {
							if len(cachedresStr) > 0 {
								var cachedres []map[string]interface{}
								err = json.Unmarshal([]byte(cachedresStr), &cachedres)
								if err != nil {
									return nil, err
								}
								return cachedres, nil
							}
						}
					}
					var postElasticIndex string
					if thetype == "blog" {
						postElasticIndex = blogElasticIndex
					} else {
						postElasticIndex = projectElasticIndex
					}
					var posts []map[string]interface{}
					if len(fields) > 0 {
						sourceContext := elastic.NewFetchSourceContext(true).Include(fields...)
						var numtags = len(tags)
						mustQueries := make([]elastic.Query, numtags+len(categories))
						for i, tag := range tags {
							mustQueries[i] = elastic.NewTermQuery("tags", tag)
						}
						for i, category := range categories {
							mustQueries[i+numtags] = elastic.NewTermQuery("categories", category)
						}
						query := elastic.NewBoolQuery().Must(mustQueries...)
						if len(searchterm) > 0 {
							mainquery := elastic.NewMultiMatchQuery(searchterm, postSearchFields...)
							query = query.Filter(mainquery)
						}
						searchResult, err := elasticClient.Search().
							Index(postElasticIndex).
							Query(query).
							Sort(sort, ascending).
							From(page * perpage).Size(perpage).
							Pretty(false).
							FetchSourceContext(sourceContext).
							Do(ctxElastic)
						if err != nil {
							return nil, err
						}
						posts = make([]map[string]interface{}, len(searchResult.Hits.Hits))
						for i, hit := range searchResult.Hits.Hits {
							if hit.Source == nil {
								return nil, errors.New("no hit source found")
							}
							var postData map[string]interface{}
							err := json.Unmarshal(hit.Source, &postData)
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
							posts[i] = postData
						}
					}
					postsBytes, err := json.Marshal(posts)
					if err != nil {
						return nil, err
					}
					err = redisClient.Set(cachepath, string(postsBytes), cacheTime).Err()
					if err != nil {
						return nil, err
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
					"cache": &graphql.ArgumentConfig{
						Type: graphql.Boolean,
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
					if params.Args["cache"] == nil {
						return nil, errors.New("no cache argument found")
					}
					cache, ok := params.Args["cache"].(bool)
					if !ok {
						return nil, errors.New("cache could not be cast to bool")
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
					fieldarray := params.Info.FieldASTs
					fieldselections := fieldarray[0].SelectionSet.Selections
					fields := make([]string, len(fieldselections))
					for i, field := range fieldselections {
						fieldast, ok := field.(*ast.Field)
						if !ok {
							return nil, errors.New("field cannot be converted to *ast.FIeld")
						}
						fields[i] = fieldast.Name.Value
					}
					params.Args["cache"] = true
					pathMap := map[string]interface{}{
						"path":   "post",
						"args":   params.Args,
						"fields": fields,
					}
					cachepathBytes, err := json.Marshal(pathMap)
					if err != nil {
						return nil, err
					}
					cachepath := string(cachepathBytes)
					if cache {
						cachedresStr, err := redisClient.Get(cachepath).Result()
						if err != nil {
							if err != redis.Nil {
								return nil, err
							}
						} else {
							if len(cachedresStr) > 0 {
								var cachedres map[string]interface{}
								err = json.Unmarshal([]byte(cachedresStr), &cachedres)
								if err != nil {
									return nil, err
								}
								return cachedres, nil
							}
						}
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
					if err != nil {
						return nil, err
					}
					postBytes, err := json.Marshal(postData)
					if err != nil {
						return nil, err
					}
					err = redisClient.Set(cachepath, string(postBytes), cacheTime).Err()
					if err != nil {
						return nil, err
					}
					return postData, nil
				},
			},
		},
	})
}
