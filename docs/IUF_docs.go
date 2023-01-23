/*
 *
 *  MIT License
 *
 *  (C) Copyright 2023 Hewlett Packard Enterprise Development LP
 *
 *  Permission is hereby granted, free of charge, to any person obtaining a
 *  copy of this software and associated documentation files (the "Software"),
 *  to deal in the Software without restriction, including without limitation
 *  the rights to use, copy, modify, merge, publish, distribute, sublicense,
 *  and/or sell copies of the Software, and to permit persons to whom the
 *  Software is furnished to do so, subject to the following conditions:
 *
 *  The above copyright notice and this permission notice shall be included
 *  in all copies or substantial portions of the Software.
 *
 *  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *  IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *  FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
 *  THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR
 *  OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
 *  ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
 *  OTHER DEALINGS IN THE SOFTWARE.
 *
 */
// Package docs GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import "github.com/swaggo/swag"

const docTemplateIUF = `{
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
        "/iuf/v1/activities": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Activities"
                ],
                "summary": "List IUF activities",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/iuf.Activity"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/ResponseError"
                        }
                    }
                }
            },
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Activities"
                ],
                "summary": "Create an IUF activity",
                "parameters": [
                    {
                        "description": "IUF activity",
                        "name": "activity",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/iuf.CreateActivityRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/iuf.Activity"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
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
        "/iuf/v1/activities/{activity_name}": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Activities"
                ],
                "summary": "Get an IUF activity",
                "parameters": [
                    {
                        "type": "string",
                        "description": "activity name",
                        "name": "activity_name",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/iuf.Activity"
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
            },
            "patch": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Activities"
                ],
                "summary": "Patches an existing IUF activity",
                "parameters": [
                    {
                        "type": "string",
                        "description": "activity name",
                        "name": "activity_name",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/iuf.Activity"
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
        "/iuf/v1/activities/{activity_name}/history": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "History"
                ],
                "summary": "List history of an iuf activity",
                "parameters": [
                    {
                        "type": "string",
                        "description": "activity name",
                        "name": "activity_name",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/iuf.History"
                            }
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
        "/iuf/v1/activities/{activity_name}/history/abort": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "History"
                ],
                "summary": "Abort a session",
                "parameters": [
                    {
                        "type": "string",
                        "description": "activity name",
                        "name": "activity_name",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Action Request",
                        "name": "action_request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/iuf.HistoryActionRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created"
                    },
                    "501": {
                        "description": "Not Implemented"
                    }
                }
            }
        },
        "/iuf/v1/activities/{activity_name}/history/blocked": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "History"
                ],
                "summary": "Mark a session blocked",
                "parameters": [
                    {
                        "type": "string",
                        "description": "activity name",
                        "name": "activity_name",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Action Request",
                        "name": "action_request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/iuf.HistoryActionRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created"
                    },
                    "501": {
                        "description": "Not Implemented"
                    }
                }
            }
        },
        "/iuf/v1/activities/{activity_name}/history/paused": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "History"
                ],
                "summary": "Pause a session",
                "parameters": [
                    {
                        "type": "string",
                        "description": "activity name",
                        "name": "activity_name",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Action Request",
                        "name": "action_request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/iuf.HistoryActionRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created"
                    },
                    "501": {
                        "description": "Not Implemented"
                    }
                }
            }
        },
        "/iuf/v1/activities/{activity_name}/history/resume": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "History"
                ],
                "summary": "Resume an activity",
                "parameters": [
                    {
                        "type": "string",
                        "description": "activity name",
                        "name": "activity_name",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Action Request",
                        "name": "action_request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/iuf.HistoryActionRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created"
                    },
                    "501": {
                        "description": "Not Implemented"
                    }
                }
            }
        },
        "/iuf/v1/activities/{activity_name}/history/run": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "History"
                ],
                "summary": "Run a session",
                "parameters": [
                    {
                        "type": "string",
                        "description": "activity name",
                        "name": "activity_name",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Action Request",
                        "name": "action_request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/iuf.HistoryRunActionRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/iuf.Session"
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
        "/iuf/v1/activities/{activity_name}/history/{start_time}": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "History"
                ],
                "summary": "Get a history item of an iuf activity",
                "parameters": [
                    {
                        "type": "string",
                        "description": "start time of a history item",
                        "name": "start_time",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/iuf.History"
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
            },
            "patch": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "History"
                ],
                "summary": "replace comment of a history item of an iuf activity",
                "parameters": [
                    {
                        "type": "string",
                        "description": "activity name",
                        "name": "activity_name",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "start time of a history item",
                        "name": "start_time",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Modify comment of a history",
                        "name": "activity",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/iuf.ReplaceHistoryCommentRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/iuf.History"
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
        "/iuf/v1/activities/{activity_name}/sessions": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Sessions"
                ],
                "summary": "List sessions of an IUF activity",
                "parameters": [
                    {
                        "type": "string",
                        "description": "activity name",
                        "name": "activity_name",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/iuf.Session"
                            }
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
        "/iuf/v1/activities/{activity_name}/sessions/{session_name}": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Sessions"
                ],
                "summary": "Get a session of an IUF activity",
                "parameters": [
                    {
                        "type": "string",
                        "description": "activity name",
                        "name": "activity_name",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "session name",
                        "name": "session_name",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/iuf.Session"
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
        "/iuf/v1/stages": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Stages"
                ],
                "summary": "Get the IUF stages",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/iuf.Stages"
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
        "iuf.Activity": {
            "type": "object",
            "required": [
                "activity_state",
                "input_parameters",
                "operation_outputs",
                "products",
                "site_parameters"
            ],
            "properties": {
                "activity_state": {
                    "description": "State of activity",
                    "enum": [
                        "paused",
                        "in_progress",
                        "debug",
                        "blocked",
                        "wait_for_admin"
                    ],
                    "allOf": [
                        {
                            "$ref": "#/definitions/iuf.ActivityState"
                        }
                    ]
                },
                "input_parameters": {
                    "description": "Input parameters by admin",
                    "allOf": [
                        {
                            "$ref": "#/definitions/iuf.InputParameters"
                        }
                    ]
                },
                "name": {
                    "description": "Name of activity",
                    "type": "string"
                },
                "operation_outputs": {
                    "description": "Operation outputs from argo",
                    "type": "object",
                    "additionalProperties": true
                },
                "products": {
                    "description": "List of products included in an activity",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/iuf.Product"
                    }
                },
                "site_parameters": {
                    "description": "Site parameters set by the admin",
                    "allOf": [
                        {
                            "$ref": "#/definitions/iuf.SiteParameters"
                        }
                    ]
                }
            }
        },
        "iuf.ActivityState": {
            "type": "string",
            "enum": [
                "in_progress",
                "paused",
                "debug",
                "blocked",
                "wait_for_admin"
            ],
            "x-enum-varnames": [
                "ActivityStateInProgress",
                "ActivityStatePaused",
                "ActivityStateDebug",
                "ActivityStateBlocked",
                "ActivityStateWaitForAdmin"
            ]
        },
        "iuf.CreateActivityRequest": {
            "type": "object",
            "required": [
                "name"
            ],
            "properties": {
                "name": {
                    "description": "Name of activity",
                    "type": "string"
                }
            }
        },
        "iuf.History": {
            "type": "object",
            "required": [
                "activity_state"
            ],
            "properties": {
                "activity_state": {
                    "description": "State of activity",
                    "enum": [
                        "paused",
                        "in_progress",
                        "debug",
                        "blocked",
                        "wait_for_admin"
                    ],
                    "allOf": [
                        {
                            "$ref": "#/definitions/iuf.ActivityState"
                        }
                    ]
                },
                "comment": {
                    "description": "Comment",
                    "type": "string"
                },
                "name": {
                    "description": "Comment",
                    "type": "string"
                },
                "session_name": {
                    "description": "Name of the session",
                    "type": "string"
                },
                "start_time": {
                    "description": "Epoch timestamp",
                    "type": "integer"
                }
            }
        },
        "iuf.HistoryActionRequest": {
            "type": "object",
            "properties": {
                "comment": {
                    "description": "Comment",
                    "type": "string"
                },
                "start_time": {
                    "description": "Epoch timestamp",
                    "type": "integer"
                }
            }
        },
        "iuf.HistoryRunActionRequest": {
            "type": "object",
            "required": [
                "input_parameters"
            ],
            "properties": {
                "comment": {
                    "description": "Comment",
                    "type": "string"
                },
                "input_parameters": {
                    "$ref": "#/definitions/iuf.InputParameters"
                },
                "site_parameters": {
                    "$ref": "#/definitions/iuf.SiteParameters"
                }
            }
        },
        "iuf.InputParameters": {
            "type": "object",
            "properties": {
                "bootprep_config_managed": {
                    "description": "The path to the bootprep config file for managed nodes, relative to the media_dir",
                    "type": "string"
                },
                "bootprep_config_management": {
                    "description": "The path to the bootprep config file for management nodes, relative to the media_dir",
                    "type": "string"
                },
                "concurrency": {
                    "description": "An integer defining how many products / operations can we concurrently execute.",
                    "type": "integer"
                },
                "force": {
                    "description": "Force re-execution of stage operations",
                    "type": "boolean"
                },
                "limit_managed_nodes": {
                    "description": "Each item is the xname of a managed node",
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "limit_management_nodes": {
                    "description": "Each item is the xname of a management node",
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "media_dir": {
                    "description": "Location of media",
                    "type": "string"
                },
                "media_host": {
                    "description": "A string containing the hostname of where the media is located",
                    "type": "string"
                },
                "site_parameters": {
                    "description": "DEPRECATED: use site_parameters at the top level of the activity or session resource. The inline contents of the site_parameters.yaml file.",
                    "type": "string"
                },
                "stages": {
                    "description": "Stages to execute",
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "iuf.Operations": {
            "type": "object",
            "required": [
                "name",
                "static-parameters"
            ],
            "properties": {
                "name": {
                    "description": "Name of the operation",
                    "type": "string"
                },
                "static-parameters": {
                    "type": "object",
                    "additionalProperties": true
                }
            }
        },
        "iuf.Product": {
            "type": "object",
            "required": [
                "name",
                "original_location",
                "validated",
                "version"
            ],
            "properties": {
                "manifest": {
                    "description": "the content of manifest",
                    "type": "string"
                },
                "name": {
                    "description": "The name of the product",
                    "type": "string"
                },
                "original_location": {
                    "description": "The original location of the extracted tar in on the physical storage.",
                    "type": "string"
                },
                "validated": {
                    "description": "The flag indicates md5 of a product tarball file has been validated",
                    "type": "boolean"
                },
                "version": {
                    "description": "The version of the product.",
                    "type": "string"
                }
            }
        },
        "iuf.ReplaceHistoryCommentRequest": {
            "type": "object",
            "properties": {
                "comment": {
                    "description": "Comment",
                    "type": "string"
                }
            }
        },
        "iuf.Session": {
            "type": "object",
            "required": [
                "products"
            ],
            "properties": {
                "current_state": {
                    "enum": [
                        "paused",
                        "in_progress",
                        "debug",
                        "completed"
                    ],
                    "allOf": [
                        {
                            "$ref": "#/definitions/iuf.SessionState"
                        }
                    ]
                },
                "input_parameters": {
                    "$ref": "#/definitions/iuf.InputParameters"
                },
                "name": {
                    "type": "string"
                },
                "products": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/iuf.Product"
                    }
                },
                "site_parameters": {
                    "$ref": "#/definitions/iuf.SiteParameters"
                },
                "stage": {
                    "type": "string"
                },
                "workflows": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/iuf.SessionWorkflow"
                    }
                }
            }
        },
        "iuf.SessionState": {
            "type": "string",
            "enum": [
                "in_progress",
                "paused",
                "debug",
                "completed"
            ],
            "x-enum-varnames": [
                "SessionStateInProgress",
                "SessionStatePaused",
                "SessionStateDebug",
                "SessionStateCompleted"
            ]
        },
        "iuf.SessionWorkflow": {
            "type": "object",
            "properties": {
                "id": {
                    "description": "id of argo workflow",
                    "type": "string"
                },
                "url": {
                    "description": "url to the argo workflow",
                    "type": "string"
                }
            }
        },
        "iuf.SiteParameters": {
            "type": "object",
            "properties": {
                "global": {
                    "description": "global parameters applicable to all products",
                    "type": "object",
                    "additionalProperties": true
                },
                "products": {
                    "description": "Product-specific parameters",
                    "type": "object",
                    "additionalProperties": {
                        "type": "object",
                        "additionalProperties": true
                    }
                }
            }
        },
        "iuf.Stage": {
            "type": "object",
            "required": [
                "name",
                "operations",
                "type"
            ],
            "properties": {
                "name": {
                    "description": "Name of the stage",
                    "type": "string"
                },
                "no-hooks": {
                    "description": "no-hook indicates that there are no hooks that should be run for this stage",
                    "type": "boolean"
                },
                "operations": {
                    "description": "operations",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/iuf.Operations"
                    }
                },
                "type": {
                    "description": "Type of the stage",
                    "type": "string"
                }
            }
        },
        "iuf.Stages": {
            "type": "object",
            "required": [
                "stages",
                "version"
            ],
            "properties": {
                "hooks": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "string"
                    }
                },
                "stages": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/iuf.Stage"
                    }
                },
                "version": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfoIUF holds exported Swagger Info so clients can modify it
var SwaggerInfoIUF = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "/apis",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "IUF",
	SwaggerTemplate:  docTemplateIUF,
}

func init() {
	swag.Register(SwaggerInfoIUF.InstanceName(), SwaggerInfoIUF)
}
