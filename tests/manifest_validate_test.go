package tests

import (
	"testing"

	iuf "github.com/Cray-HPE/cray-nls/src/api/models/iuf"
	mdv "github.com/Cray-HPE/cray-nls/src/api/models/iuf/manifestDataValidation"

	"sigs.k8s.io/yaml"
)

// Unit testcase for Basic sanity of manifest validate
func TestValidateManifestSaninty(t *testing.T) {
	file := "data/manifests/iuf-product-manifest.yaml"
	err := iuf.ValidateFile(file)
	if err != nil {
		t.Fatal("Issue seen in manifest file:", file, "error:", err)
	}
}

// Unit testcase for checking invalid 3s path
func TestValidateManifestInvalidS3(t *testing.T) {

	data := []byte(`
content:
  docker:
  - path: docker	
  s3:
  - path: data/s3/dummy_upload_1.txt
    bucket: dummy-bucket
    key: dummy-key
  - path: not/exists/path/dummy_upload_2.txt
    bucket: dummy-bucket
    key: dummy-key
`)
	var dataObject map[string]interface{}
	inputErr := yaml.Unmarshal(data, &dataObject)

	if inputErr != nil {
		t.Fatal("Test setup has issues, error details:", inputErr)
	}

	err := mdv.Validate(dataObject)

	if err == nil {
		t.Fatal("s3 path validation logic is not working properly")
	}
}

// Unit testcase for checking invalid nexus blob path
func TestValidateManifestInvalidNexusBlob(t *testing.T) {

	data := []byte(`
content:
  nexus_blob_stores:
    yaml_path: 'do/not/exist/nexus-blobstores.yaml'
`)
	var dataObject map[string]interface{}
	inputErr := yaml.Unmarshal(data, &dataObject)

	if inputErr != nil {
		t.Fatal("Test setup has issues, error details:", inputErr)
	}

	err := mdv.Validate(dataObject)

	if err == nil {
		t.Fatal("nexus blob validation logic is not working properly")
	}
}

// Unit testcase for checking empty vcs path
func TestValidateManifestEmptyVcs(t *testing.T) {

	data := []byte(`
content:
  vcs:
    path: 'data/vcs_empty'
`)
	var dataObject map[string]interface{}
	inputErr := yaml.Unmarshal(data, &dataObject)

	if inputErr != nil {
		t.Fatal("Test setup has issues, error details:", inputErr)
	}

	err := mdv.Validate(dataObject)

	if err == nil {
		t.Fatal("vcs validation logic is not working properly")
	}
}

// Unit testcase for checking invalid vcs path
func TestValidateManifestInvalidVcs(t *testing.T) {
	data := []byte(`
content:
  vcs:
    path: 'data/do/not/exist'
`)
	var dataObject map[string]interface{}
	inputErr := yaml.Unmarshal(data, &dataObject)

	if inputErr != nil {
		t.Fatal("Test setup has issues, error details:", inputErr)
	}

	err := mdv.Validate(dataObject)

	if err == nil {
		t.Fatal("vcs validation logic is not working properly")
	}
}

/*
	Test cases for rpm validations
*/
// checking empty rpm path
func TestValidateManifestEmptyRpm(t *testing.T) {
	data := []byte(`
content:
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

   - path: data/rpms/empty_rpm
     repository_name: cos-2.5.97-sle-15sp4-compute
     repository_type: raw
`)
	var dataObject map[string]interface{}
	inputErr := yaml.Unmarshal(data, &dataObject)

	if inputErr != nil {
		t.Fatal("Test setup has issues, error details:", inputErr)
	}

	err := mdv.Validate(dataObject)

	if err == nil {
		t.Fatal("rpm validation logic is not working properly")
	}
}

// checking invalid rpm path
func TestValidateManifestInvalidRpm(t *testing.T) {
	data := []byte(`
content:
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

   - path: data/rpms/wrong/path
     repository_name: cos-2.5.97-sle-15sp4-compute
     repository_type: raw
`)
	var dataObject map[string]interface{}
	inputErr := yaml.Unmarshal(data, &dataObject)

	if inputErr != nil {
		t.Fatal("Test setup has issues, error details:", inputErr)
	}

	err := mdv.Validate(dataObject)

	if err == nil {
		t.Fatal("rpm validation logic is not working properly")
	}
}

// Handling of rpm not present in nexus repo
func TestValidateHostedLogic(t *testing.T) {
	data := []byte(`
content:
  nexus_repositories:
    yaml_path: 'data/np/nexus-repositories.yaml'
  rpms:
   - path: data/rpms/rpm_dummy_1
     repository_name: cos-2.5.97-sle-15sp4
     repository_type: raw

   - path: data/rpms/rpm_dummy_2
     repository_name: repo_non_existent
     repository_type: raw

   - path: data/rpms/rpm_dummy_3
     repository_name: cos-2.5.97-sle-15sp4-compute
     repository_type: raw
`)
	var dataObject map[string]interface{}
	inputErr := yaml.Unmarshal(data, &dataObject)

	if inputErr != nil {
		t.Fatal("Test setup has issues, error details:", inputErr)
	}

	err := mdv.Validate(dataObject)

	if err == nil {
		t.Fatal("rpm validation logic is not working properly")
	}
}

