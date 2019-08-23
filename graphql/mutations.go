package main

import (
	"errors"

	"net/url"

	"cloud.google.com/go/storage"
	"github.com/graphql-go/graphql"

	// medium "github.com/medium/medium-sdk-go"
	// "gopkg.in/russross/blackfriday.v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	// "github.com/Depado/bfchroma"
)

func interfaceListToStringList(interfaceList []interface{}) ([]string, error) {
	result := make([]string, len(interfaceList))
	for i, item := range interfaceList {
		itemStr, ok := item.(string)
		if !ok {
			return nil, errors.New("item in list cannot be cast to string")
		}
		itemStrDecoded, err := url.QueryUnescape(itemStr)
		if err != nil {
			return nil, err
		}
		result[i] = itemStrDecoded
	}
	return result, nil
}

func deleteAccount(idstring string) (interface{}, error) {
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
	shortlinksInterface, ok := userData["shortlinks"].([]interface{})
	if !ok {
		return nil, errors.New("unable to cast shortlinks to array")
	}
	shortlinks, err := interfaceListToStringList(shortlinksInterface)
	if err != nil {
		return nil, err
	}
	for _, link := range shortlinks {
		err = deleteShortLink(link)
		if err != nil {
			return nil, err
		}
	}
	_, err = userCollection.DeleteOne(ctxMongo, bson.M{
		"_id": id,
	})
	if err != nil {
		return nil, err
	}

	return userData, nil
}

