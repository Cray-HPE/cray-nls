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

	response, err := mutils.SplitMultiYamlFile(data)

	if err != nil {
		t.Fatal("Issue in parsing multidoc yaml")
	}

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

	response, err := mutils.SplitMultiYamlFile(data)

	if err != nil {
		t.Fatal("Issue in parsing multidoc yaml")
	}

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

	response, err := mutils.SplitMultiYamlFile(data)

	if err != nil {
		t.Fatal("Issue in parsing multidoc yaml")
	}

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

	response, err := mutils.SplitMultiYamlFile(data)

	if err != nil {
		t.Fatal("Issue in parsing multidoc yaml")
	}

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

	response, err := mutils.SplitMultiYamlFile(data)

	if err != nil {
		t.Fatal("Issue in parsing multidoc yaml")
	}

	for i, b := range expected {

		if string(b) != string(response[i]) {
			t.Fatal("Spilt operations is not working properly, expected:", string(b), "response got:", string(response[i]))
		}
	}
}