// Handling of group rpm in manifest file
func TestValidateGroupRpmExpectionLogic(t *testing.T) {
	data := []byte(`
content:
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
     repository_name: cos-2.5-net-sle-15sp4-shs-2.0
     repository_type: raw
`) // last path is a group repo, defined in nexus-repositories.yaml

	var dataObject map[string]interface{}
	inputErr := yaml.Unmarshal(data, &dataObject)

	if inputErr != nil {
		t.Fatal("Test setup has issues, error details:", inputErr)
	}

	err := mdv.Validate(dataObject)

	if err == nil {
		t.Fatal("rpm validation logic is not working properly", err)
	}
}

/*
	Test cases for Nexus repo validations
*/
// checking invalid nexus repo path
func TestValidateManifestInvalidNexusRepo(t *testing.T) {
	data := []byte(`
content:
  nexus_repositories:
    yaml_path: 'data/np/dummy/filenotpresent.yaml'
`)
	var dataObject map[string]interface{}
	inputErr := yaml.Unmarshal(data, &dataObject)

	if inputErr != nil {
		t.Fatal("Test setup has issues, error details:", inputErr)
	}

	err := mdv.Validate(dataObject)

	if err == nil {
		t.Fatal("nexus repo validation logic is not working properly")
	}
}

// Handling of missing hosted
func TestValidateNexusRepoMissingHosted(t *testing.T) {
	data := []byte(`
content:
  nexus_repositories:
    yaml_path: 'data/np/nexus-repositories-dummy.yaml'
`)
	var dataObject map[string]interface{}
	inputErr := yaml.Unmarshal(data, &dataObject)

	if inputErr != nil {
		t.Fatal("Test setup has issues, error details:", inputErr)
	}

	nexusCotent := []byte(`
---
cleanup: null
format: raw
name: cos-2.5.97-sle-15sp4
online: true
storage:
  blobStoreName: cos
  strictContentTypeValidation: false
  writePolicy: ALLOW_ONCE
type: hosted
---
format: raw
group:
  memberNames:
  - cos-2.5.97-sle-15sp4
name: cos-2.5-sle-15sp4
online: true
storage:
  blobStoreName: cos
  strictContentTypeValidation: false
type: group
---
format: raw
group:
  memberNames:
  - cos-2.5.97-net-sle-15sp4-shs-2.0
name: cos-2.5-net-sle-15sp4-shs-2.0
online: true
storage:
  blobStoreName: cos
  strictContentTypeValidation: false
type: group
`)
	//mocking the file IO
	mdv.FileReader = func(filePath string) ([]byte, error) {
		return nexusCotent, nil
	}
	err := mdv.Validate(dataObject)

	if err == nil {
		t.Fatal("nexus repo validation logic is not working properly", err)
	}
}

// Handling of hosted defined below group
func TestValidateNexusRepoMisplacedHosted(t *testing.T) {
	data := []byte(`
content:
  nexus_repositories:
    yaml_path: 'data/np/nexus-repositories-dummy.yaml'
`)
	var dataObject map[string]interface{}
	inputErr := yaml.Unmarshal(data, &dataObject)

	if inputErr != nil {
		t.Fatal("Test setup has issues, error details:", inputErr)
	}

	nexusCotent := []byte(`
---
cleanup: null
format: raw
name: cos-2.5.97-sle-15sp4
online: true
storage:
  blobStoreName: cos
  strictContentTypeValidation: false
  writePolicy: ALLOW_ONCE
type: hosted
---
format: raw
group:
  memberNames:
  - cos-2.5.97-sle-15sp4
name: cos-2.5-sle-15sp4
online: true
storage:
  blobStoreName: cos
  strictContentTypeValidation: false
type: group
---
format: raw
group:
  memberNames:
  - cos-2.5.97-net-sle-15sp4-shs-2.0
name: cos-2.5-net-sle-15sp4-shs-2.0
online: true
storage:
  blobStoreName: cos
  strictContentTypeValidation: false
type: group
---
cleanup: null
format: raw
name: cos-2.5.97-net-sle-15sp4-shs-2.0
online: true
storage:
  blobStoreName: cos
  strictContentTypeValidation: false
  writePolicy: ALLOW_ONCE
type: hosted
`)
	//mocking the file IO
	mdv.FileReader = func(filePath string) ([]byte, error) {
		return nexusCotent, nil
	}
	err := mdv.Validate(dataObject)

	if err == nil {
		t.Fatal("nexus repo validation logic is not working properly")
	}
}

