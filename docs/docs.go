// Code generated by swaggo/swag. DO NOT EDIT.

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
        "/ping": {
            "get": {
                "description": "Ping",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Ping"
                ],
                "summary": "Ping",
                "responses": {
                    "200": {
                        "description": "pong",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/scan/job-name": {
            "get": {
                "description": "Get the name of the last / current job",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Scan"
                ],
                "summary": "Get Job Name",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/scan/logs": {
            "get": {
                "description": "Get the logs of the last / current job",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Scan"
                ],
                "summary": "Get Job Logs",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/controllers.jobLogResponse"
                            }
                        }
                    }
                }
            }
        },
        "/scan/movie": {
            "post": {
                "description": "Scan movies from the configured folder",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Scan"
                ],
                "summary": "Scan Movies",
                "responses": {
                    "200": {
                        "description": "Scan started",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/controllers.errorResponse"
                        }
                    }
                }
            }
        },
        "/scan/pop-logs": {
            "get": {
                "description": "Get the logs of the last / current job",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Scan"
                ],
                "summary": "Get Job Logs",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/controllers.jobLogResponse"
                            }
                        }
                    }
                }
            }
        },
        "/scan/tv": {
            "post": {
                "description": "Scan TV Shows from the configured folder",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Scan"
                ],
                "summary": "Scan TV Shows",
                "responses": {
                    "200": {
                        "description": "Scan started",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/controllers.errorResponse"
                        }
                    }
                }
            }
        },
        "/upload/movie": {
            "post": {
                "description": "Upload movies from the configured folder",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Upload"
                ],
                "summary": "Upload Movies",
                "parameters": [
                    {
                        "type": "file",
                        "description": "Files to upload",
                        "name": "upload[]",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controllers.uploadResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/controllers.errorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/controllers.errorResponse"
                        }
                    }
                }
            }
        },
        "/upload/tv": {
            "post": {
                "description": "Upload TV Shows from the configured folder",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Upload"
                ],
                "summary": "Upload TV Shows",
                "parameters": [
                    {
                        "type": "file",
                        "description": "Files to upload",
                        "name": "upload[]",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controllers.uploadResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/controllers.errorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/controllers.errorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "controllers.errorResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                }
            }
        },
        "controllers.jobLogResponse": {
            "type": "object",
            "properties": {
                "date": {
                    "type": "string",
                    "example": "2021-01-01"
                },
                "job_name": {
                    "type": "string",
                    "example": "upload movie"
                },
                "message": {
                    "type": "string",
                    "example": "Uploading movie test.mp4"
                }
            }
        },
        "controllers.uploadResponse": {
            "type": "object",
            "properties": {
                "count": {
                    "type": "integer"
                },
                "message": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Media Indexer API",
	Description:      "This is the API for the Media Indexer application",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
