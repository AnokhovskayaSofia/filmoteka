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
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/actors": {
            "get": {
                "security": [
                    {
                        "BasicAuth": []
                    }
                ],
                "description": "Availible only for authenticated user, getting actors list from db",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "actors"
                ],
                "summary": "List actors",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api_models.ActorsResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResponse"
                        }
                    }
                }
            },
            "put": {
                "security": [
                    {
                        "BasicAuth": []
                    }
                ],
                "description": "Availible only for admin user, updating actor using id from request params and return actor",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "actors"
                ],
                "summary": "Update actor",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Actors Id",
                        "name": "actorID",
                        "in": "query",
                        "required": true
                    },
                    {
                        "description": "actor info",
                        "name": "Actor",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/db.Actor"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api_models.ActorResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResponse"
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "BasicAuth": []
                    }
                ],
                "description": "Availible only for admin user, creating actor using data from request body and return new actor",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "actors"
                ],
                "summary": "Create actor",
                "parameters": [
                    {
                        "description": "actor info",
                        "name": "Actor",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/db.Actor"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api_models.ActorResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResponse"
                        }
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "BasicAuth": []
                    }
                ],
                "description": "Availible only for admin user, deleting actor using id from request params",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "actors"
                ],
                "summary": "Delete actor",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Actors Id",
                        "name": "actorID",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/films": {
            "get": {
                "security": [
                    {
                        "BasicAuth": []
                    }
                ],
                "description": "Availible only for authenticated user, getting films list, they can be sorted by fields, default is rate. Also you can use filters in field.value template.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "films"
                ],
                "summary": "Get films list",
                "parameters": [
                    {
                        "type": "string",
                        "example": "name",
                        "description": "Sort by field, default rate",
                        "name": "sortBy",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "example": "name.Name1",
                        "description": "Filter by field (field.value), can be user all except actors",
                        "name": "filter",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api_models.FilmsResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResponse"
                        }
                    }
                }
            },
            "put": {
                "security": [
                    {
                        "BasicAuth": []
                    }
                ],
                "description": "Availible only for admin user, updating film using data from request body and return new film",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "films"
                ],
                "summary": "Update film",
                "parameters": [
                    {
                        "description": "film info",
                        "name": "Film",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/db.Film"
                        }
                    },
                    {
                        "type": "string",
                        "description": "Film Id",
                        "name": "filmID",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api_models.FilmResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResponse"
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "BasicAuth": []
                    }
                ],
                "description": "Availible only for admin user, creating film using data from request body and return new film",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "films"
                ],
                "summary": "Create film",
                "parameters": [
                    {
                        "description": "film info",
                        "name": "Film",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/db.Film"
                        }
                    },
                    {
                        "type": "string",
                        "description": "Film Id",
                        "name": "filmID",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api_models.FilmResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResponse"
                        }
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "BasicAuth": []
                    }
                ],
                "description": "Availible only for admin user, deleting film by id from params",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "films"
                ],
                "summary": "Delete film",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Film Id",
                        "name": "filmID",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "api.ErrorResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string",
                    "example": "Error msg"
                },
                "success": {
                    "type": "boolean",
                    "example": false
                }
            }
        },
        "api_models.ActorResponse": {
            "type": "object",
            "properties": {
                "actor": {
                    "$ref": "#/definitions/filmoteka_db.Actor"
                },
                "error": {
                    "type": "string"
                },
                "success": {
                    "type": "boolean"
                }
            }
        },
        "api_models.ActorsResponse": {
            "type": "object",
            "properties": {
                "actors": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/filmoteka_db.Actor"
                    }
                },
                "error": {
                    "type": "string"
                },
                "success": {
                    "type": "boolean"
                }
            }
        },
        "api_models.FilmResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                },
                "film": {
                    "$ref": "#/definitions/filmoteka_db.Film"
                },
                "success": {
                    "type": "boolean"
                }
            }
        },
        "api_models.FilmsResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                },
                "film": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/filmoteka_db.Film"
                    }
                },
                "success": {
                    "type": "boolean"
                }
            }
        },
        "db.Actor": {
            "type": "object",
            "properties": {
                "birthday": {
                    "type": "string"
                },
                "films": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/db.Film"
                    }
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "sex": {
                    "type": "string",
                    "enum": [
                        "male",
                        "female"
                    ]
                }
            }
        },
        "db.Film": {
            "type": "object",
            "properties": {
                "actors": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/db.Actor"
                    }
                },
                "date": {
                    "type": "string"
                },
                "description": {
                    "type": "string",
                    "maxLength": 1000
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string",
                    "maxLength": 150,
                    "minLength": 1
                },
                "rate": {
                    "type": "integer",
                    "maximum": 10,
                    "minimum": 0
                }
            }
        },
        "filmoteka_db.Actor": {
            "type": "object",
            "properties": {
                "birthday": {
                    "type": "string"
                },
                "films": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/db.Film"
                    }
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "sex": {
                    "type": "string",
                    "enum": [
                        "male",
                        "female"
                    ]
                }
            }
        },
        "filmoteka_db.Film": {
            "type": "object",
            "properties": {
                "actors": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/db.Actor"
                    }
                },
                "date": {
                    "type": "string"
                },
                "description": {
                    "type": "string",
                    "maxLength": 1000
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string",
                    "maxLength": 150,
                    "minLength": 1
                },
                "rate": {
                    "type": "integer",
                    "maximum": 10,
                    "minimum": 0
                }
            }
        }
    },
    "securityDefinitions": {
        "BasicAuth": {
            "type": "basic"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8084",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Filmoteka API",
	Description:      "This is a sample Filmoteka server.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
