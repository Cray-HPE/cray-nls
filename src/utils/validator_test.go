//
//  MIT License
//
//  (C) Copyright 2022 Hewlett Packard Enterprise Development LP
//
//  Permission is hereby granted, free of charge, to any person obtaining a
//  copy of this software and associated documentation files (the "Software"),
//  to deal in the Software without restriction, including without limitation
//  the rights to use, copy, modify, merge, publish, distribute, sublicense,
//  and/or sell copies of the Software, and to permit persons to whom the
//  Software is furnished to do so, subject to the following conditions:
//
//  The above copyright notice and this permission notice shall be included
//  in all copies or substantial portions of the Software.
//
//  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
//  IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
//  FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
//  THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR
//  OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
//  ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
//  OTHER DEALINGS IN THE SOFTWARE.
//
package utils

import (
	"testing"
)

func TestWorkerHostnameValidator(t *testing.T) {
	var tests = []struct {
		hostnames []string
		wantErr  bool
	}{
		{[]string{"ncn-m001"}, true},
		{[]string{"ncn-w001"}, false},
		{[]string{"ncn-w001", "ncn-m002"}, true},
		{[]string{"ncn-s001"}, true},
		{[]string{"ncn-m011"}, true},
		{[]string{"ncn-w001", "ncn-s002"}, true},
		{[]string{"ncn-w001", "ncn-w002"}, false},
		{[]string{"ncn-x001"}, true},
		{[]string{"sccn-m001"}, true},
		{[]string{"ncn-x001"}, true},
		{[]string{"ncn-m001asdf"}, true},
	}
	validator := NewValidator()
	for _, tt := range tests {
		t.Run(tt.hostnames[0], func(t *testing.T) {
			err := validator.ValidateWorkerHostnames(tt.hostnames)
			if (err != nil) != tt.wantErr {
				t.Errorf("got %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStorageHostnamesValidator(t *testing.T) {
	var tests = []struct {
		hostnames []string
		wantErr   bool
	}{
		{[]string{"ncn-m001"}, true},
		{[]string{"ncn-w001"}, true},
		{[]string{"ncn-w001", "ncn-m002"}, true},
		{[]string{"ncn-s001"}, false},
		{[]string{"ncn-m011"}, true},
		{[]string{"ncn-w001", "ncn-s002"}, true},
		{[]string{"ncn-s001", "ncn-s002"}, false},
		{[]string{"ncn-x001"}, true},
		{[]string{"sccn-m001"}, true},
		{[]string{"ncn-x001"}, true},
		{[]string{"ncn-m001asdf"}, true},
	}
	validator := NewValidator()
	for _, tt := range tests {
		t.Run(tt.hostnames[0], func(t *testing.T) {
			err := validator.ValidateStorageHostnames(tt.hostnames)
			if (err != nil) != tt.wantErr {
				t.Errorf("got %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}


func TestHostnamesValidator(t *testing.T) {
	var tests = []struct {
		hostnames []string
		wantErr   bool
	}{
		{[]string{"ncn-m001"}, false},
		{[]string{"ncn-w001"}, false},
		{[]string{"ncn-w001", "ncn-m002"}, true},
		{[]string{"ncn-w001", "ncn-w002"}, false},
		{[]string{"ncn-s001"}, false},
		{[]string{"ncn-m011"}, false},
		{[]string{"ncn-w001", "ncn-s002"}, true},
		{[]string{"ncn-s001", "ncn-s002"}, false},
		{[]string{"ncn-x001"}, true},
		{[]string{"sccn-m001"}, true},
		{[]string{"ncn-x001"}, true},
		{[]string{"ncn-m001asdf"}, true},
	}
	validator := NewValidator()
	for _, tt := range tests {
		t.Run(tt.hostnames[0], func(t *testing.T) {
			err := validator.ValidateHostnames(tt.hostnames)
			if (err != nil) != tt.wantErr {
				t.Errorf("got %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestLimitManagementNodesInputValidator(t *testing.T) {
	var tests = []struct {
		LimitManagedNodes []string
		wantErr   bool
		expectedWorkflow string
	}{
		{[]string{"ncn-m001"}, false, "master"},
		{[]string{"ncn-m001", "ncn-m002"}, false, "master"},
		{[]string{"ncn-m001", "ncn-s001"}, true, ""},
		{[]string{"ncn-w001"}, false, "worker"},
		{[]string{"ncn-w001", "ncn-w002", "ncn-w003"}, false, "worker"},
		{[]string{"ncn-w001", "ncn-s002", "ncn-w003"}, true, ""},
		{[]string{"ncn-s001"}, false, "storage"},
		{[]string{"ncn-s001", "ncn-s004"}, false, "storage"},
		{[]string{"ncn-s001", "ncn-m004", "ncn-w003"}, true, ""},
		{[]string{"Management_Storage"}, false, "storage"},
		{[]string{"Management_Worker"}, false, "worker"},
		{[]string{"Management_Master"}, false, "master"},
		{[]string{"Management_Storage", "Management_Master"}, true, ""},
		{[]string{"Management_Storage", "ncn-s001"}, true, ""},
		{[]string{"Management_BADNAME"}, true, ""},
		{[]string{"ncn-bad"}, true, ""},
	}
	validator := NewValidator()
	for _, tt := range tests {
		t.Run(tt.LimitManagedNodes[0], func(t *testing.T) {
			workFlowType, err := validator.ValidateLimitManagementNodesInput(tt.LimitManagedNodes)
			if (err != nil) != tt.wantErr {
				t.Errorf("got %v, wantErr %v", err, tt.wantErr)
				return
			}
			if workFlowType != tt.expectedWorkflow {
				t.Errorf("got %v, wantType %v", workFlowType, tt.expectedWorkflow)
				return
			}
		})
	}
}