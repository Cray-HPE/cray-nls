package utils

import (
	"fmt"
	"testing"
)

func TestValidator(t *testing.T) {
	var tests = []struct {
		hostname string
		want     bool
	}{
		{"ncn-m001", true},
		{"ncn-w001", true},
		{"ncn-s001", true},
		{"ncn-m011", true},
		{"ncn-x001", false},
		{"sccn-m001", false},
		{"ncn-x001", false},
	}
	validator := NewValidator()
	for _, tt := range tests {
		testname := fmt.Sprintf("%s", tt.hostname)
		t.Run(testname, func(t *testing.T) {
			ans, _ := validator.ValidateHostname(tt.hostname)
			if ans != tt.want {
				t.Errorf("got %v, want %v", ans, tt.want)
			}
		})
	}
}
