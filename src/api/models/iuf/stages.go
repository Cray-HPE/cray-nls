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

// Stage
type Stage struct {
	Name                               string       `yaml:"name" json:"name" binding:"required"`                                                // Name of the stage
	Type                               string       `yaml:"type" json:"type" binding:"required"`                                                // Type of the stage
	Operations                         []Operations `yaml:"operations" json:"operations" binding:"required"`                                    // operations
	NoHooks                            bool         `yaml:"no-hooks" json:"no-hooks"`                                                           // no-hook indicates that there are no hooks that should be run for this stage
	ProcessProductVariantsSequentially bool         `yaml:"process-product-variants-sequentially" json:"process-product-variants-sequentially"` // this stage wants to make sure all products with the same name (but different versions) are processed sequentially, not in parallel, to avoid operational race conditions
} //	@name	Stage

type Stages struct {
	Version string            `yaml:"version" json:"version" binding:"required"`
	Stages  []Stage           `yaml:"stages" json:"stages" binding:"required"`
	Hooks   map[string]string `yaml:"hooks" json:"hooks"`
} //	@name	Stages

type Operations struct {
	Name                              string                 `yaml:"name" json:"name" binding:"required"` // Name of the operation
	StaticParameters                  map[string]interface{} `yaml:"static-parameters" json:"static-parameters" binding:"required"`
	RequiredManifestAttributes        []string               `yaml:"required-manifest-attributes" json:"required-manifest-attributes"`
	IncludeDefaultProductInSiteParams bool                   `yaml:"include-default-product-in-site-params" json:"include-default-product-in-site-params"`
} //	@name	Operations
