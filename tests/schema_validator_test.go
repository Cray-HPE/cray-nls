package tests

import (
	"testing"

	sv "github.com/Cray-HPE/cray-nls/src/api/models/iuf/schemaValidator"
	"sigs.k8s.io/yaml"
)

// Unit testcase for Basic sanity - schema validation failure
func TestSchemaValidatorFalureCheck(t *testing.T) {
	data := []byte(`
---
iuf_version: ^0.5.0
`)
	var dataObject map[string]interface{}
	inputErr := yaml.Unmarshal(data, &dataObject)

	if inputErr != nil {
		t.Fatal("Test setup has issues, error details:", inputErr)
	}

	schemaFile := "schemas/iuf-manifest-schema.yaml"

	err := sv.Validate(dataObject, schemaFile)
	if err == nil {
		t.Fatal("Basic schema validation sanity failing, must return validation error which is not returning")
	}

}

// Unit testcase for Basic sanity - schema validation pass case
func TestSchemaValidatorPassCheck(t *testing.T) {
	data := []byte(`
---
iuf_version: ^0.5.0
name: cos
description: >
  The Cray Operating System (COS).
version: 2.5.97
hooks:
  deliver_product:
    post:
      script_path: hooks/deliver_product-posthook.sh
  post_install_service_check:
    pre:
      script_path: hooks/post_install_service_check-prehook.sh

content:
  docker:
  - path: docker

  s3:
  - path: data/s3/dummy_upload_1.txt
    bucket: dummy-bucket
    key: dummy-key

  - path: data/s3/dummy_upload_2.txt
    bucket: dummy-bucket
    key: dummy-key

  helm:
  - path: helm

  loftsman:
  - path: manifests/cos-services.yaml
    use_manifestgen: true
    deploy: true

  nexus_blob_stores:
    yaml_path: 'data/nexus-blobstores.yaml'

  nexus_repositories:
    yaml_path: 'data/np/nexus-repositories.yaml'

  rpms:
  - path: data/rpms/rpm_dummy_1
    repository_name: cos-2.5.97-sle-15sp4
    repository_type: raw
  
  - path: data/rpms/rpm_dummy_2
    repository_name: cos-2.5.97-net-sle-15sp4-shs-2.0
    repository_type: raw
  
  - path: data/rpms/rpm_dummy_3
    repository_name: cos-2.5.97-sle-15sp4-compute
    repository_type: raw

  vcs:
    path: data/vcs

  ims:
    content_dirs:
    - ims/recipes/x86_64
`)
	var dataObject map[string]interface{}
	inputErr := yaml.Unmarshal(data, &dataObject)

	if inputErr != nil {
		t.Fatal("Test setup has issues, error details:", inputErr)
	}

	schemaFile := "schemas/iuf-manifest-schema.yaml"

	err := sv.Validate(dataObject, schemaFile)
	if err != nil {
		t.Fatal("Basic schema validation sanity failing, error", err)
	}

}
