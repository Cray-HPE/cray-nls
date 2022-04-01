package utils

import (
	"testing"
)

func TestHostnameValidator(t *testing.T) {
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
		t.Run(tt.hostname, func(t *testing.T) {
			err := validator.ValidateHostname(tt.hostname)
			if (err != nil) != tt.wantErr {
				t.Errorf("got %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestWorkerHostnameValidator(t *testing.T) {
	var tests = []struct {
		hostname string
		wantErr  bool
	}{
		{"ncn-m001", true},
		{"ncn-w001", false},
		{"ncn-s001", true},
		{"ncn-m011", true},
		{"ncn-x001", true},
		{"sccn-m001", true},
		{"ncn-x001", true},
		{"ncn-m001asdf", true},
	}
	validator := NewValidator()
	for _, tt := range tests {
		t.Run(tt.hostname, func(t *testing.T) {
			err := validator.ValidateWorkerHostname(tt.hostname)
			if (err != nil) != tt.wantErr {
				t.Errorf("got %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
