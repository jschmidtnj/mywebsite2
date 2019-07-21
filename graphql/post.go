package main

import (
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/kinds"
)

func parseLiteral(astValue ast.Value) interface{} {
	kind := astValue.GetKind()

	switch kind {
	case kinds.StringValue:
		return astValue.GetValue()
	case kinds.BooleanValue:
		return astValue.GetValue()
	case kinds.IntValue:
		return astValue.GetValue()
	case kinds.FloatValue:
		return astValue.GetValue()
	case kinds.ObjectValue:
		obj := make(map[string]interface{})
		for _, v := range astValue.GetValue().([]*ast.ObjectField) {
			obj[v.Name.Value] = parseLiteral(v.Value)
		}
		return obj
	case kinds.ListValue:
		astValueList := astValue.GetValue().([]ast.Value)
		list := make([]interface{}, len(astValueList))
		for i, v := range astValueList {
			list[i] = parseLiteral(v)
		}
		return list
	default:
		return nil
	}
}

// JSON json type
var jsonType = graphql.NewScalar(
	graphql.ScalarConfig{
		Name:        "JSON",
		Description: "The `JSON` scalar type represents JSON values as specified by [ECMA-404](http://www.ecma-international.org/publications/files/ECMA-ST/ECMA-404.pdf)",
		Serialize: func(value interface{}) interface{} {
			return value
		},
		ParseValue: func(value interface{}) interface{} {
			return value
		},
		ParseLiteral: parseLiteral,
	},
)

// PostType graphql post type is a post object
var PostType *graphql.Object = graphql.NewObject(graphql.ObjectConfig{
	Name: "Post",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
		},
		"title": &graphql.Field{
			Type: graphql.String,
		},
		"caption": &graphql.Field{
			Type: graphql.String,
		},
		"content": &graphql.Field{
			Type: graphql.String,
		},
		"author": &graphql.Field{
			Type: graphql.String,
		},
		"color": &graphql.Field{
			Type: graphql.String,
		},
		"tags": &graphql.Field{
			Type: graphql.NewList(graphql.String),
		},
		"categories": &graphql.Field{
			Type: graphql.NewList(graphql.String),
		},
		"views": &graphql.Field{
			Type: graphql.Int,
		},
		"date": &graphql.Field{
			Type: graphql.String,
		},
		"heroimage": &graphql.Field{
			Type: graphql.String,
		},
		"tileimage": &graphql.Field{
			Type: graphql.String,
		},
		"images": &graphql.Field{
			Type: graphql.NewList(graphql.String),
		},
		"files": &graphql.Field{
			Type: graphql.NewList(graphql.String),
		},
		"comments": &graphql.Field{
			Type: jsonType,
		},
		"shortlink": &graphql.Field{
			Type: graphql.String,
		},
	},
})
