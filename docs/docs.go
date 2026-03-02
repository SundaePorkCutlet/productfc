// Package docs provides Swagger documentation for PRODUCTFC API.
// Run `swag init -g main.go` in PRODUCTFC directory to regenerate from annotations.
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "swagger": "2.0",
    "info": {
        "title": "PRODUCTFC API",
        "description": "Product catalog, categories, and inventory for Go Commerce.",
        "version": "1.0"
    },
    "host": "localhost:28081",
    "basePath": "/",
    "paths": {
        "/ping": {
            "get": {
                "summary": "Ping",
                "responses": {"200": {"description": "pong"}}
            }
        },
        "/v1/products/search": {
            "get": {
                "summary": "Search products",
                "parameters": [
                    {"in": "query", "name": "name", "type": "string"},
                    {"in": "query", "name": "category", "type": "string"},
                    {"in": "query", "name": "min_price", "type": "number"},
                    {"in": "query", "name": "max_price", "type": "number"},
                    {"in": "query", "name": "page", "type": "integer"},
                    {"in": "query", "name": "page_size", "type": "integer"},
                    {"in": "query", "name": "order_by", "type": "string"},
                    {"in": "query", "name": "sort", "type": "string"}
                ],
                "responses": {"200": {"description": "List of products with pagination"}}
            }
        },
        "/v1/products/{id}": {
            "get": {
                "summary": "Get product by ID",
                "parameters": [{"in": "path", "name": "id", "type": "integer", "required": true}],
                "responses": {"200": {"description": "Product"}, "400": {"description": "Invalid id"}}
            }
        },
        "/v1/product-categories/{id}": {
            "get": {
                "summary": "Get category by ID",
                "parameters": [{"in": "path", "name": "id", "type": "integer", "required": true}],
                "responses": {"200": {"description": "Category"}}
            }
        },
        "/api/v1/products": {
            "post": {
                "security": [{"BearerAuth": []}],
                "summary": "Create product",
                "parameters": [{"in": "body", "name": "body", "schema": {"type": "object"}}],
                "responses": {"201": {"description": "Created"}, "401": {"description": "Unauthorized"}}
            }
        },
        "/api/v1/products/{id}": {
            "put": {
                "security": [{"BearerAuth": []}],
                "summary": "Update product",
                "parameters": [{"in": "path", "name": "id", "type": "integer"}, {"in": "body", "name": "body", "schema": {"type": "object"}}],
                "responses": {"200": {"description": "OK"}, "401": {"description": "Unauthorized"}}
            },
            "delete": {
                "security": [{"BearerAuth": []}],
                "summary": "Delete product",
                "parameters": [{"in": "path", "name": "id", "type": "integer"}],
                "responses": {"200": {"description": "OK"}, "401": {"description": "Unauthorized"}}
            }
        },
        "/api/v1/product-categories": {
            "post": {
                "security": [{"BearerAuth": []}],
                "summary": "Create category",
                "parameters": [{"in": "body", "name": "body", "schema": {"type": "object"}}],
                "responses": {"201": {"description": "Created"}, "401": {"description": "Unauthorized"}}
            }
        },
        "/api/v1/product-categories/{id}": {
            "put": {
                "security": [{"BearerAuth": []}],
                "summary": "Update category",
                "parameters": [{"in": "path", "name": "id", "type": "integer"}, {"in": "body", "name": "body", "schema": {"type": "object"}}],
                "responses": {"200": {"description": "OK"}, "401": {"description": "Unauthorized"}}
            },
            "delete": {
                "security": [{"BearerAuth": []}],
                "summary": "Delete category",
                "parameters": [{"in": "path", "name": "id", "type": "integer"}],
                "responses": {"200": {"description": "OK"}, "401": {"description": "Unauthorized"}}
            }
        }
    },
    "securityDefinitions": {
        "BearerAuth": {"type": "apiKey", "name": "Authorization", "in": "header"}
    }
}`

func init() {
	swag.Register(swag.Name, &s{})
}

type s struct{}

func (s *s) ReadDoc() string {
	return docTemplate
}