func rootMutation() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"addPost": &graphql.Field{
				Type:        PostType,
				Description: "Create a Post Post",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"type": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"title": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"caption": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"content": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"author": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"color": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"tags": &graphql.ArgumentConfig{
						Type: graphql.NewList(graphql.String),
					},
					"categories": &graphql.ArgumentConfig{
						Type: graphql.NewList(graphql.String),
					},
					"heroimage": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"tileimage": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"images": &graphql.ArgumentConfig{
						Type: graphql.NewList(graphql.String),
					},
					"files": &graphql.ArgumentConfig{
						Type: graphql.NewList(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					_, err := validateAdmin(params.Context.Value(tokenKey).(string))
					if err != nil {
						return nil, err
					}
					if params.Args["id"] == nil || params.Args["title"] == nil || params.Args["content"] == nil ||
						params.Args["author"] == nil || params.Args["type"] == nil ||
						params.Args["heroimage"] == nil || params.Args["images"] == nil ||
						params.Args["files"] == nil || params.Args["caption"] == nil ||
						params.Args["color"] == nil || params.Args["tags"] == nil ||
						params.Args["categories"] == nil || params.Args["tileimage"] == nil {
						return nil, errors.New("title or content or author or type or heroimage or images or files or caption or color or tags or categories or tileimage not provided")
					}
					idstring, ok := params.Args["id"].(string)
					if !ok {
						return nil, errors.New("cannot cast id to string")
					}
					id, err := primitive.ObjectIDFromHex(idstring)
					if err != nil {
						return nil, err
					}
					title, ok := params.Args["title"].(string)
					if !ok {
						return nil, errors.New("problem casting title to string")
					}
					caption, ok := params.Args["caption"].(string)
					if !ok {
						return nil, errors.New("problem casting caption to string")
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
					if !validType(thetype) {
						return nil, errors.New("invalid type given")
					}
					color, ok := params.Args["color"].(string)
					if !ok {
						return nil, errors.New("problem casting color to string")
					}
					decodedColor, err := url.QueryUnescape(color)
					if err != nil {
						return nil, err
					}
					if !validHexcode.MatchString(decodedColor) {
						return nil, errors.New("invalid hex code for color")
					}
					tagsinterface, ok := params.Args["tags"].([]interface{})
					if !ok {
						return nil, errors.New("problem casting tags to interface array")
					}
					tags, err := interfaceListToStringList(tagsinterface)
					if err != nil {
						return nil, err
					}
					categoriesinterface, ok := params.Args["categories"].([]interface{})
					if !ok {
						return nil, errors.New("problem casting categories to interface array")
					}
					categories, err := interfaceListToStringList(categoriesinterface)
					if err != nil {
						return nil, err
					}
					heroimage, ok := params.Args["heroimage"].(string)
					if !ok {
						return nil, errors.New("problem casting heroimage to string")
					}
					tileimage, ok := params.Args["tileimage"].(string)
					if !ok {
						return nil, errors.New("problem casting tileimage to string")
					}
					imagesinterface, ok := params.Args["images"].([]interface{})
					if !ok {
						return nil, errors.New("problem casting images to interface array")
					}
					images, err := interfaceListToStringList(imagesinterface)
					if err != nil {
						return nil, err
					}
					filesinterface, ok := params.Args["files"].([]interface{})
					if !ok {
						return nil, errors.New("problem casting files to interface array")
					}
					files, err := interfaceListToStringList(filesinterface)
					if err != nil {
						return nil, err
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
					shortlink, err := generateShortLink(websiteURL + "/" + thetype + "/" + idstring)
					if err != nil {
						return nil, err
					}
					postData := bson.M{
						"_id":        id,
						"title":      title,
						"caption":    caption,
						"content":    content,
						"author":     author,
						"color":      color,
						"tags":       tags,
						"categories": categories,
						"views":      0,
						"heroimage":  heroimage,
						"tileimage":  tileimage,
						"images":     images,
						"files":      files,
						"comments":   []string{},
						"shortlink":  shortlink,
					}
					_, err = mongoCollection.InsertOne(ctxMongo, postData)
					if err != nil {
						return nil, err
					}
					timestamp := objectidtimestamp(id)
					postData["date"] = timestamp.Unix()
					delete(postData, "_id")
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
					"id": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"type": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"title": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"caption": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"content": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"author": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"color": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"tags": &graphql.ArgumentConfig{
						Type: graphql.NewList(graphql.String),
					},
					"categories": &graphql.ArgumentConfig{
						Type: graphql.NewList(graphql.String),
					},
					"heroimage": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"tileimage": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"images": &graphql.ArgumentConfig{
						Type: graphql.NewList(graphql.String),
					},
					"files": &graphql.ArgumentConfig{
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
					if params.Args["caption"] != nil {
						caption, ok := params.Args["caption"].(string)
						if !ok {
							return nil, errors.New("problem casting caption to string")
						}
						updateData["caption"] = caption
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
					if params.Args["color"] != nil {
						color, ok := params.Args["color"].(string)
						if !ok {
							return nil, errors.New("problem casting color to string")
						}
						decodedColor, err := url.QueryUnescape(color)
						if err != nil {
							return nil, err
						}
						if !validHexcode.MatchString(decodedColor) {
							return nil, errors.New("invalid hex code for color")
						}
						updateData["color"] = color
					}
					if params.Args["tags"] != nil {
						tagsinterface, ok := params.Args["tags"].([]interface{})
						if !ok {
							return nil, errors.New("problem casting tags to interface array")
						}
						tags, err := interfaceListToStringList(tagsinterface)
						if err != nil {
							return nil, err
						}
						updateData["tags"] = tags
					}
					if params.Args["categories"] != nil {
						categoriesinterface, ok := params.Args["categories"].([]interface{})
						if !ok {
							return nil, errors.New("problem casting categories to interface array")
						}
						categories, err := interfaceListToStringList(categoriesinterface)
						if err != nil {
							return nil, err
						}
						updateData["categories"] = categories
					}
					if params.Args["heroimage"] != nil {
						heroimage, ok := params.Args["heroimage"].(string)
						if !ok {
							return nil, errors.New("problem casting heroimage to string")
						}
						updateData["heroimage"] = heroimage
					}
					if params.Args["tileimage"] != nil {
						tileimage, ok := params.Args["tileimage"].(string)
						if !ok {
							return nil, errors.New("problem casting tileimage to string")
						}
						updateData["tileimage"] = tileimage
					}
					if params.Args["images"] != nil {
						imagesinterface, ok := params.Args["images"].([]interface{})
						if !ok {
							return nil, errors.New("problem casting images to interface array")
						}
						images, err := interfaceListToStringList(imagesinterface)
						if err != nil {
							return nil, err
						}
						updateData["images"] = images
					}
					if params.Args["files"] != nil {
						filesinterface, ok := params.Args["files"].([]interface{})
						if !ok {
							return nil, errors.New("problem casting files to interface array")
						}
						files, err := interfaceListToStringList(filesinterface)
						if err != nil {
							return nil, err
						}
						updateData["files"] = files
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
					err = deleteShortLink(postData["shortlink"].(string))
					if err != nil {
						return nil, err
					}
					heroimageid, ok := postData["heroimage"].(string)
					if !ok {
						return nil, errors.New("cannot convert heroimage to string")
					}
					if len(heroimageid) > 0 {
						var heroobjblur *storage.ObjectHandle
						var heroobjoriginal *storage.ObjectHandle
						if thetype == "blog" {
							heroobjblur = storageBucket.Object(blogImageIndex + "/" + idstr + "/" + heroimageid + "/blur")
							heroobjoriginal = storageBucket.Object(blogImageIndex + "/" + idstr + "/" + heroimageid + "/original")
						} else {
							heroobjblur = storageBucket.Object(projectImageIndex + "/" + idstr + "/" + heroimageid + "/blur")
							heroobjoriginal = storageBucket.Object(projectImageIndex + "/" + idstr + "/" + heroimageid + "/original")
						}
						if err := heroobjblur.Delete(ctxStorage); err != nil {
							return nil, err
						}
						if err := heroobjoriginal.Delete(ctxStorage); err != nil {
							return nil, err
						}
					}
					tileimageid, ok := postData["tileimage"].(string)
					if !ok {
						return nil, errors.New("cannot convert tileimage to string")
					}
					if len(tileimageid) > 0 {
						var tileobjblur *storage.ObjectHandle
						var tileobjoriginal *storage.ObjectHandle
						if thetype == "blog" {
							tileobjblur = storageBucket.Object(blogImageIndex + "/" + idstr + "/" + tileimageid + "/blur")
							tileobjoriginal = storageBucket.Object(blogImageIndex + "/" + idstr + "/" + tileimageid + "/original")
						} else {
							tileobjblur = storageBucket.Object(projectImageIndex + "/" + idstr + "/" + tileimageid + "/blur")
							tileobjoriginal = storageBucket.Object(projectImageIndex + "/" + idstr + "/" + tileimageid + "/original")
						}
						if err := tileobjblur.Delete(ctxStorage); err != nil {
							return nil, err
						}
						if err := tileobjoriginal.Delete(ctxStorage); err != nil {
							return nil, err
						}
					}
					primativeimageids, ok := postData["images"].(primitive.A)
					if !ok {
						return nil, errors.New("cannot convert imageids to primitive")
					}
					for _, primativeimageid := range primativeimageids {
						imageid := primativeimageid.(string)
						var imageobjblur *storage.ObjectHandle
						var imageobjoriginal *storage.ObjectHandle
						if thetype == "blog" {
							imageobjblur = storageBucket.Object(blogImageIndex + "/" + idstr + "/" + imageid + "/blur")
							imageobjoriginal = storageBucket.Object(blogImageIndex + "/" + idstr + "/" + imageid + "/original")
						} else {
							imageobjblur = storageBucket.Object(projectImageIndex + "/" + idstr + "/" + imageid + "/blur")
							imageobjoriginal = storageBucket.Object(projectImageIndex + "/" + idstr + "/" + imageid + "/original")
						}
						if err := imageobjblur.Delete(ctxStorage); err != nil {
							return nil, err
						}
						if err := imageobjoriginal.Delete(ctxStorage); err != nil {
							return nil, err
						}
					}
					primativefileids, ok := postData["files"].(primitive.A)
					if !ok {
						return nil, errors.New("cannot convert fileids to primitive")
					}
					for _, primativefileid := range primativefileids {
						fileid := primativefileid.(string)
						var fileobj *storage.ObjectHandle
						if thetype == "blog" {
							fileobj = storageBucket.Object(blogFileIndex + "/" + idstr + "/" + fileid)
						} else {
							fileobj = storageBucket.Object(projectFileIndex + "/" + idstr + "/" + fileid)
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
			"addShortlink": &graphql.Field{
				Type:        ShortLinkType,
				Description: "Add Short Link",
				Args: graphql.FieldConfigArgument{
					"link": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"recaptcha": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					claims, err := validateLoggedIn(params.Context.Value(tokenKey).(string))
					if err != nil {
						return nil, err
					}
					idstring, ok := claims["id"].(string)
					if !ok {
						return nil, errors.New("cannot cast id to string")
					}
					id, err := primitive.ObjectIDFromHex(idstring)
					if err != nil {
						return nil, err
					}
					if params.Args["link"] == nil {
						return nil, errors.New("link not provided")
					}
					linkstring, ok := params.Args["link"].(string)
					if !ok {
						return nil, errors.New("cannot cast link to string")
					}
					if params.Args["recaptcha"] == nil {
						return nil, errors.New("recaptcha not provided")
					}
					recaptchastring, ok := params.Args["recaptcha"].(string)
					if !ok {
						return nil, errors.New("cannot cast recaptcha to string")
					}
					err = verifyRecaptcha(recaptchastring, shortlinkRecaptchaSecret)
					if err != nil {
						return nil, err
					}
					linkid, err := generateShortLink(linkstring)
					if err != nil {
						return nil, err
					}
					_, err = userCollection.UpdateOne(ctxMongo, bson.M{
						"_id": id,
					}, bson.M{
						"$push": bson.M{
							"shortlinks": linkid,
						},
					})
					if err != nil {
						return nil, err
					}
					shortLinkData := bson.M{
						"id":   linkid,
						"link": linkstring,
					}
					return shortLinkData, nil
				},
			},
			"removeShortlink": &graphql.Field{
				Type:        ShortLinkType,
				Description: "Remove Short Link",
				Args: graphql.FieldConfigArgument{
					"linkid": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					claims, err := validateLoggedIn(params.Context.Value(tokenKey).(string))
					if err != nil {
						return nil, err
					}
					idstring, ok := claims["id"].(string)
					if !ok {
						return nil, errors.New("cannot cast id to string")
					}
					id, err := primitive.ObjectIDFromHex(idstring)
					if err != nil {
						return nil, err
					}
					if params.Args["linkid"] == nil {
						return nil, errors.New("link not provided")
					}
					linkid, ok := params.Args["linkid"].(string)
					if !ok {
						return nil, errors.New("cannot cast linkid to string")
					}
					shortLink, err := getShortLink(linkid)
					if err != nil {
						return nil, err
					}
					err = deleteShortLink(linkid)
					if err != nil {
						return nil, err
					}
					_, err = userCollection.UpdateOne(ctxMongo, bson.M{
						"_id": id,
					}, bson.M{
						"$pull": bson.M{
							"shortlinks": linkid,
						},
					})
					if err != nil {
						return nil, err
					}
					shortLinkData := bson.M{
						"id":   linkid,
						"link": shortLink,
					}
					return shortLinkData, nil
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
					return deleteAccount(idstring)
				},
			},
			"deleteAccount": &graphql.Field{
				Type:        AccountType,
				Description: "Delete a User",
				Args:        graphql.FieldConfigArgument{},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					claims, err := validateLoggedIn(params.Context.Value(tokenKey).(string))
					if err != nil {
						return nil, err
					}
					idstring, ok := claims["id"].(string)
					if !ok {
						return nil, errors.New("cannot cast id to string")
					}
					return deleteAccount(idstring)
				},
			},
		},
	})
}
