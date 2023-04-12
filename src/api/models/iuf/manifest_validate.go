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
	mdv "cray-nls/src/api/models/iuf/manifestDataValidation"
	"embed"
	"fmt"
	"os"

	"github.com/santhosh-tekuri/jsonschema/v5"
	"sigs.k8s.io/yaml"
)

//go:embed iuf-manifest-schema.yaml
var rebuildWorkflowFS embed.FS

const iuf_manifest_schema_file string = "iuf-manifest-schema.yaml"

func ValidateFile(file_path string) error {
	fmt.Printf("Validating file %v against IUF Product Manifest Schema.\n", file_path)

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

	schema_file_contents, err := rebuildWorkflowFS.ReadFile(iuf_manifest_schema_file)
	if err != nil {
		return fmt.Errorf("failed to load IUF Product Manifest schema: %v", err)
	}

	schema_json, err := yaml.YAMLToJSON(schema_file_contents)
	if err != nil {
		return fmt.Errorf("failed to convert IUF Product Manifest from YAML to JSON: %v", err)
	}

	schema, err := jsonschema.CompileString(iuf_manifest_schema_file, string(schema_json))
	if err != nil {
		return fmt.Errorf("failed to validate IUF Product Manifest JSON Schema: %v", err)
	}

	err = schema.Validate(manifest)
	if err != nil {
		return fmt.Errorf("failed to validate IUF Product Manifest schema: %#v", err)
	}

	err = mdv.Validate(manifest)
	if err != nil {
		return fmt.Errorf("failed to validate IUF Product Manifest data: %#v", err)
	}
	return nil
}
