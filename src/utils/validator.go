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
	var firstHostname = hostnames[0]
	isWorker, err := validator.IsWorkerHostname(firstHostname)
	if err != nil {
		return err
	}
	if isWorker {
		return validator.ValidateWorkerHostnames(hostnames)
	}
	isStorage, err := validator.IsStorageHostname(firstHostname)
	if err != nil {
		return err
	}
	if isStorage {
		return validator.ValidateStorageHostnames(hostnames)
	}
	isMaster, err :=  validator.IsMasterHostname(firstHostname)
	if err != nil {
		return err
	}
	if isMaster {
		return validator.ValidateMasterHostnames(hostnames)
	}
	return fmt.Errorf("Invalid worker, storage, or master hostname: %s", firstHostname)
}

func (validator Validator) ValidateWorkerHostnames(hostnames []string) error {
	for _, hostname := range hostnames {
		isValid, err := validator.IsWorkerHostname(hostname)
		if err != nil {
			return err
		}
		if !isValid {
			isStorage, err := validator.IsStorageHostname(hostname)
			if err != nil {
				return err
			}
			isMaster, err := validator.IsMasterHostname(hostname)
			if err != nil {
				return err
			}
			if isStorage || isMaster {
				return fmt.Errorf("Cannot have both worker hostnames and storage or master hostnames. All hostnames must belong to exactly one management node type.")
			} else {
				return fmt.Errorf("Invalid worker hostname: %s", hostname)
			}
		}
	}
	return nil
}

func (validator Validator) ValidateStorageHostnames(hostnames []string) error {
	for _, hostname := range hostnames {
		isValid, err := validator.IsStorageHostname(hostname)
		if err != nil {
			return err
		}
		if !isValid {
			isWorker, err := validator.IsWorkerHostname(hostname)
			if err != nil {
				return err
			}
			isMaster, err := validator.IsMasterHostname(hostname)
			if err != nil {
				return err
			}
			if isWorker || isMaster {
				return fmt.Errorf("Cannot have both storage hostnames and worker or master hostnames. All hostnames must belong to exactly one management node type.")
			} else {
				return fmt.Errorf("Invalid storage hostname: %s", hostname)
			}
		}
	}
	return nil
}

func (validator Validator) ValidateMasterHostnames(hostnames []string) error {
	for _, hostname := range hostnames {
		isValid, err := validator.IsMasterHostname(hostname)
		if err != nil {
			return err
		}
		if !isValid {
			isWorker, err := validator.IsWorkerHostname(hostname)
			if err != nil {
				return err
			}
			isStorage, err := validator.IsStorageHostname(hostname)
			if err != nil {
				return err
			}
			if isWorker || isStorage {
				return fmt.Errorf("Cannot have both master hostnames and worker or storage hostnames. All hostnames must belong to exactly one management node type.")
			} else {
				return fmt.Errorf("Invalid storage hostname: %s", hostname)
			}
		}
	}
	return nil
}

func (validator Validator) IsWorkerHostname(hostname string) (bool, error) {
	return regexp.Match(`^ncn-w[0-9]*$`, []byte(hostname))
}

func (validator Validator) IsStorageHostname(hostname string) (bool, error) {
	return regexp.Match(`^ncn-s[0-9]*$`, []byte(hostname))
}

func (validator Validator) IsMasterHostname(hostname string) (bool, error) {
	return regexp.Match(`^ncn-m[0-9]*$`, []byte(hostname))
}

func (validator Validator) ValidateManagementRole(managementRole []string) error {
	
	if len(managementRole) > 1 {
		return fmt.Errorf("Exactly one Role_Subrole should be provided to '--limit-management-rollout'. Recieved %s", managementRole)
	}
	switch managementRole[0] {
		case
			"Management_Master",
			"Management_Storage",
			"Management_Worker":
			return nil
		}
		return fmt.Errorf("Invalid Role_Subrole: %s. Valid Role_Subroles are 'Management_Master', 'Management_Storage', and 'Management_Worker'.", managementRole[0])
}

func (validator Validator) ValidateLimitManagementNodesInput(limitManagementNodes []string) (string, error) {
	firstEntry := limitManagementNodes[0]
	if len(firstEntry) > 0 {
		matched_Management, err := regexp.MatchString(`Management.`, firstEntry)
		if err != nil {
			return "", err
		}
		if matched_Management {
			err := validator.ValidateManagementRole(limitManagementNodes)
			if err != nil {
				return "", err
			}
			return validator.getNCNWorkflowType(firstEntry)
		} else {
			matched_ncn, err := regexp.MatchString(`ncn.`,firstEntry)
			if err != nil {
				return "", err
			}
			if matched_ncn {
				err := validator.ValidateHostnames(limitManagementNodes)
				if err != nil {
					return "", err
				}
				return validator.getNCNWorkflowType(firstEntry)
			} else {
				return "", fmt.Errorf("getManagementNodesRolloutSubOperation: found input '--limit-management-nodes' did not have valid HSM Role_Subrole or hostname(s).")
			}
		}
	} else {
		return "", fmt.Errorf("getManagementNodesRolloutSubOperation: found input '--limit-management-nodes' had length 0. Management-nodes-rollout requires HSM Role_Subrole or hostnames to be specified.")
	}
}

func (validator Validator) getNCNWorkflowType(input string) (string, error) {
	// check if worker role or hostnames
	worker_role, err := regexp.MatchString(`.Worker`, input)
	if err != nil {
		return "", err
	}
	worker_hostname, err := validator.IsWorkerHostname(input)
	if err != nil {
		return "", err
	}
	if worker_role || worker_hostname {
		return "worker", nil
	}
	// check if storage role or hostnames
	storage_role, err := regexp.MatchString(`.Storage`, input)
	if err != nil {
		return "", err
	}
	storage_hostname, err := validator.IsStorageHostname(input)
	if err != nil {
		return "", err
	}
	if storage_role || storage_hostname {
		return "storage", nil
	}
	// check if master role or hostnames
	master_role, err := regexp.MatchString(`.Master`, input)
	if err != nil {
		return "", err
	}
	master_hostname, err := validator.IsMasterHostname(input)
	if err != nil {
		return "", err
	}
	if master_role || master_hostname {
		return "master", nil
	}
	return "", nil
}
