
### /iuf/v1/activities

#### GET
##### Summary

List IUF activities

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [ [Iuf.Activity](#iufactivity) ] |
| 501 | Not Implemented |  |

### /iuf/v1/activities/{id}

#### GET
##### Summary

Get an IUF activities

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path | activity id | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [Iuf.Activity](#iufactivity) |
| 501 | Not Implemented |  |

#### PATCH
##### Summary

Patch an IUF activities

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path | activity id | Yes | string |
| partial_activity | body | partial IUF activity | Yes | [Iuf.CreateOrPatchActivityRequest](#iufcreateorpatchactivityrequest) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [Iuf.Activity](#iufactivity) |
| 501 | Not Implemented |  |

### /iuf/v1/activity

#### POST
##### Summary

Create an IUF activity

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| activity | body | IUF activity | Yes | [Iuf.CreateOrPatchActivityRequest](#iufcreateorpatchactivityrequest) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [Iuf.Activity](#iufactivity) |
| 501 | Not Implemented |  |

### Models

#### Iuf.Activity

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| inputs | [IufSession.InputParams](#iufsessioninputparams) |  | No |
| name | string |  | No |
| products | [ [Iuf.Product](#iufproduct) ] |  | No |
| sessions | [ object ] |  | No |

#### Iuf.CreateOrPatchActivityRequest

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| inputs | [IufSession.InputParams](#iufsessioninputparams) |  | No |
| products | [ [Iuf.Product](#iufproduct) ] |  | No |

#### Iuf.Product

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| name | string |  | No |
| version | string |  | No |

#### IufSession.InputParams

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| vcs_working_branch_pattern | string | The pattern to use for all products. Use the following variables in braces {} to specify the pattern:  {product_name} {product_version}  E.g.  {product_name}-{product_version}-test-branch | No |
| vcs_working_branch_per_product | object | Specify the working branch name per product. This is an object where the key is the product name, and the value is the exact name (not a pattern) of the VCS branch for that product. | No |

#### IufSession.Product

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| after_hook_scripts | object | Any after hook scripts for this product. This is an object where the key is operation name, and value is the CR name of the hook script for this product.  Hook scripts are executed either before or after an execution of a operation. They are specified in each product's distribution file, as part of the iuf-manifest.yaml.  The hook scripts are initially taken from the product distribution file and stored in S3, so that they can later be referenced. | No |
| before_hook_scripts | object | Any before hook scripts for this product. This is an object where the key is operation name, and value is the CR name of the hook script for this product.  Hook scripts are executed either before or after an execution of a operation. They are specified in each product's distribution file, as part of the iuf-manifest.yaml.  The hook scripts are initially taken from the product distribution file and stored in S3, so that they can later be referenced. | No |
| name | string | The name of the product | No |
| original_location | string | The original location of the extracted tar in on the physical storage. | No |
| version | string | The version of the product. | No |

#### IufSession.Spec

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| input_params | [IufSession.InputParams](#iufsessioninputparams) |  | No |
| products | [ [IufSession.Product](#iufsessionproduct) ] | The products that need to be installed, as specified by the Admin. | No |
| stages | [ string ] | The stages that need to be executed. This is either explicitly specified by the Admin, or it is computed from the workflow type. An Stage is a group of Operations. Stages represent the overall workflow at a high-level, and executing a stage means executing a bunch of Operations in a predefined manner.  An Admin can specify the stages that must be executed for an install-upgrade workflow. And Product Developers can extend each stage with custom hook scripts that they would like to run before and after the stage's execution.  The high-level stages allow their configuration would revealing too many details to the consumers of IUF. if not specified, we apply all stages | No |
| workflow_type | string | What type of workflow are we executing? install or upgrade | No |

#### IufSession.Status

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| argo_workflow | string | The unique name of the Argo workflow that is created from all the input parameters above. | No |
| message | string |  | No |
| observedGeneration | integer |  | No |
| operations | [ [ string ] ] | A 2-level DAG of Operations derived from stages that would be executed for each of the products that are specified. This is not specified by the Admin -- it is computed from the list of stages above.  This is an array of array of CR names of Operations that are installed as part of IUF, and determined by the Stages supplied. | No |
| phase | string |  | No |
