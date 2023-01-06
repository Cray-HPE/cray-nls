/*
 *
 *  MIT License
 *
 *  (C) Copyright 2022 Hewlett Packard Enterprise Development LP
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
	// The name of the product
	Name string `json:"name" validate:"required"`
	// The version of the product.
	Version string `json:"version" validate:"required"`
	// The original location of the extracted tar in on the physical storage.
	OriginalLocation string `json:"original_location"  validate:"required"`
	// The flag indicates md5 of a product tarball file has been validated
	Validated bool `json:"validated"  validate:"required"`
	// the content of manifest
	Manifest string `json:"manifest"`
} //	@name	Product

type InputParameters struct {
	MediaDir                 string   `json:"media_dir"`                  // Location of media
	SiteParameters           string   `json:"site_parameters"`            // DEPRECATED: use site_parameters at the top level of the activity or session resource. The inline contents of the site_parameters.yaml file.
	LimitNodes               []string `json:"limit_nodes"`                // Each item is the xname of a node
	BootprepConfigManaged    []string `json:"bootprep_config_managed"`    // Each item is the path to the bootprep config file for managed nodes, relative to the media_dir
	BootprepConfigManagement []string `json:"bootprep_config_management"` // Each item is the path to the bootprep config file for management nodes, relative to the media_dir
	Stages                   []string `json:"stages"`                     // Stages to execute
	Force                    bool     `json:"force"`                      // Force re-execution of stage operations
} //	@name	InputParameters

type SiteParameters struct {
	Global   map[string]interface{}            `json:"global"`   // global parameters applicable to all products
	Products map[string]map[string]interface{} `json:"products"` // Product-specific parameters
} //	@name	SiteParameters

type SiteParametersForOperationsAndHooks struct {
	SiteParameters
	CurrentProduct map[string]interface{} `json:"current_product"` // Current product's site parameters.
} //	@name	SiteParametersForOperationsAndHooks
