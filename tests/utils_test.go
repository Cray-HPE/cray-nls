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

	"github.com/Cray-HPE/cray-nls/src/api/models/iuf/mutils"
)

func TestMultiSchemaYamlDocValid(t *testing.T) {
	data := []byte(`
---
iuf_version: ^0.5.0
---
iuf_version: ^0.5.1
`)

	expected := [][]byte{
		[]byte(`
iuf_version: ^0.5.0
`),
		[]byte(`
iuf_version: ^0.5.1
`)}

	response := mutils.SplitMultiYamlFile(data)

	for i, b := range expected {

		if string(b) != string(response[i]) {
			t.Fatal("Spilt operations is not working properly, expected:", string(b), "response got:", string(response[i]))
		}
	}
}

func TestMultiSchemaYamlDocEmptyYaml(t *testing.T) {
	data := []byte(`
---
---
`)

	expected := [][]byte{}

	response := mutils.SplitMultiYamlFile(data)

	for i, b := range expected {

		if string(b) != string(response[i]) {
			t.Fatal("Spilt operations is not working properly, expected:", string(b), "response got:", string(response[i]))
		}
	}
}

func TestMultiSchemaYamlDocEmptyDoc(t *testing.T) {
	data := []byte(`

`)

	expected := [][]byte{}

	response := mutils.SplitMultiYamlFile(data)

	for i, b := range expected {

		if string(b) != string(response[i]) {
			t.Fatal("Spilt operations is not working properly, expected:", string(b), "response got:", string(response[i]))
		}
	}
}

func TestMultiSchemaYamlDocSingle(t *testing.T) {
	data := []byte(`
iuf_version: ^0.5.0
`)

	expected := [][]byte{
		[]byte(`
iuf_version: ^0.5.0
`)}

	response := mutils.SplitMultiYamlFile(data)

	for i, b := range expected {

		if string(b) != string(response[i]) {
			t.Fatal("Spilt operations is not working properly, expected:", string(b), "response got:", string(response[i]))
		}
	}
}

func TestMultiSchemaYamlDocSingle_2(t *testing.T) {
	data := []byte(`
---
iuf_version: ^0.5.0
`)

	expected := [][]byte{
		[]byte(`
iuf_version: ^0.5.0
`)}

	response := mutils.SplitMultiYamlFile(data)

	for i, b := range expected {

		if string(b) != string(response[i]) {
			t.Fatal("Spilt operations is not working properly, expected:", string(b), "response got:", string(response[i]))
		}
	}
}

func TestMultiSchemaYamlDocSingle_3(t *testing.T) {
	data := []byte(`---
iuf_version: ^0.5.0
`)

	expected := [][]byte{
		[]byte(`
iuf_version: ^0.5.0
`)}

	response := mutils.SplitMultiYamlFile(data)

	for i, b := range expected {

		if string(b) != string(response[i]) {
			t.Fatal("Spilt operations is not working properly, expected:", string(b), "response got:", string(response[i]))
		}
	}
}
