package main

import (
	"cloud.google.com/go/storage"
	"errors"
	"github.com/graphql-go/graphql"
	// medium "github.com/medium/medium-sdk-go"
	// "gopkg.in/russross/blackfriday.v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	// "github.com/Depado/bfchroma"
)

func rootMutation() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"addPost": &graphql.Field{
				Type:        PostType,
				Description: "Create a Post Post",
				Args: graphql.FieldConfigArgument{
					"type": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"title": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"content": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"author": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"heroimage": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"images": &graphql.ArgumentConfig{
						Type: graphql.NewList(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					_, err := validateAdmin(params.Context.Value(tokenKey).(string))
					if err != nil {
						return nil, err
					}
					if params.Args["title"] == nil || params.Args["content"] == nil ||
						params.Args["author"] == nil || params.Args["type"] == nil ||
						params.Args["heroimage"] == nil || params.Args["images"] == nil {
						return nil, errors.New("title or content or author or type or heroimage or images not provided")
					}
					title, ok := params.Args["title"].(string)
					if !ok {
						return nil, errors.New("problem casting title to string")
					}
					author, ok := params.Args["author"].(string)
					if !ok {
						return nil, errors.New("problem casting author to string")
					}
					content, ok := params.Args["content"].(string)
					if !ok {
						return nil, errors.New("problem casting content to string")
					}
					thetype, ok := params.Args["type"].(string)
					if !ok {
						return nil, errors.New("problem casting type to string")
					}
					heroimage, ok := params.Args["heroimage"].(string)
					if !ok {
						return nil, errors.New("problem casting heroimage to string")
					}
					images, ok := params.Args["images"].([]interface{})
					if !ok {
						return nil, errors.New("problem casting images to string array")
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
					postData := bson.M{
						"title":     title,
						"content":   content,
						"author":    author,
						"views":     0,
						"heroimage": heroimage,
						"images":    images,
					}
					res, err := mongoCollection.InsertOne(ctxMongo, postData)
					if err != nil {
						return nil, err
					}
					id := res.InsertedID.(primitive.ObjectID)
					idstring := id.Hex()
					timestamp := objectidtimestamp(id)
					postData["date"] = timestamp.Unix()
					_, err = elasticClient.Index().
						Index(postElasticIndex).
						Type(postElasticType).
						Id(idstring).
						BodyJson(postData).
						Do(ctxElastic)
					if err != nil {
						return nil, err
					}
					postData["date"] = timestamp.Format(dateFormat)
					postData["id"] = idstring
					/*
						mediumContentHTML := string(blackfriday.Run([]byte(content), blackfriday.WithRenderer(bfchroma.NewRenderer())))
						_, err = mediumClient.CreatePost(medium.CreatePostOptions{
							UserID:        mediumUser.ID,
							Title:         title,
							Content:       mediumContentHTML,
							ContentFormat: medium.ContentFormatHTML,
							PublishStatus: medium.PublishStatusDraft,
						})
						if err != nil {
							return nil, err
						}
					*/
					return postData, nil
				},
			},
			"updatePost": &graphql.Field{
				Type:        PostType,
				Description: "Update a Post",
				Args: graphql.FieldConfigArgument{
					"type": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"id": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"title": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"content": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"author": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"heroimage": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"images": &graphql.ArgumentConfig{
						Type: graphql.NewList(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					_, err := validateAdmin(params.Context.Value(tokenKey).(string))
					if err != nil {
						return nil, err
					}
					if params.Args["id"] == nil || params.Args["type"] == nil {
						return nil, errors.New("post id or type not provided")
					}
					idstring, ok := params.Args["id"].(string)
					if !ok {
						return nil, errors.New("cannot cast id to string")
					}
					id, err := primitive.ObjectIDFromHex(idstring)
					if err != nil {
						return nil, err
					}
					updateData := bson.M{}
					if params.Args["title"] != nil {
						title, ok := params.Args["title"].(string)
						if !ok {
							return nil, errors.New("problem casting title to string")
						}
						updateData["title"] = title
					}
					if params.Args["author"] != nil {
						author, ok := params.Args["author"].(string)
						if !ok {
							return nil, errors.New("problem casting author to string")
						}
						updateData["author"] = author
					}
					if params.Args["content"] != nil {
						content, ok := params.Args["content"].(string)
						if !ok {
							return nil, errors.New("problem casting content to string")
						}
						updateData["content"] = content
					}
					if params.Args["heroimage"] != nil {
						heroimage, ok := params.Args["heroimage"].(string)
						if !ok {
							return nil, errors.New("problem casting heroimage to string")
						}
						updateData["heroimage"] = heroimage
					}
					if params.Args["images"] != nil {
						images, ok := params.Args["images"].([]interface{})
						if !ok {
							return nil, errors.New("problem casting images to string array")
						}
						updateData["images"] = images
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
					_, err = elasticClient.Update().
						Index(postElasticIndex).
						Type(postElasticType).
						Id(idstring).
						Doc(updateData).
						Do(ctxElastic)
					if err != nil {
						return nil, err
					}
					_, err = mongoCollection.UpdateOne(ctxMongo, bson.M{
						"_id": id,
					}, bson.M{
						"$set": updateData,
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
						id := postData["_id"].(primitive.ObjectID)
						postData["date"] = objectidtimestamp(id).Format(dateFormat)
						postData["id"] = id.Hex()
						delete(postData, "_id")
						foundstuff = true
						break
					}
					if !foundstuff {
						return nil, errors.New("post not found with given id")
					}
					if err != nil {
						return nil, err
					}
					return postData, nil
				},
			},
			"deletePost": &graphql.Field{
				Type:        PostType,
				Description: "Delete a Post Post",
				Args: graphql.FieldConfigArgument{
					"type": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"id": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					_, err := validateAdmin(params.Context.Value(tokenKey).(string))
					if err != nil {
						return nil, err
					}
					if params.Args["id"] == nil || params.Args["type"] == nil {
						return nil, errors.New("post id or type not provided")
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
					idstring, ok := params.Args["id"].(string)
					if !ok {
						return nil, errors.New("cannot cast id to string")
					}
					id, err := primitive.ObjectIDFromHex(idstring)
					if err != nil {
						return nil, err
					}
					_, err = elasticClient.Delete().
						Index(postElasticIndex).
						Type(postElasticType).
						Id(idstring).
						Do(ctxElastic)
					if err != nil {
						return nil, err
					}
					cursor, err := mongoCollection.Find(ctxMongo, bson.M{
						"_id": id,
					})
					if err != nil {
						return nil, err
					}
					defer cursor.Close(ctxMongo)
					var postData map[string]interface{}
					idstr := id.Hex()
					var foundstuff = false
					for cursor.Next(ctxMongo) {
						postPrimitive := &bson.D{}
						err = cursor.Decode(postPrimitive)
						if err != nil {
							return nil, err
						}
						postData = postPrimitive.Map()
						id := postData["_id"].(primitive.ObjectID)
						postData["date"] = objectidtimestamp(id).Format(dateFormat)
						postData["id"] = idstr
						delete(postData, "_id")
						foundstuff = true
						break
					}
					if !foundstuff {
						return nil, errors.New("post not found with given id")
					}
					_, err = mongoCollection.DeleteOne(ctxMongo, bson.M{
						"_id": id,
					})
					if err != nil {
						return nil, err
					}
					heroimageid, ok := postData["heroimage"].(string)
					if !ok {
						return nil, errors.New("cannot convert heroimage to string")
					}
					if len(heroimageid) > 0 {
						var herofileobj *storage.ObjectHandle
						if thetype == "blog" {
							herofileobj = imageBucket.Object(blogImageIndex + "/" + heroimageid)
						} else {
							herofileobj = imageBucket.Object(projectImageIndex + "/" + heroimageid)
						}
						if err := herofileobj.Delete(ctxStorage); err != nil {
							return nil, err
						}
					}
					primativeimageids, ok := postData["images"].(primitive.A)
					if !ok {
						return nil, errors.New("cannot convert imageids to primitive")
					}
					for _, primativeimageid := range primativeimageids {
						imageid := primativeimageid.(string)
						var fileobj *storage.ObjectHandle
						if thetype == "blog" {
							fileobj = imageBucket.Object(blogImageIndex + "/" + imageid)
						} else {
							fileobj = imageBucket.Object(projectImageIndex + "/" + imageid)
						}
						if err := fileobj.Delete(ctxStorage); err != nil {
							return nil, err
						}
					}
					if err != nil {
						return nil, err
					}
					return postData, nil
				},
			},
			"deleteUser": &graphql.Field{
				Type:        AccountType,
				Description: "Delete a User",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					_, err := validateAdmin(params.Context.Value(tokenKey).(string))
					if err != nil {
						return nil, err
					}
					if params.Args["id"] == nil {
						return nil, errors.New("user id not provided")
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
						userPrimitive := &bson.D{}
						err = cursor.Decode(userPrimitive)
						if err != nil {
							return nil, err
						}
						userData = userPrimitive.Map()
						id := userData["_id"].(primitive.ObjectID)
						userData["date"] = objectidtimestamp(id).Format(dateFormat)
						userData["id"] = id.Hex()
						delete(userData, "_id")
						foundstuff = true
						break
					}
					if !foundstuff {
						return nil, errors.New("user not found with given id")
					}
					_, err = userCollection.DeleteOne(ctxMongo, bson.M{
						"_id": id,
					})
					if err != nil {
						return nil, err
					}
					return userData, nil
				},
			},
		},
	})
}
