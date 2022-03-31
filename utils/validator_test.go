package utils

import (
	"fmt"
	"testing"
)

func TestValidator(t *testing.T) {
	var tests = []struct {
		hostname string
		wantErr  bool
	}{
		{"ncn-m001", false},
		{"ncn-w001", false},
		{"ncn-s001", false},
		{"ncn-m011", false},
		{"ncn-x001", true},
		{"sccn-m001", true},
		{"ncn-x001", true},
		{"ncn-m001asdf", true},
	}
	validator := NewValidator()
	for _, tt := range tests {
		testname := fmt.Sprintf("%s", tt.hostname)
		t.Run(testname, func(t *testing.T) {
			err := validator.ValidateHostname(tt.hostname)
			if (err != nil) != tt.wantErr {
				t.Errorf("got %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
