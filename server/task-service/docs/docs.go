// Package docs provides Swagger documentation for the Task Service
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/tasks": {
            "get": {
                "description": "List tasks for the authenticated user (or all tasks for admin)",
                "produces": ["application/json"],
                "tags": ["tasks"],
                "summary": "List tasks",
                "security": [{"BearerAuth": []}],
                "parameters": [
                    {
                        "type": "integer",
                        "default": 1,
                        "description": "Page number",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "default": 10,
                        "description": "Items per page",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "enum": ["PENDING", "IN_PROGRESS", "DONE"],
                        "description": "Filter by status",
                        "name": "status",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/APIResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/APIResponse"
                        }
                    }
                }
            },
            "post": {
                "description": "Create a new task for the authenticated user",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "tags": ["tasks"],
                "summary": "Create a new task",
                "security": [{"BearerAuth": []}],
                "parameters": [
                    {
                        "description": "Task details",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/CreateTaskRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/APIResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/APIResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/APIResponse"
                        }
                    }
                }
            }
        },
        "/tasks/{id}": {
            "get": {
                "description": "Get a specific task by its ID",
                "produces": ["application/json"],
                "tags": ["tasks"],
                "summary": "Get task by ID",
                "security": [{"BearerAuth": []}],
                "parameters": [
                    {
                        "type": "string",
                        "description": "Task ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/APIResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/APIResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/APIResponse"
                        }
                    }
                }
            },
            "put": {
                "description": "Update an existing task",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "tags": ["tasks"],
                "summary": "Update a task",
                "security": [{"BearerAuth": []}],
                "parameters": [
                    {
                        "type": "string",
                        "description": "Task ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Updated task details",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/UpdateTaskRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/APIResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/APIResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/APIResponse"
                        }
                    }
                }
            },
            "delete": {
                "description": "Delete an existing task",
                "produces": ["application/json"],
                "tags": ["tasks"],
                "summary": "Delete a task",
                "security": [{"BearerAuth": []}],
                "parameters": [
                    {
                        "type": "string",
                        "description": "Task ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/APIResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/APIResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/APIResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "CreateTaskRequest": {
            "type": "object",
            "required": ["title"],
            "properties": {
                "title": {
                    "type": "string",
                    "minLength": 1,
                    "maxLength": 255
                },
                "description": {
                    "type": "string",
                    "maxLength": 5000
                },
                "status": {
                    "type": "string",
                    "enum": ["PENDING", "IN_PROGRESS", "DONE"]
                }
            }
        },
        "UpdateTaskRequest": {
            "type": "object",
            "properties": {
                "title": {
                    "type": "string",
                    "minLength": 1,
                    "maxLength": 255
                },
                "description": {
                    "type": "string",
                    "maxLength": 5000
                },
                "status": {
                    "type": "string",
                    "enum": ["PENDING", "IN_PROGRESS", "DONE"]
                }
            }
        },
        "APIResponse": {
            "type": "object",
            "properties": {
                "success": {
                    "type": "boolean"
                },
                "data": {
                    "type": "object"
                },
                "message": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "BearerAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8082",
	BasePath:         "/api/v1",
	Schemes:          []string{"http"},
	Title:            "Prime Trade Task Service API",
	Description:      "Task Management Service for Prime Trade",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
