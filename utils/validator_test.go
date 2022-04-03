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
