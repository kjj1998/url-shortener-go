// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "url": "http://www.swagger.io/support",
            "email": "support@swagger.io"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/data/shorten": {
            "post": {
                "description": "generate shortened urls",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "example"
                ],
                "summary": "generate shortened urls",
                "parameters": [
                    {
                        "description": "Add URL for shortening",
                        "name": "longUrl",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.LongUrl"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/models.ShortenedUrl"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.HTTPError"
                        }
                    }
                }
            }
        },
        "/{shortUrl}": {
            "get": {
                "description": "redirect shortened urls to the actual urls",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    ""
                ],
                "summary": "redirect shortened urls to the actual urls",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Short URL",
                        "name": "shortUrl",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "307": {
                        "description": "Temporary Redirect"
                    }
                }
            }
        }
    },
    "definitions": {
        "models.LongUrl": {
            "type": "object",
            "properties": {
                "longUrl": {
                    "type": "string"
                }
            }
        },
        "models.ShortenedUrl": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "longUrl": {
                    "type": "string"
                },
                "shortUrl": {
                    "type": "string"
                }
            }
        },
        "utils.HTTPError": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer",
                    "example": 400
                },
                "message": {
                    "type": "string",
                    "example": "status bad request"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "URL shortening API",
	Description:      "This is a URL shortener service.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
