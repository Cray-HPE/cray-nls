/*
 *
 *  MIT License
 *
 *  (C) Copyright 2022-2023 Hewlett Packard Enterprise Development LP
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
package iuf

import (
	"fmt"
	"os"
	"path/filepath"

	manifestDataValidator "github.com/Cray-HPE/cray-nls/src/api/models/iuf/manifestDataValidation"
	sv "github.com/Cray-HPE/cray-nls/src/api/models/iuf/schemaValidator"
	"sigs.k8s.io/yaml"
)

const iuf_manifest_schema_file string = "schemas/iuf-manifest-schema.yaml"

var manifestRootPath string

func extractDirPath(filePath string, storageVariable *string) {
	*storageVariable = filepath.Dir(filePath)
}

func ValidateFile(file_path string) error {
	fmt.Printf("Validating file %v against IUF Product Manifest Schema.\n", file_path)

	// Extracting root dir of manifest
	extractDirPath(file_path, &manifestRootPath)

	file_contents, err := os.ReadFile(file_path)

	if err != nil {
		return fmt.Errorf("failed to open IUF Product Manifest file: %v", err)
	}

	err = Validate(file_contents)
	if err != nil {
		return fmt.Errorf("failed to validate IUF Product Manifest schema: %#v", err)
	}

	fmt.Printf("File %v is valid against IUF Product Manifest schema.\n", file_path)
	return nil
}

func Validate(file_contents []byte) error {

	var manifest interface{}
	err := yaml.Unmarshal(file_contents, &manifest)
	if err != nil {
		return fmt.Errorf("failed to parse IUF Product Manifest as YAML: %v", err)
	}

	// manifest schema validation
	err = sv.Validate(manifest, iuf_manifest_schema_file)
	if err != nil {
		return fmt.Errorf("failed to validate IUF Product Manifest schema: %#v", err)
	}

	// manifest data validation
	manifestDataValidator.SetManifestRootDir(manifestRootPath)
	err = manifestDataValidator.Validate(manifest)
	if err != nil {
		return fmt.Errorf("failed to validate IUF Product Manifest data: %#v", err)
	}
	return nil
}
