// Package docs GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import "github.com/swaggo/swag"

const docTemplateNLS = `{
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
        "/nls/v1/ncns/reboot": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "NCN Lifecycle Events"
                ],
                "summary": "End to end rolling reboot ncns",
                "parameters": [
                    {
                        "description": "hostnames to include",
                        "name": "include",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.CreateRebootWorkflowRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.CreateRebootWorkflowResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/ResponseError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/ResponseError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/ResponseError"
                        }
                    }
                }
            }
        },
        "/nls/v1/ncns/rebuild": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "NCN Lifecycle Events"
                ],
                "summary": "End to end rolling rebuild ncns",
                "parameters": [
                    {
                        "description": "hostnames to include",
                        "name": "include",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.CreateRebuildWorkflowRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.CreateRebuildWorkflowResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/ResponseError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/ResponseError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/ResponseError"
                        }
                    }
                }
            }
        },
        "/nls/v1/workflows": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Workflow Management"
                ],
                "summary": "Get status of a ncn workflow",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Label Selector",
                        "name": "labelSelector",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.GetWorkflowResponse"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/ResponseError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/ResponseError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/ResponseError"
                        }
                    }
                }
            }
        },
        "/nls/v1/workflows/{name}": {
            "delete": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Workflow Management"
                ],
                "summary": "Delete a ncn workflow",
                "parameters": [
                    {
                        "type": "string",
                        "description": "name of workflow",
                        "name": "name",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/ResponseOk"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/ResponseError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/ResponseError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/ResponseError"
                        }
                    }
                }
            }
        },
        "/nls/v1/workflows/{name}/rerun": {
            "put": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Workflow Management"
                ],
                "summary": "Rerun a workflow, all steps will run",
                "parameters": [
                    {
                        "type": "string",
                        "description": "name of workflow",
                        "name": "name",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/ResponseOk"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/ResponseError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/ResponseError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/ResponseError"
                        }
                    }
                }
            }
        },
        "/nls/v1/workflows/{name}/retry": {
            "put": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Workflow Management"
                ],
                "summary": "Retry a failed ncn workflow, skip passed steps",
                "parameters": [
                    {
                        "type": "string",
                        "description": "name of workflow",
                        "name": "name",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "retry options",
                        "name": "retryOptions",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.RetryWorkflowRequestBody"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/ResponseOk"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/ResponseError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/ResponseError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/ResponseError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "ResponseError": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "ResponseOk": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "models.CreateRebootWorkflowRequest": {
            "type": "object",
            "properties": {
                "dryRun": {
                    "type": "boolean"
                },
                "hosts": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "wipeOsd": {
                    "type": "boolean"
                }
            }
        },
        "models.CreateRebootWorkflowResponse": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string"
                },
                "targetNcns": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "models.CreateRebuildWorkflowRequest": {
            "type": "object",
            "properties": {
                "bootTimeoutInSeconds": {
                    "type": "integer"
                },
                "desiredCfsConfig": {
                    "type": "string"
                },
                "dryRun": {
                    "type": "boolean"
                },
                "hosts": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "imageId": {
                    "type": "string"
                },
                "labels": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "string"
                    }
                },
                "workflowType": {
                    "description": "used to determine storage rebuild vs upgrade",
                    "type": "string"
                },
                "zapOsds": {
                    "description": "this is necessary for storage rebuilds when unable to wipe the node prior to rebuild",
                    "type": "boolean"
                }
            }
        },
        "models.CreateRebuildWorkflowResponse": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string"
                },
                "targetNcns": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "models.GetWorkflowResponse": {
            "type": "object",
            "properties": {
                "label": {
                    "type": "object"
                },
                "name": {
                    "type": "string"
                },
                "status": {
                    "type": "object"
                }
            }
        },
        "models.RetryWorkflowRequestBody": {
            "type": "object",
            "properties": {
                "restartSuccessful": {
                    "type": "boolean"
                },
                "stepName": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfoNLS holds exported Swagger Info so clients can modify it
var SwaggerInfoNLS = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "/apis",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "NLS",
	SwaggerTemplate:  docTemplateNLS,
}

func init() {
	swag.Register(SwaggerInfoNLS.InstanceName(), SwaggerInfoNLS)
}
