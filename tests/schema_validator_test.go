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
