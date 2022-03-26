# NCN Lifecycle Management API
This doc descibes REST API for ncn lifecycle management. Note that in this version, we only provide APIs for individual operation. A full end to end lifecycle management API is out of scope in Phase I

---

## Argo workflow Demo

---

[API Doc](swagger.md)

## Version: 1.0

**License:** [Apache 2.0](http://www.apache.org/licenses/LICENSE-2.0.html)

### Security
**OAuth2Application**  

|oauth2|*OAuth 2.0*|
|---|---|
|Flow|application|
|**Scopes**||
|admin|                             Grants read and write access to administrative information|
|read|                              Grants read access|
|Token URL|<https://example.com/oauth/token>|

### /v1/ncns/{hostname}/reboot

#### POST
##### Summary

End to end reboot of a single ncn

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| hostname | path | hostname | Yes | string |

##### Responses

| Code | Description |
| ---- | ----------- |
| 501 | Not Implemented |

##### Security

| Security Schema | Scopes |
| --- | --- |
| OAuth2Application | admin |

### /v1/ncns/{hostname}/rebuild

#### POST
##### Summary

End to end rebuild of a single ncn

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| hostname | path | hostname | Yes | string |

##### Responses

| Code | Description |
| ---- | ----------- |
| 501 | Not Implemented |

##### Security

| Security Schema | Scopes |
| --- | --- |
| OAuth2Application | admin |

### /v1/workflows

#### GET
##### Summary

Get status of a ncn workflow

##### Responses

| Code | Description |
| ---- | ----------- |
| 501 | Not Implemented |

##### Security

| Security Schema | Scopes | |
| --- | --- | --- |
| OAuth2Application | admin | read |

### /v1/workflows/{name}

#### DELETE
##### Summary

Delete a ncn workflow

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| name | path | name of workflow | Yes | string |

##### Responses

| Code | Description |
| ---- | ----------- |
| 501 | Not Implemented |

##### Security

| Security Schema | Scopes | |
| --- | --- | --- |
| OAuth2Application | admin | read |

### /v1/workflows/{name}/rerun

#### PUT
##### Summary

Rerun a workflow, all steps will run

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| name | path | name of workflow | Yes | string |

##### Responses

| Code | Description |
| ---- | ----------- |
| 501 | Not Implemented |

##### Security

| Security Schema | Scopes | |
| --- | --- | --- |
| OAuth2Application | admin | read |

### /v1/workflows/{name}/retry

#### PUT
##### Summary

Retry a failed ncn workflow, skip passed steps

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| name | path | name of workflow | Yes | string |

##### Responses

| Code | Description |
| ---- | ----------- |
| 501 | Not Implemented |

##### Security

| Security Schema | Scopes | |
| --- | --- | --- |
| OAuth2Application | admin | read |

### /v2/ncn

#### POST
##### Summary

Add a ncn

##### Responses

| Code | Description |
| ---- | ----------- |
| 501 | Not Implemented |

##### Security

| Security Schema | Scopes |
| --- | --- |
| OAuth2Application | admin |

### /v2/ncns/{hostname}

#### DELETE
##### Summary

Remove a ncn

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| hostname | path | hostname | Yes | string |

##### Responses

| Code | Description |
| ---- | ----------- |
| 501 | Not Implemented |

##### Security

| Security Schema | Scopes |
| --- | --- |
| OAuth2Application | admin |

### /v2/ncns/reboot

#### POST
##### Summary

End to end rolling reboot request

##### Responses

| Code | Description |
| ---- | ----------- |
| 501 | Not Implemented |

##### Security

| Security Schema | Scopes |
| --- | --- |
| OAuth2Application | admin |

### /v2/ncns/rebuild

#### POST
##### Summary

End to end rolling rebuild request

##### Responses

| Code | Description |
| ---- | ----------- |
| 501 | Not Implemented |

##### Security

| Security Schema | Scopes |
| --- | --- |
| OAuth2Application | admin |
