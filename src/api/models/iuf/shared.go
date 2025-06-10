/*
 *
 *  MIT License
 *
 *  (C) Copyright 2022,2025 Hewlett Packard Enterprise Development LP
 *
 *  Permission is hereby granted, free of charge, to any person obtaining a
 *  copy of this software and associated documentation files (the "Software"),
 *  to deal in the Software without restriction, including without limitation
 *  the rights to use, copy, modify, merge, publish, distribute, sublicense,
 *  and/or sell copies of the Software, and to permit persons to whom the
 *  Software is furnished to do so, subject to the following conditions:
 *
 *  The above copyright notice and this permission notice shall be included
 *  in all copies or substantial portions of the Software.
 *
 *  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *  IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *  FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
 *  THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR
 *  OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
 *  ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
 *  OTHER DEALINGS IN THE SOFTWARE.
 *
 */
package iuf

type Product struct {
	Name             string `json:"name" validate:"required"`               // The name of the product
	Version          string `json:"version" validate:"required"`            // The version of the product.
	OriginalLocation string `json:"original_location"  validate:"required"` // The original location of the extracted tar in on the physical storage.
	Validated        bool   `json:"validated"  validate:"required"`         // The flag indicates md5 of a product tarball file has been validated
	Manifest         string `json:"manifest"`                               // the content of manifest
} //	@name	Product

type InputParameters struct {
	MediaDir                              string                     `json:"media_dir"`                                          // Location of media
	SiteParameters                        string                     `json:"site_parameters"`                                    // DEPRECATED: use site_parameters at the top level of the activity or session resource. The inline contents of the site_parameters.yaml file.
	LimitManagementNodes                  []string                   `json:"limit_management_nodes"`                             // Must in the form <role>_<subrole>. E.g. Management_Master, Management_Worker, Management_Storage
	LimitManagedNodes                     []string                   `json:"limit_managed_nodes"`                                // Anything accepted by BOS v2 as the value to a session's limit parameter.
	ManagementRolloutStrategy             EManagementRolloutStrategy `json:"management_rollout_strategy" enums:"reboot,rebuild"` // Whether to use a reboot or rebuild strategy for management nodes.
	ManagedRolloutStrategy                EManagedRolloutStrategy    `json:"managed_rollout_strategy" enums:"reboot,stage"`      // Whether to use a reboot or staged rollout strategy for managed nodes. Refer to BOS v2 for more details.
	ConcurrentManagementRolloutPercentage int64                      `json:"concurrent_management_rollout_percentage"`           // The percentage of management nodes to reboot in parallel before moving on to the next set of management nodes to reboot.
	MediaHost                             string                     `json:"media_host"`                                         // A string containing the hostname of where the media is located
	Concurrency                           int64                      `json:"concurrency"`                                        // An integer defining how many products / operations can we concurrently execute.
	BootprepConfigManaged                 string                     `json:"bootprep_config_managed"`                            // The path to the bootprep config file for managed nodes, relative to the media_dir
	BootprepConfigManagement              string                     `json:"bootprep_config_management"`                         // The path to the bootprep config file for management nodes, relative to the media_dir
	CfsConfigurationManagement            string                     `json:"cfs_configuration_management"`                       // The name of the cfs configuration for management nodes
	BootImageManagement                   string                     `json:"boot_image_management"`                              // The name of the boot image to be used for management nodes
	Stages                                []string                   `json:"stages"`                                             // Stages to execute
	Force                                 bool                       `json:"force"`                                              // Force re-execution of stage operations
} //	@name	InputParameters

type InputParametersPatch struct {
	MediaDir                              *string                     `json:"media_dir"`                                          // Location of media
	SiteParameters                        *string                     `json:"site_parameters"`                                    // DEPRECATED: use site_parameters at the top level of the activity or session resource. The inline contents of the site_parameters.yaml file.
	LimitManagementNodes                  *[]string                   `json:"limit_management_nodes"`                             // Must in the form <role>_<subrole>. E.g. Management_Master, Management_Worker, Management_Storage
	LimitManagedNodes                     *[]string                   `json:"limit_managed_nodes"`                                // Anything accepted by BOS v2 as the value to a session's limit parameter.
	ManagementRolloutStrategy             *EManagementRolloutStrategy `json:"management_rollout_strategy" enums:"reboot,rebuild"` // Whether to use a reboot or rebuild strategy for management nodes.
	ManagedRolloutStrategy                *EManagedRolloutStrategy    `json:"managed_rollout_strategy" enums:"reboot,stage"`      // Whether to use a reboot or staged rollout strategy for managed nodes. Refer to BOS v2 for more details.
	ConcurrentManagementRolloutPercentage *int64                      `json:"concurrent_management_rollout_percentage"`           // The percentage of management nodes to reboot in parallel before moving on to the next set of management nodes to reboot.
	MediaHost                             *string                     `json:"media_host"`                                         // A string containing the hostname of where the media is located
	Concurrency                           *int64                      `json:"concurrency"`                                        // An integer defining how many products / operations can we concurrently execute.
	BootprepConfigManaged                 *string                     `json:"bootprep_config_managed"`                            // The path to the bootprep config file for managed nodes, relative to the media_dir
	BootprepConfigManagement              *string                     `json:"bootprep_config_management"`                         // The path to the bootprep config file for management nodes, relative to the media_dir
	CfsConfigurationManagement            *string                     `json:"cfs_configuration_management"`                       // The name of the cfs configuration for management nodes
	BootImageManagement                   *string                     `json:"boot_image_management"`                              // The name of the boot image to be used for management nodes
	Stages                                *[]string                   `json:"stages"`                                             // Stages to execute
	Force                                 *bool                       `json:"force"`                                              // Force re-execution of stage operations
} //	@name	InputParameters

type EManagementRolloutStrategy string

const (
	EManagementRolloutStrategyReboot  EManagementRolloutStrategy = "reboot"
	EManagementRolloutStrategyRebuild EManagementRolloutStrategy = "rebuild"
)

type EManagedRolloutStrategy string

const (
	EManagedRolloutStrategyReboot EManagedRolloutStrategy = "reboot"
	EManagedRolloutStrategyStaged EManagedRolloutStrategy = "stage"
)

type SiteParameters struct {
	Global   map[string]interface{}            `json:"global"`   // global parameters applicable to all products
	Products map[string]map[string]interface{} `json:"products"` // Product-specific parameters
} //	@name	SiteParameters

type SiteParametersForOperationsAndHooks struct {
	SiteParameters
	CurrentProduct map[string]interface{} `json:"current_product"` // Current product's site parameters.
} //	@name	SiteParametersForOperationsAndHooks
