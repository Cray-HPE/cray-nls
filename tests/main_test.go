package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestBadLabels(t *testing.T) {

	tt := []struct {
		name           string
		workflow_label string
		double         int
		status         int
		err            string
		response       string
	}{
		{name: "bad label", workflow_label: "bad-label", response: "null"},

		{name: "another bad label", workflow_label: "another-bad-label", response: "null"},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			var myurl = "http://localhost:3000/apis/nls/v1/workflows?labelSelector=" + tc.workflow_label

			response, err := http.Get(myurl)

			if err != nil {
				t.Fatalf("could not complete request %v", err)
			}

			defer response.Body.Close()

			b, err := ioutil.ReadAll(response.Body)

			if response.StatusCode != http.StatusOK {
				t.Errorf("expected status OK got %v", response.StatusCode)
			}

			if msg := string(bytes.TrimSpace(b)); msg != tc.response {
				t.Errorf("expected response to be %v, got %v", tc.response, msg)
			}

		})
	}

}
func TestBadLabel(t *testing.T) {

	workflow_label := "bad-label"

	expected_response := "null"

	var myurl = "http://localhost:3000/apis/nls/v1/workflows?labelSelector=" + workflow_label

	response, err := http.Get(myurl)

	if err != nil {
		t.Fatalf("could not complete request %v", err)
	}

	defer response.Body.Close()

	b, err := ioutil.ReadAll(response.Body)

	if response.StatusCode != http.StatusOK {
		t.Errorf("expected status OK got %v", response.StatusCode)
	}

	if msg := string(bytes.TrimSpace(b)); msg != expected_response {
		t.Errorf("expected response to be %v, got %v", expected_response, msg)
	}

}
