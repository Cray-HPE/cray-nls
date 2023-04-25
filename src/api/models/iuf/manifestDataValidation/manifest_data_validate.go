package manifestDataValidation

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	mutils "github.com/Cray-HPE/cray-nls/src/api/models/iuf/mutils"

	//v2yaml "gopkg.in/yaml.v2"
	goyaml "github.com/go-yaml/yaml"
)

// const for yaml keys
const (
	CONTENT_KEY         string = "content"
	S3_KEY              string = "s3"
	S3_PATH_KEY         string = "path"
	NEXUS_REPO_KEY      string = "nexus_repositories"
	NEXUS_REPO_PATH_KEY string = "yaml_path"
	NEXUS_BLOB_KEY      string = "nexus_blob_stores"
	NEXUS_BLOB_PATH_KEY string = "yaml_path"
	VCS_KEY             string = "vcs"
	VCS_PATH_KEY        string = "path"
	RPM_KEY             string = "rpms"
	RPM_PATH_KEY        string = "path"
	REPO_NAME           string = "repository_name"
	HOSTED_REPO_NAME    string = "name"
	REPO_TYPE           string = "type"
	FORMAT              string = "format"
)

// Getting the file reader
var FileReader = mutils.ReadYamFile

// Struct with list of validators
type validators struct {
	content           map[string]interface{}
	nexusRepoFileName string
	hostedRepoNames   []string
}

// Method to process s3 content, returns error in case of issues
func (vs *validators) processS3() error {

	s3, s3_present := vs.content[S3_KEY] // extracting s3 key

	if !s3_present {
		return nil // s3 key missing
	}

	s3_array := s3.([]interface{}) // assuming array, is it validated in schema??
	for _, s3 := range s3_array {
		s3_element := s3.(map[string]interface{})
		file_path := s3_element[S3_PATH_KEY].(string)

		exist := mutils.IsPathExist(file_path) // do we need error details??
		if !exist {                            // if path is invalid
			return fmt.Errorf("error in processing s3 file %v", file_path)
		}
	}
	return nil
}

// Method to process nexus repo content, returns repo file path and error(in case of issues),
func (vs *validators) processNexusRepo() error {
	nr, nr_present := vs.content[NEXUS_REPO_KEY] // extracting nexus repo key

	if !nr_present {
		return nil // nexus repo key missing
	}

	nr_map := nr.(map[string]interface{}) // assuming map, is it validated in schema??
	file_path := nr_map[NEXUS_REPO_PATH_KEY].(string)
	vs.nexusRepoFileName = file_path

	exist := mutils.IsPathExist(file_path) // do we need error details??
	if !exist {                            // if path is invalid
		return fmt.Errorf("error in processing nexus repo file %v", file_path)
	}

	return nil
}

// Method to process nexus repo file and get hosted repo names
func (vs *validators) processNexusRepoFile() error {

	if vs.nexusRepoFileName == "" {
		return nil // skip processing of nexus repo file
	}

	nexusFile_contents, err := FileReader(vs.nexusRepoFileName)
	var temp_repo_names []string

	if err != nil {
		return fmt.Errorf("failed to open Nexus Repository file: %v", err)
	}

	dec := goyaml.NewDecoder(bytes.NewReader(nexusFile_contents))

	skipFormats := []string{"docker", "helm"}

	for {
		var nexusContent map[string]interface{}

		err := dec.Decode(&nexusContent)
		if errors.Is(err, io.EOF) {
			break
		}

		format := nexusContent[FORMAT].(string)

		isFormatToBeSkipped, _ := mutils.StringFoundInArray(skipFormats, format)

		if isFormatToBeSkipped {
			continue //Skip doc which has format that does not require validataion
		}

		if nexusContent[REPO_TYPE] == "hosted" {

			vs.hostedRepoNames = append(vs.hostedRepoNames, nexusContent["name"].(string))
			temp_repo_names = append(temp_repo_names, nexusContent["name"].(string))

		} else if nexusContent[REPO_TYPE] == "group" {
			group_map := nexusContent["group"].(map[interface{}]interface{})

			for _, v := range group_map {
				memNames_array := v.([]interface{})
				for _, m := range memNames_array {

					memberRepo := m.(string)

					isHostedRepo, index := mutils.StringFoundInArray(temp_repo_names, memberRepo)
					if isHostedRepo {
						temp_repo_names, err = mutils.Delete(temp_repo_names, index)

						if err != nil {
							fmt.Println(err)
						}

					} else {
						return fmt.Errorf("Repo referenced in group does not match hosted repo or Hosted Repos are not listed before group repos")
					}
				}
			}
		}
	}
	if len(temp_repo_names) > 0 {
		return fmt.Errorf("Repo defined in host repo is not listed in group repo")
	}

	return nil
}

