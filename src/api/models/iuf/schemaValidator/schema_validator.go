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
package schemaValidator

import (
	"embed"
	"fmt"

	"github.com/santhosh-tekuri/jsonschema/v5"
	"sigs.k8s.io/yaml"
)

//go:embed schemas/*
var schemas embed.FS

// Validates schema of a yaml file
// Provide the yaml data as object and schema file to process
func Validate(yamlObject interface{}, schemaFile string) error {

	schema_file_contents, err := schemas.ReadFile(schemaFile)
	if err != nil {
		return fmt.Errorf("failed to load schema file: %s, error: %v", schemaFile, err)
	}

	schema_json, err := yaml.YAMLToJSON(schema_file_contents)
	if err != nil {
		return fmt.Errorf("failed to convert schema file: %s from YAML to JSON: %v", schemaFile, err)
	}

	schema, err := jsonschema.CompileString(schemaFile, string(schema_json))
	if err != nil {
		return fmt.Errorf("failed to validate schema file: %s JSON Schema: %v", schemaFile, err)
	}

	err = schema.Validate(yamlObject)
	if err != nil {
		return fmt.Errorf("failed to validate yaml: %#v", err)
	}

	return nil
}