// Handling of format skipping logic
func TestValidateNexusRepoSkipFormat(t *testing.T) {
	data := []byte(`
content:
  nexus_repositories:
    yaml_path: 'data/np/nexus-repositories-dummy.yaml'
`)
	var dataObject map[string]interface{}
	inputErr := yaml.Unmarshal(data, &dataObject)

	if inputErr != nil {
		t.Fatal("Test setup has issues, error details:", inputErr)
	}

	nexusCotent := []byte(`
---
cleanup: null
format: raw
name: cos-2.5.97-sle-15sp4
online: true
storage:
  blobStoreName: cos
  strictContentTypeValidation: false
  writePolicy: ALLOW_ONCE
type: hosted
---
format: raw
group:
  memberNames:
  - cos-2.5.97-sle-15sp4
name: cos-2.5-sle-15sp4
online: true
storage:
  blobStoreName: cos
  strictContentTypeValidation: false
type: group
---
cleanup: null
format: helm
name: charts
online: true
storage:
  blobStoreName: csm
  strictContentTypeValidation: false
  writePolicy: ALLOW_ONCE
type: hosted
---
cleanup: null
format: raw
name: cos-2.5.97-net-sle-15sp4-shs-2.0
online: true
storage:
  blobStoreName: cos
  strictContentTypeValidation: false
  writePolicy: ALLOW_ONCE
type: hosted
---
cleanup: null
docker:
  forceBasicAuth: true
  httpPort: 5003
  httpsPort: null
  v1Enabled: false
format: docker
name: registry
online: true
storage:
  blobStoreName: csm
  strictContentTypeValidation: false
  writePolicy: ALLOW
type: hosted
---
format: raw
group:
  memberNames:
  - cos-2.5.97-net-sle-15sp4-shs-2.0
name: cos-2.5-net-sle-15sp4-shs-2.0
online: true
storage:
  blobStoreName: cos
  strictContentTypeValidation: false
type: group
`)
	//mocking the file IO
	mdv.FileReader = func(filePath string) ([]byte, error) {
		return nexusCotent, nil
	}
	err := mdv.Validate(dataObject)

	if err != nil {
		t.Fatal("nexus repo validation logic is not working properly")
	}
}

// Handling of multiple repo member logic logic -> + ve test case
func TestValidateNexusRepoMultiMemberPos(t *testing.T) {
	data := []byte(`
content:
  nexus_repositories:
    yaml_path: 'data/np/nexus-repositories-dummy.yaml'
`)
	var dataObject map[string]interface{}
	inputErr := yaml.Unmarshal(data, &dataObject)

	if inputErr != nil {
		t.Fatal("Test setup has issues, error details:", inputErr)
	}

	nexusCotent := []byte(`
---
cleanup: null
format: raw
name: cos-2.5.98-sle-15sp4
online: true
storage:
  blobStoreName: cos
  strictContentTypeValidation: false
  writePolicy: ALLOW_ONCE
type: hosted
---
cleanup: null
format: raw
name: cos-2.5.97-sle-15sp4
online: true
storage:
  blobStoreName: cos
  strictContentTypeValidation: false
  writePolicy: ALLOW_ONCE
type: hosted
---
format: raw
group:
  memberNames:
  - cos-2.5.97-sle-15sp4
  - cos-2.5.98-sle-15sp4
name: cos-2.5-sle-15sp4
online: true
storage:
  blobStoreName: cos
  strictContentTypeValidation: false
type: group
`)
	//mocking the file IO
	mdv.FileReader = func(filePath string) ([]byte, error) {
		return nexusCotent, nil
	}
	err := mdv.Validate(dataObject)

	if err != nil {
		t.Fatal("nexus repo validation logic is not working properly")
	}
}

// Handling of multiple repo member logic logic -> - ve test case
func TestValidateNexusRepoMultiMemberNeg(t *testing.T) {
	data := []byte(`
content:
  nexus_repositories:
    yaml_path: 'data/np/nexus-repositories-dummy.yaml'
`)
	var dataObject map[string]interface{}
	inputErr := yaml.Unmarshal(data, &dataObject)

	if inputErr != nil {
		t.Fatal("Test setup has issues, error details:", inputErr)
	}

	nexusCotent := []byte(`
---
cleanup: null
format: raw
name: cos-2.5.98-sle-15sp4
online: true
storage:
  blobStoreName: cos
  strictContentTypeValidation: false
  writePolicy: ALLOW_ONCE
type: hosted
---
format: raw
group:
  memberNames:
  - cos-2.5.97-sle-15sp4
  - cos-2.5.98-sle-15sp4
name: cos-2.5-sle-15sp4
online: true
storage:
  blobStoreName: cos
  strictContentTypeValidation: false
type: group
---
cleanup: null
format: raw
name: cos-2.5.97-sle-15sp4
online: true
storage:
  blobStoreName: cos
  strictContentTypeValidation: false
  writePolicy: ALLOW_ONCE
type: hosted
`)
	//mocking the file IO
	mdv.FileReader = func(filePath string) ([]byte, error) {
		return nexusCotent, nil
	}
	err := mdv.Validate(dataObject)

	if err == nil {
		t.Fatal("nexus repo validation logic is not working properly")
	}
}
