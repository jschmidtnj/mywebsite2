package main

import (
	"errors"
	"github.com/graphql-go/graphql"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func rootMutation() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"addBlog": &graphql.Field{
				Type:        BlogType,
				Description: "Create a Blog Post",
				Args: graphql.FieldConfigArgument{
					"title": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"content": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"author": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					_, err := validateAdmin(params.Context.Value("token").(string))
					if err != nil {
						return nil, err
					}
					if params.Args["title"] == nil || params.Args["content"] == nil || params.Args["author"] == nil {
						return nil, errors.New("title or content or author not provided")
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
					blogData := bson.M{
						"title":     title,
						"content":   content,
						"author":    author,
						"views":     0,
						"heroimage": "",
						"images":    []string{},
					}
					res, err := blogCollection.InsertOne(ctxMongo, blogData)
					if err != nil {
						return nil, err
					}
					id := res.InsertedID.(primitive.ObjectID)
					idstring := id.Hex()
					timestamp := objectidtimestamp(id)
					blogData["date"] = timestamp.Unix()
					_, err = elasticClient.Index().
						Index(blogElasticIndex).
						Type("blog").
						Id(idstring).
						BodyJson(blogData).
						Do(ctxElastic)
					if err != nil {
						return nil, err
					}
					blogData["date"] = timestamp.Format(dateFormat)
					blogData["id"] = idstring
					return blogData, nil
				},
			},
			"updateBlog": &graphql.Field{
				Type:        BlogType,
				Description: "Update a Blog Post",
				Args: graphql.FieldConfigArgument{
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
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					_, err := validateAdmin(params.Context.Value("token").(string))
					if err != nil {
						return nil, err
					}
					if params.Args["id"] == nil {
						return nil, errors.New("blog id not provided")
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
					_, err = elasticClient.Update().
						Index(blogElasticIndex).
						Type("blog").
						Id(idstring).
						Doc(updateData).
						Do(ctxElastic)
					if err != nil {
						return nil, err
					}
					_, err = blogCollection.UpdateOne(ctxMongo, bson.M{
						"_id": id,
					}, bson.M{
						"$set": updateData,
					})
					if err != nil {
						return nil, err
					}
					cursor, err := blogCollection.Find(ctxMongo, bson.M{
						"_id": id,
					})
					defer cursor.Close(ctxMongo)
					if err != nil {
						return nil, err
					}
					var blogData map[string]interface{}
					var foundstuff = false
					for cursor.Next(ctxMongo) {
						blogPrimitive := &bson.D{}
						err = cursor.Decode(blogPrimitive)
						if err != nil {
							return nil, err
						}
						blogData = blogPrimitive.Map()
						id := blogData["_id"].(primitive.ObjectID)
						blogData["date"] = objectidtimestamp(id).Format(dateFormat)
						blogData["id"] = id.Hex()
						delete(blogData, "_id")
						foundstuff = true
						break
					}
					if !foundstuff {
						return nil, errors.New("blog not found with given id")
					}
					return blogData, nil
				},
			},
			"deleteBlog": &graphql.Field{
				Type:        BlogType,
				Description: "Delete a Blog Post",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					_, err := validateAdmin(params.Context.Value("token").(string))
					if err != nil {
						return nil, err
					}
					if params.Args["id"] == nil {
						return nil, errors.New("blog id not provided")
					}
					idstring, ok := params.Args["id"].(string)
					if !ok {
						return nil, errors.New("cannot cast id to string")
					}
					id, err := primitive.ObjectIDFromHex(idstring)
					if err != nil {
						return nil, err
					}
					cursor, err := blogCollection.Find(ctxMongo, bson.M{
						"_id": id,
					})
					if err != nil {
						return nil, err
					}
					defer cursor.Close(ctxMongo)
					var blogData map[string]interface{}
					idstr := id.Hex()
					var foundstuff = false
					for cursor.Next(ctxMongo) {
						blogPrimitive := &bson.D{}
						err = cursor.Decode(blogPrimitive)
						if err != nil {
							return nil, err
						}
						blogData = blogPrimitive.Map()
						id := blogData["_id"].(primitive.ObjectID)
						blogData["date"] = objectidtimestamp(id).Format(dateFormat)
						blogData["id"] = idstr
						delete(blogData, "_id")
						foundstuff = true
						break
					}
					if !foundstuff {
						return nil, errors.New("blog not found with given id")
					}
					_, err = blogCollection.DeleteOne(ctxMongo, bson.M{
						"_id": id,
					})
					if err != nil {
						return nil, err
					}
					_, err = elasticClient.Delete().
						Index(blogElasticIndex).
						Type("blog").
						Id(idstring).
						Do(ctxElastic)
					if err != nil {
						return nil, err
					}
					pictureids := blogData["images"].([]string)
					for _, pictureid := range pictureids {
						logger.Info("pictureid: " + pictureid + ", blogid: " + idstr)
						fileobj := blogImageBucket.Object(blogPictureIndex + "/" + idstr + "/" + pictureid)
						if err := fileobj.Delete(ctxStorage); err != nil {
							return nil, err
						}
					}
					return blogData, nil
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
					_, err := validateAdmin(params.Context.Value("token").(string))
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
