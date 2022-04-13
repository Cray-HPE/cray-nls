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
	"fmt"
	"regexp"
)

type Validator struct {
}

// NewValidator creates a new environment
func NewValidator() Validator {

	return Validator{}
}

func (validator Validator) ValidateWorkerHostnames(hostnames []string) error {
	for _, hostname := range hostnames {
		err := validator.validateWorkerHostname(hostname)
		if err != nil {
			return err
		}
	}
	return nil
}

func (validator Validator) validateWorkerHostname(hostname string) error {
	isValid, err := regexp.Match(`^ncn-w[0-9]*$`, []byte(hostname))
	if err != nil {
		return err
	}

	if !isValid {
		return fmt.Errorf("invalid worker hostname: %s", hostname)
	}
	return nil
}
