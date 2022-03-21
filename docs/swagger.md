# NCN Lifecycle Management API
# (WIP)

This doc descibes REST API for ncn lifecycle management. Note that in this version, we only provide APIs for individual operation. A full end to end lifecycle management API is out of scope in Phase I

> TIP: This is Descrption is rendered from `docs/api.md`

## Version: 1.0

**License:** [Apache 2.0](http://www.apache.org/licenses/LICENSE-2.0.html)

### /etcd/{hostname}/prepare

#### PUT
##### Summary

Prepare baremetal etcd for a master node to rejoin

##### Description

# Prepare baremetal ETCD for rejoining

Prepare a master ncn to rejoin baremetal etcd cluster

## Pre-condition

1. **NCN** is a **master** node
1. Baremetal etcd cluster is in **healthy** state
1. quorum after removal

## Action

1. Remove a ncn from baremetal etcd cluster
1. Stop etcd services on the ncn
1. Add the ncn back to etcd cluster so it can rejoin on boot

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| hostname | path | Hostname of target ncn | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | ok | string |
| 400 | Bad Request | [utils.ResponseError](#utilsresponseerror) |
| 404 | Not Found | [utils.ResponseError](#utilsresponseerror) |
| 500 | Internal Server Error | [utils.ResponseError](#utilsresponseerror) |

### /kubernetes/{hostname}/drain

#### POST
##### Summary

Drain a Kubernetes node

##### Description

# Drain Kubernetes Node

Before we can safely drain/remove a node from k8s cluster, we need to run some `CSM specific logic` to make sure a node can be drained from k8s cluster safely

---

## Master

#### Pre-condition

1. **NCN** is a **master** node
1. quorum after removal

#### Actions

1. drain node

---

## Worker

#### Pre-condition

#### Actions

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| hostname | path | Hostname | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 400 | Bad Request | [utils.ResponseError](#utilsresponseerror) |
| 404 | Not Found | [utils.ResponseError](#utilsresponseerror) |
| 500 | Internal Server Error | [utils.ResponseError](#utilsresponseerror) |

### /kubernetes/{hostname}/move-first-master

#### POST
##### Summary

Move first master from a master k8s node

##### Description

# Move First Master

We need to make sure first master is not the node being rebuit. We need to move `first_master` to a different master node

### Pre-condition

1. **NCN** is a **master** node
1. **NCN** is already the **first master**

### Action

1. Loop through other master nodes until `scripts/k8s/promote-initial-master.sh` returns 0
2. Update `meta-data.first-master-hostname`

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| hostname | path | Hostname | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 400 | Bad Request | [utils.ResponseError](#utilsresponseerror) |
| 404 | Not Found | [utils.ResponseError](#utilsresponseerror) |
| 500 | Internal Server Error | [utils.ResponseError](#utilsresponseerror) |

### /kubernetes/{hostname}/post-rebuild

#### POST
##### Summary

Kubernetes node post rebuild action

##### Description

# Post Rebuild

After a node rejoined k8s cluster after rebuild, certain `CSM specific steps` are required. We need to perform such action so we put a system back up health state.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| hostname | path | Hostname | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 400 | Bad Request | [utils.ResponseError](#utilsresponseerror) |
| 404 | Not Found | [utils.ResponseError](#utilsresponseerror) |
| 500 | Internal Server Error | [utils.ResponseError](#utilsresponseerror) |

### /ncn/{hostname}/backup

#### POST
##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| hostname | path | Hostname | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 400 | Bad Request | [utils.ResponseError](#utilsresponseerror) |
| 404 | Not Found | [utils.ResponseError](#utilsresponseerror) |
| 500 | Internal Server Error | [utils.ResponseError](#utilsresponseerror) |

### /ncn/{hostname}/post-rebuild

#### POST
##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| hostname | path | Hostname | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 400 | Bad Request | [utils.ResponseError](#utilsresponseerror) |
| 404 | Not Found | [utils.ResponseError](#utilsresponseerror) |
| 500 | Internal Server Error | [utils.ResponseError](#utilsresponseerror) |

### /ncn/{hostname}/reboot

#### POST
##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| hostname | path | Hostname | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 400 | Bad Request | [utils.ResponseError](#utilsresponseerror) |
| 404 | Not Found | [utils.ResponseError](#utilsresponseerror) |
| 500 | Internal Server Error | [utils.ResponseError](#utilsresponseerror) |

### /ncn/{hostname}/restore

#### POST
##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| hostname | path | Hostname | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 400 | Bad Request | [utils.ResponseError](#utilsresponseerror) |
| 404 | Not Found | [utils.ResponseError](#utilsresponseerror) |
| 500 | Internal Server Error | [utils.ResponseError](#utilsresponseerror) |

### /ncn/{hostname}/validate

#### POST
##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| hostname | path | Hostname | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 400 | Bad Request | [utils.ResponseError](#utilsresponseerror) |
| 404 | Not Found | [utils.ResponseError](#utilsresponseerror) |
| 500 | Internal Server Error | [utils.ResponseError](#utilsresponseerror) |

### /ncn/{hostname}/wipe

#### POST
##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| hostname | path | Hostname | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 400 | Bad Request | [utils.ResponseError](#utilsresponseerror) |
| 404 | Not Found | [utils.ResponseError](#utilsresponseerror) |
| 500 | Internal Server Error | [utils.ResponseError](#utilsresponseerror) |

### Models

#### utils.ResponseError

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| message | string |  | No |
