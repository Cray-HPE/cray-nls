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

func (validator Validator) ValidateHostnames(hostnames []string) error {
	firstHostname = hostnames[0]
	isWorker, err := regexp.Match(`^ncn-w[0-9]*$`, []byte(firstHostname))
	if err != nil {
		return err
	}
	if isWorker {
		return ValidateWorkerHostnames(hostnames)
	}
	else {
		isStorage, err := regexp.Match(`^ncn-s[0-9]*$`, []byte(firstHostname))
		if err != nil {
			return err
		}
		if isStorage {
			return ValidateStorageHostnames(hostnames)
		} else {
			return fmt.Errorf("invalid worker or storage hostname: %s", firstHostname)
		}
	}
}

func (validator Validator) ValidateWorkerHostnames(hostnames []string) error {
	for _, hostname := range hostnames {
		isValid, err := regexp.Match(`^ncn-w[0-9]*$`, []byte(hostname))
		if err != nil {
			return err
		}
		if !isValid {
			isStorage, err := regexp.Match(`^ncn-s[0-9]*$`, []byte(hostname))
			if isStorage {
				return fmt.Errorf("Cannot have both storage and worker node hostnames")
			} else {
				return fmt.Errorf("Invalid worker hostname: %s", hostname)
			}
		}
	}
	return nil
}

func (validator Validator) ValidateStorageHostnames(hostnames []string) error {
	for _, hostname := range hostnames {
		isValid, err := regexp.Match(`^ncn-s[0-9]*$`, []byte(hostname))
		if err != nil {
			return err
		}
		if !isValid {
			isWorker, err := regexp.Match(`^ncn-w[0-9]*$`, []byte(hostname))
			if isWorker {
				return fmt.Errorf("Cannot have both storage and worker node hostnames")
			} else {
				return fmt.Errorf("Invalid storage hostname: %s", hostname)
			}
		}
	return nil
	}
}
