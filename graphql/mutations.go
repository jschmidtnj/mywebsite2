package main

import (
	"errors"
	"github.com/graphql-go/graphql"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func RootMutation() *graphql.Object {
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
					_, err := ValidateAdmin(params.Context.Value("token").(string))
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
						"title":   title,
						"content": content,
						"author":  author,
						"views":   0,
					}
					res, err := BlogCollection.InsertOne(CTX, blogData)
					if err != nil {
						return nil, err
					}
					id := res.InsertedID.(primitive.ObjectID)
					idstring := id.Hex()
					timestamp := objectidtimestamp(id)
					blogData["date"] = timestamp.Unix()
					_, err = Elastic.Index().
						Index(BlogElasticIndex).
						Type("blog").
						Id(idstring).
						BodyJson(blogData).
						Do(CTXElastic)
					if err != nil {
						return nil, err
					}
					blogData["date"] = timestamp.Format(DateFormat)
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
					_, err := ValidateAdmin(params.Context.Value("token").(string))
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
					_, err = Elastic.Update().
						Index(BlogElasticIndex).
						Type("blog").
						Id(idstring).
						Doc(updateData).
						Do(CTXElastic)
					if err != nil {
						return nil, err
					}
					_, err = BlogCollection.UpdateOne(CTX, bson.M{
						"_id": id,
					}, bson.M{
						"$set": updateData,
					})
					if err != nil {
						return nil, err
					}
					cursor, err := BlogCollection.Find(CTX, bson.M{
						"_id": id,
					})
					defer cursor.Close(CTX)
					if err != nil {
						return nil, err
					}
					var blogData map[string]interface{}
					var foundstuff = false
					for cursor.Next(CTX) {
						blogPrimitive := &bson.D{}
						err = cursor.Decode(blogPrimitive)
						if err != nil {
							return nil, err
						}
						blogData = blogPrimitive.Map()
						id := blogData["_id"].(primitive.ObjectID)
						blogData["date"] = objectidtimestamp(id).Format(DateFormat)
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
					_, err := ValidateAdmin(params.Context.Value("token").(string))
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
					cursor, err := BlogCollection.Find(CTX, bson.M{
						"_id": id,
					})
					if err != nil {
						return nil, err
					}
					defer cursor.Close(CTX)
					var blogData map[string]interface{}
					var foundstuff = false
					for cursor.Next(CTX) {
						blogPrimitive := &bson.D{}
						err = cursor.Decode(blogPrimitive)
						if err != nil {
							return nil, err
						}
						blogData = blogPrimitive.Map()
						id := blogData["_id"].(primitive.ObjectID)
						blogData["date"] = objectidtimestamp(id).Format(DateFormat)
						blogData["id"] = id.Hex()
						delete(blogData, "_id")
						foundstuff = true
						break
					}
					if !foundstuff {
						return nil, errors.New("blog not found with given id")
					}
					_, err = BlogCollection.DeleteOne(CTX, bson.M{
						"_id": id,
					})
					if err != nil {
						return nil, err
					}
					_, err = Elastic.Delete().
						Index(BlogElasticIndex).
						Type("blog").
						Id(idstring).
						Do(CTXElastic)
					if err != nil {
						return nil, err
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
					_, err := ValidateAdmin(params.Context.Value("token").(string))
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
					cursor, err := UserCollection.Find(CTX, bson.M{
						"_id": id,
					})
					defer cursor.Close(CTX)
					if err != nil {
						return nil, err
					}
					var userData map[string]interface{}
					var foundstuff = false
					for cursor.Next(CTX) {
						userPrimitive := &bson.D{}
						err = cursor.Decode(userPrimitive)
						if err != nil {
							return nil, err
						}
						userData = userPrimitive.Map()
						id := userData["_id"].(primitive.ObjectID)
						userData["date"] = objectidtimestamp(id).Format(DateFormat)
						userData["id"] = id.Hex()
						delete(userData, "_id")
						foundstuff = true
						break
					}
					if !foundstuff {
						return nil, errors.New("user not found with given id")
					}
					_, err = UserCollection.DeleteOne(CTX, bson.M{
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
