# NCN Lifecycle Management API
## NCN Lifecycle Events

A set of REST APIs that allow system admin to create ncn lifecycle event such as reboot/rebuild. It uses argo workflow to run defined procedures for each event.

## NCN Lifecycle Hooks [details](https://github.com/Cray-HPE/cray-nls/blob/master/docs/NCN%20Lifecycle%20Hooks.md)

A set of REST APIs that allow system admin or developers to register customize hooks as part of argo workflow.

## Security [details](https://github.com/Cray-HPE/cray-nls/blob/master/docs/security.md)

## Argo UI [details](https://github.com/Cray-HPE/cray-nls/blob/master/docs/argo.md)

## Version: 1.0

### /v1/liveness

#### GET
##### Summary

K8s Liveness endpoint

##### Responses

| Code | Description |
| ---- | ----------- |
| 204 |  |

### /v1/ncns/hooks

#### GET
##### Summary

Get ncn lifecycle hooks

##### Responses

| Code | Description |
| ---- | ----------- |
| 501 | Not Implemented |

#### POST
##### Summary

Get ncn lifecycle hooks

##### Responses

| Code | Description |
| ---- | ----------- |
| 501 | Not Implemented |

### /v1/ncns/hooks/{hook_id}

#### DELETE
##### Summary

Get ncn lifecycle hooks

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| hook_id | path | id of a hook | Yes | string |

##### Responses

| Code | Description |
| ---- | ----------- |
| 501 | Not Implemented |

#### PUT
##### Summary

Update a ncn lifecycle hook

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| hook_id | path | id of a hook | Yes | string |

##### Responses

| Code | Description |
| ---- | ----------- |
| 501 | Not Implemented |

### /v1/ncns/reboot

#### POST
##### Summary

End to end rolling reboot ncns

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| include | body | hostnames to include | Yes | [models.CreateRebootWorkflowRequest](#modelscreaterebootworkflowrequest) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [models.CreateRebootWorkflowResponse](#modelscreaterebootworkflowresponse) |
| 400 | Bad Request | [ResponseError](#responseerror) |
| 404 | Not Found | [ResponseError](#responseerror) |
| 500 | Internal Server Error | [ResponseError](#responseerror) |

### /v1/ncns/rebuild

#### POST
##### Summary

End to end rolling rebuild ncns

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| include | body | hostnames to include | Yes | [models.CreateRebuildWorkflowRequest](#modelscreaterebuildworkflowrequest) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [models.CreateRebuildWorkflowResponse](#modelscreaterebuildworkflowresponse) |
| 400 | Bad Request | [ResponseError](#responseerror) |
| 404 | Not Found | [ResponseError](#responseerror) |
| 500 | Internal Server Error | [ResponseError](#responseerror) |

### /v1/readiness

#### GET
##### Summary

K8s Readiness endpoint

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 204 |  |  |
| 500 | Internal Server Error | [ResponseError](#responseerror) |

### /v1/version

#### GET
##### Summary

Get version of cray-nls service

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [ResponseOk](#responseok) |
| 500 | Internal Server Error | [ResponseError](#responseerror) |

### /v1/workflows

#### GET
##### Summary

Get status of a ncn workflow

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| labelSelector | query | Label Selector | No | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [ [models.GetWorkflowResponse](#modelsgetworkflowresponse) ] |
| 400 | Bad Request | [ResponseError](#responseerror) |
| 404 | Not Found | [ResponseError](#responseerror) |
| 500 | Internal Server Error | [ResponseError](#responseerror) |

### /v1/workflows/{name}

#### DELETE
##### Summary

Delete a ncn workflow

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| name | path | name of workflow | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [ResponseOk](#responseok) |
| 400 | Bad Request | [ResponseError](#responseerror) |
| 404 | Not Found | [ResponseError](#responseerror) |
| 500 | Internal Server Error | [ResponseError](#responseerror) |

### /v1/workflows/{name}/rerun

#### PUT
##### Summary

Rerun a workflow, all steps will run

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| name | path | name of workflow | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [ResponseOk](#responseok) |
| 400 | Bad Request | [ResponseError](#responseerror) |
| 404 | Not Found | [ResponseError](#responseerror) |
| 500 | Internal Server Error | [ResponseError](#responseerror) |

### /v1/workflows/{name}/retry

#### PUT
##### Summary

Retry a failed ncn workflow, skip passed steps

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| name | path | name of workflow | Yes | string |
| retryOptions | body | retry options | Yes | [models.RetryWorkflowRequestBody](#modelsretryworkflowrequestbody) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [ResponseOk](#responseok) |
| 400 | Bad Request | [ResponseError](#responseerror) |
| 404 | Not Found | [ResponseError](#responseerror) |
| 500 | Internal Server Error | [ResponseError](#responseerror) |

### Models

#### ResponseError

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| message | string |  | No |

#### ResponseOk

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| message | string |  | No |

#### models.CreateRebootWorkflowRequest

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| dryRun | boolean |  | No |
| hosts | [ string ] |  | No |
| switchPassword | string |  | No |
| wipeOsd | boolean |  | No |

#### models.CreateRebootWorkflowResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| name | string |  | No |
| targetNcns | [ string ] |  | No |

#### models.CreateRebuildWorkflowRequest

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| dryRun | boolean |  | No |
| hosts | [ string ] |  | No |
| switchPassword | string |  | No |
| wipeOsd | boolean |  | No |

#### models.CreateRebuildWorkflowResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| name | string |  | No |
| targetNcns | [ string ] |  | No |

#### models.GetWorkflowResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| label | object |  | No |
| name | string |  | No |
| status | object |  | No |

#### models.RetryWorkflowRequestBody

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| restartSuccessful | boolean |  | No |
| stepName | string |  | No |
