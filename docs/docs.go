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
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controllers.movieScanResponse"
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
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controllers.tvScanResponse"
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
        "controllers.movieScanResponse": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/features.MovieScannerResult"
                    }
                }
            }
        },
        "controllers.tvScanResponse": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/features.TVScannerResult"
                    }
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
        },
        "features.MovieScannerResult": {
            "type": "object",
            "properties": {
                "destination": {
                    "description": "Full destination path of the moved file.",
                    "type": "string"
                },
                "movie": {
                    "description": "Movie details returned by TMDB.",
                    "allOf": [
                        {
                            "$ref": "#/definitions/pkg.Movie"
                        }
                    ]
                },
                "source": {
                    "description": "Source filename.",
                    "type": "string"
                }
            }
        },
        "features.TVScannerResult": {
            "type": "object",
            "properties": {
                "destination": {
                    "description": "Full destination path of the moved file.",
                    "type": "string"
                },
                "source": {
                    "description": "Source filename.",
                    "type": "string"
                },
                "tvepisode": {
                    "description": "TV episode details returned by TMDB.",
                    "allOf": [
                        {
                            "$ref": "#/definitions/pkg.TVEpisode"
                        }
                    ]
                }
            }
        },
        "pkg.Category": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "pkg.Movie": {
            "type": "object",
            "properties": {
                "categories": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/pkg.Category"
                    }
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "releaseDate": {
                    "type": "string"
                }
            }
        },
        "pkg.TVEpisode": {
            "type": "object",
            "properties": {
                "categories": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/pkg.Category"
                    }
                },
                "episode": {
                    "type": "integer"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "releaseDate": {
                    "type": "string"
                },
                "season": {
                    "type": "integer"
                },
                "tvReleaseDate": {
                    "type": "string"
                },
                "tvShowID": {
                    "type": "integer"
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