// Method to process nexus blob content, returns error in case of issues
func (vs *validators) processNexusBlob() error {
	nb, nb_present := vs.content[NEXUS_BLOB_KEY] // extracting nexus blob key

	if !nb_present {
		return nil // nexus blob key missing
	}

	nb_map := nb.(map[string]interface{}) // assuming map, is it validated in schema??
	file_path := nb_map[NEXUS_BLOB_PATH_KEY].(string)

	exist := mutils.IsPathExist(file_path) // do we need error details??
	if !exist {                            // if path is invalid
		return fmt.Errorf("error in processing nexus repo file %v", file_path)
	}
	return nil
}

// Method to process vcs content, return error in case of issues
func (vs *validators) processVcs() error {

	vcs, vcs_present := vs.content[VCS_KEY] // extracting nexus repo key

	if !vcs_present {
		return nil // vcs key missing
	}
	vcs_map := vcs.(map[string]interface{})

	dir_path := vcs_map[VCS_PATH_KEY].(string)
	empty := mutils.IsEmptyDirectory(dir_path) // do we need error details??
	if empty {                                 // if path is invalid
		return fmt.Errorf("error in processing vcs directory %v", dir_path)
	}
	return nil
}

// Method to process rpm content, return error in case of issues
func (vs *validators) processRpm() error {

	rpm, rpm_present := vs.content[RPM_KEY] // extracting rpm key

	if !rpm_present {
		return nil // rpm key missing
	}

	rpm_array := rpm.([]interface{})

	for _, rpm := range rpm_array {
		rpm_map := rpm.(map[string]interface{})

		dir_path := rpm_map[RPM_PATH_KEY].(string)
		empty := mutils.IsEmptyDirectory(dir_path) // do we need error details??
		if empty {                                 // if path is invalid
			return fmt.Errorf("error in processing rpm directory %v", dir_path)
		}
		repoName := rpm_map[REPO_NAME].(string)

		found := false

		found, _ = mutils.StringFoundInArray(vs.hostedRepoNames, repoName)
		if !found {
			return fmt.Errorf("Repo referenced in  rpms section is not a hosted repo")
		}
	}

	return nil
}

// Function to extract content data
// Assumtion: Schema validation is done be the data comes to this function
func getManifestContentMap(manifest interface{}) map[string]interface{} {
	manifest_map := manifest.(map[string]interface{})
	content_map := manifest_map[CONTENT_KEY].(map[string]interface{})
	return content_map
}

// Function to validate manifest data post schema validation
func Validate(manifest interface{}) error {

	var pipeline *validators = &validators{}
	pipeline.content = getManifestContentMap(manifest)

	// var hosted_repo_names []string
	// content := getManifestContentMap(manifest)

	// content.s3 checks
	if err := pipeline.processS3(); err != nil {
		return fmt.Errorf("issue in processing s3 content, details %v", err)
	}

	// content.nexus_repositories checks
	if err := pipeline.processNexusRepo(); err != nil {
		return fmt.Errorf("issue in processing nexus repo file content, details %v", err)
	}

	// Validate content of nexus_repositories.yaml
	if err := pipeline.processNexusRepoFile(); err != nil {
		return fmt.Errorf("issue in processing nexus repo file content, details %v", err)
	}

	// content.nexus_blob_stores checks
	if err := pipeline.processNexusBlob(); err != nil {
		return fmt.Errorf("issue in processing nexus repo file content, details %v", err)
	}

	// contents.vcs checks
	if err := pipeline.processVcs(); err != nil {
		return fmt.Errorf("issue in processing vcs directory content, details %v", err)
	}

	// contents.rpms checks
	if err := pipeline.processRpm(); err != nil {
		return fmt.Errorf("issue in processing rpm directory content, details %v", err)
	}

	return nil
}
