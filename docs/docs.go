// Package docs GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/etcd/{hostname}/prepare": {
            "put": {
                "description": "## Prepare baremetal ETCD for rejoining\n\nPrepare a master ncn to rejoin baremetal etcd cluster\n\n### Pre-condition\n\n1. **NCN** is a **master** node\n1. Baremetal etcd cluster is in **healthy** state\n\n### Action\n\n1. Remove a ncn from baremetal etcd cluster\n1. Stop etcd services on the ncn\n1. Add the ncn back to etcd cluster so it can rejoin on boot\n",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Etcd"
                ],
                "summary": "Prepare baremetal etcd for a master node to rejoin",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Hostname of target ncn",
                        "name": "hostname",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.ResponseError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.ResponseError"
                        }
                    }
                }
            }
        },
        "/kubernetes/first-master": {
            "put": {
                "description": "## Move First Master\n\nWe need to make sure first master is not the node being rebuit. We need to move ` + "`" + `first_master` + "`" + ` to a different master node\n",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Kubernetes"
                ],
                "summary": "Move first master to a master k8s",
                "parameters": [
                    {
                        "description": "Hostname of target first master",
                        "name": "hostname",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    }
                ],
                "responses": {
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.ResponseError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/utils.ResponseError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.ResponseError"
                        }
                    }
                }
            }
        },
        "/kubernetes/{hostname}/drain": {
            "post": {
                "description": "## Drain Kubernetes Node\n\nBefore we can safely drain/remove a node from k8s cluster, we need to run some ` + "`" + `CSM specific logic` + "`" + ` to make sure a node can be drained from k8s cluster safely\n",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Kubernetes"
                ],
                "summary": "Drain a Kubernetes node",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Hostname",
                        "name": "hostname",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.ResponseError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/utils.ResponseError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.ResponseError"
                        }
                    }
                }
            }
        },
        "/kubernetes/{hostname}/post-rebuild": {
            "post": {
                "description": "## Post Rebuild\n\nAfter a node rejoined k8s cluster after rebuild, certain ` + "`" + `CSM specific steps` + "`" + ` are required. We need to perform such action so we put a system back up health state.\n",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Kubernetes"
                ],
                "summary": "Kubernetes node post rebuild action",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Hostname",
                        "name": "hostname",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.ResponseError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/utils.ResponseError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.ResponseError"
                        }
                    }
                }
            }
        },
        "/ncn/{hostname}/backup": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "NCN"
                ],
                "parameters": [
                    {
                        "type": "string",
                        "description": "Hostname",
                        "name": "hostname",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.ResponseError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/utils.ResponseError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.ResponseError"
                        }
                    }
                }
            }
        },
        "/ncn/{hostname}/post-rebuild": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "NCN"
                ],
                "parameters": [
                    {
                        "type": "string",
                        "description": "Hostname",
                        "name": "hostname",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.ResponseError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/utils.ResponseError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.ResponseError"
                        }
                    }
                }
            }
        },
        "/ncn/{hostname}/reboot": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "NCN"
                ],
                "parameters": [
                    {
                        "type": "string",
                        "description": "Hostname",
                        "name": "hostname",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.ResponseError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/utils.ResponseError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.ResponseError"
                        }
                    }
                }
            }
        },
        "/ncn/{hostname}/restore": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "NCN"
                ],
                "parameters": [
                    {
                        "type": "string",
                        "description": "Hostname",
                        "name": "hostname",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.ResponseError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/utils.ResponseError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.ResponseError"
                        }
                    }
                }
            }
        },
        "/ncn/{hostname}/validate": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "NCN"
                ],
                "parameters": [
                    {
                        "type": "string",
                        "description": "Hostname",
                        "name": "hostname",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.ResponseError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/utils.ResponseError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.ResponseError"
                        }
                    }
                }
            }
        },
        "/ncn/{hostname}/wipe": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "NCN"
                ],
                "parameters": [
                    {
                        "type": "string",
                        "description": "Hostname",
                        "name": "hostname",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.ResponseError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/utils.ResponseError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.ResponseError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "utils.ResponseError": {
            "type": "object",
            "properties": {
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
	Host:             "",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "NCN Lifecycle Management API",
	Description:      "# (WIP)\n\nThis doc descibes REST API for ncn lifecycle management. Note that in this version, we only provide APIs for individual operation. A full end to end lifecycle management API is out of scope in Phase I\n\n\n> TIP: This is Descrption is rendered from `docs/api.md`\n\n",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
