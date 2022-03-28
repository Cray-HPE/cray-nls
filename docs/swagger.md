# NCN Lifecycle Management API
## TODO

we need some doc here @Alex

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

### /v2/ncns/hooks

#### GET
##### Summary

Add additional steps after a ncn boot(reboot)

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| filter | query | filter | No | string |

##### Responses

| Code | Description |
| ---- | ----------- |
| 501 | Not Implemented |

##### Security

| Security Schema | Scopes |
| --- | --- |
| OAuth2Application | admin |

### /v2/ncns/hooks/{hook_name}

#### DELETE
##### Summary

Remove a ncn

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| hook_name | path | hook_name | Yes | string |

##### Responses

| Code | Description |
| ---- | ----------- |
| 501 | Not Implemented |

##### Security

| Security Schema | Scopes |
| --- | --- |
| OAuth2Application | admin |

### /v2/ncns/hooks/before-k8s-drain

#### POST
##### Summary

Add additional steps before k8s drain

##### Responses

| Code | Description |
| ---- | ----------- |
| 501 | Not Implemented |

##### Security

| Security Schema | Scopes |
| --- | --- |
| OAuth2Application | admin |

### /v2/ncns/hooks/before-wipe

#### POST
##### Summary

Add additional steps before wipe a ncn

##### Responses

| Code | Description |
| ---- | ----------- |
| 501 | Not Implemented |

##### Security

| Security Schema | Scopes |
| --- | --- |
| OAuth2Application | admin |

### /v2/ncns/hooks/post-boot

#### POST
##### Summary

Add additional steps after a ncn boot(reboot)

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

End to end rolling reboot ncns

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

End to end rolling rebuild ncns

##### Responses

| Code | Description |
| ---- | ----------- |
| 501 | Not Implemented |

##### Security

| Security Schema | Scopes |
| --- | --- |
| OAuth2Application | admin |
