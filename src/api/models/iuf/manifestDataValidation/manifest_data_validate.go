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
package manifestDataValidation

import (
	"fmt"
	"path/filepath"

	mutils "github.com/Cray-HPE/cray-nls/src/api/models/iuf/mutils"
	"sigs.k8s.io/yaml"
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

// List of formats, will be skipped for nexus repo
var skipFormatsForNexusRepo = []string{"docker", "helm"}

// Getting the file reader
var FileReader = mutils.ReadYamFile

// Struct with list of validators
type validators struct {
	content           map[string]interface{}
	nexusRepoFileName string
	hostedRepoNames   []string
	manifestRootDir   string
}

// Method to process s3 content, returns error in case of issues
func (vs *validators) validateS3FilePath() error {

	s3, s3_present := vs.content[S3_KEY] // extracting s3 key

	if !s3_present {
		return nil // s3 key missing
	}

	s3_array := s3.([]interface{}) // assuming array, is it validated in schema??
	for _, s3 := range s3_array {
		s3_element := s3.(map[string]interface{})
		file_path := filepath.Join(vs.manifestRootDir, s3_element[S3_PATH_KEY].(string))

		exist := mutils.IsPathExist(file_path) // do we need error details??
		if !exist {                            // if path is invalid
			return fmt.Errorf("error in processing s3 file %v", file_path)
		}
	}
	return nil
}

// Method to process nexus repo content, returns repo file path and error(in case of issues),
func (vs *validators) validateNexusRepoFilePath() error {
	nr, nr_present := vs.content[NEXUS_REPO_KEY] // extracting nexus repo key

	if !nr_present {
		return nil // nexus repo key missing
	}

	nr_map := nr.(map[string]interface{}) // assuming map, is it validated in schema??
	file_path := filepath.Join(vs.manifestRootDir, nr_map[NEXUS_REPO_PATH_KEY].(string))
	vs.nexusRepoFileName = file_path

	exist := mutils.IsPathExist(file_path) // do we need error details??
	if !exist {                            // if path is invalid
		return fmt.Errorf("error in processing file %v, as it does not exist", file_path)
	}

	return nil
}

// Method to process nexus repo file and get hosted repo names
func (vs *validators) validateNexusRepoFileContent() error {

	if vs.nexusRepoFileName == "" {
		return nil // skip processing of nexus repo file
	}

	//var logger utils.Logger

	nexusFile_contents, err := FileReader(vs.nexusRepoFileName)

	var temp_repo_names []string

	if err != nil {
		return fmt.Errorf("failed to open Nexus Repository file: %v", err)
	}

	docs := mutils.SplitMultiYamlFile(nexusFile_contents)

	for _, doc := range docs {

		var nexusContentRaw interface{}
		err = yaml.Unmarshal(doc, &nexusContentRaw)
		if err != nil {
			return fmt.Errorf("failed to parse Nexus Repositorty as YAML: %v", err)
		}

		nexusContent := nexusContentRaw.(map[string]interface{})

		format := nexusContent[FORMAT].(string)

		isFormatToBeSkipped, _ := mutils.StringFoundInArray(skipFormatsForNexusRepo, format)

		if isFormatToBeSkipped {
			//logger.Infof("validateNexusRepoFileContent: Format %s skipped", format) // this line leading to panic
			//Need to add logger to print which format is skipped
			continue //Skip doc which has format that does not require validataion
		}

		if nexusContent[REPO_TYPE] == "hosted" {

			vs.hostedRepoNames = append(vs.hostedRepoNames, nexusContent["name"].(string))
			temp_repo_names = append(temp_repo_names, nexusContent["name"].(string))

		} else if nexusContent[REPO_TYPE] == "group" {
			group_map := nexusContent["group"].(map[string]interface{})

			for _, v := range group_map {
				memNames_array := v.([]interface{})
				for _, m := range memNames_array {

					memberRepo := m.(string)

					isHostedRepo, index := mutils.StringFoundInArray(temp_repo_names, memberRepo)
					if isHostedRepo {
						temp_repo_names, err = mutils.Delete(temp_repo_names, index)

						if err != nil {
							fmt.Println("repo defined in host repo is not listed in group repo")
						}

					} else {
						return fmt.Errorf("repo referenced in group does not match hosted repo or Hosted Repos are not listed before group repos")
					}
				}
			}
		}

	}

	return nil
}

// Method to process nexus blob content, returns error in case of issues
func (vs *validators) validateNexusBlobFilePath() error {
	nb, nb_present := vs.content[NEXUS_BLOB_KEY] // extracting nexus blob key

	if !nb_present {
		return nil // nexus blob key missing
	}

	nb_map := nb.(map[string]interface{})
	file_path := filepath.Join(vs.manifestRootDir, nb_map[NEXUS_BLOB_PATH_KEY].(string))

	exist := mutils.IsPathExist(file_path)
	if !exist { // if path is invalid
		return fmt.Errorf("error in processing file %v as it does not exist", file_path)
	}
	return nil
}

// Method to process vcs content, return error in case of issues
func (vs *validators) validateVcsFilePath() error {

	vcs, vcs_present := vs.content[VCS_KEY] // extracting nexus repo key

	if !vcs_present {
		return nil // vcs key missing
	}
	vcs_map := vcs.(map[string]interface{})

	dir_path := filepath.Join(vs.manifestRootDir, vcs_map[VCS_PATH_KEY].(string))
	empty := mutils.IsEmptyDirectory(dir_path)
	if empty { // if path is invalid
		return fmt.Errorf("error in processing vcs directory %v", dir_path)
	}
	return nil
}

func isRpmHosted(hostedRepoNames []string, repoName string) error {

	found := false

	found, _ = mutils.StringFoundInArray(hostedRepoNames, repoName)
	if !found {
		return fmt.Errorf("repo %v referenced in rpms section is not a hosted repo", repoName)
	}

	return nil
}

func isRpmExist(rootDir string, rpmName string) bool {
	dirPath := filepath.Join(rootDir, rpmName)
	return mutils.IsEmptyDirectory(dirPath)
}

// Method to process rpm content, return error in case of issues
func (vs *validators) validateRpmFiles() error {

	rpm, rpm_present := vs.content[RPM_KEY] // extracting rpm key

	if !rpm_present {
		return nil // rpm key missing
	}

	rpm_array := rpm.([]interface{})

	for _, rpm := range rpm_array {

		rpm_map := rpm.(map[string]interface{})
		rpmFile := rpm_map[RPM_PATH_KEY].(string)
		if isRpmExist(vs.manifestRootDir, rpmFile) { // if path is invalid
			return fmt.Errorf("error in processing rpm file %v in directory %v", rpmFile, vs.manifestRootDir)
		}

		repoName := rpm_map[REPO_NAME].(string)
		err := isRpmHosted(vs.hostedRepoNames, repoName) // Checking for hosted repo

		if err != nil {
			return fmt.Errorf("%#v", err)
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

var root_dir string

// Function of set manifest root path
func SetManifestRootDir(root_path string) {
	root_dir = root_path
}

// Function to validate manifest data post schema validation
func Validate(manifest interface{}) error {
	var pipeline *validators = &validators{}
	pipeline.manifestRootDir = root_dir
	pipeline.content = getManifestContentMap(manifest)

	// content.s3 checks
	if err := pipeline.validateS3FilePath(); err != nil {
		return fmt.Errorf("issue in processing s3 content, details %v", err)
	}

	// content.nexus_repositories checks
	if err := pipeline.validateNexusRepoFilePath(); err != nil {
		return fmt.Errorf("issue in validating nexus repo file path, details %v", err)
	}

	// Validate content of nexus_repositories.yaml
	if err := pipeline.validateNexusRepoFileContent(); err != nil {
		return fmt.Errorf("issue in validating nexus repo file content, details %v", err)
	}

	// content.nexus_blob_stores checks
	if err := pipeline.validateNexusBlobFilePath(); err != nil {
		return fmt.Errorf("issue in validating nexus blob file path, details %v", err)
	}

	// contents.vcs checks
	if err := pipeline.validateVcsFilePath(); err != nil {
		return fmt.Errorf("issue in validating vcs directory content, details %v", err)
	}

	// contents.rpms checks
	if err := pipeline.validateRpmFiles(); err != nil {
		return fmt.Errorf("issue in validating rpm directory content, details %v", err)
	}

	return nil
}
