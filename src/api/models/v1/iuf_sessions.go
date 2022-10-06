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
// +groupName=iuf.hpe.com
package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type IufSyncRequest struct {
	Parent IufSession `json:"parent"`
}

type IufSyncResponse struct {
	Status             IufStatus `json:"status,omitempty"`
	ResyncAfterSeconds int       `json:"resyncAfterSeconds,omitempty"`
}

type IufStatus struct {
	Phase string `json:"phase,omitempty"`
	// A 2-level DAG of Operations derived from stages that would be executed for each of the products that are specified. This is not specified by the Admin -- it is computed from the list of stages above.  This is an array of array of CR names of Operations that are installed as part of IUF, and determined by the Stages supplied.
	Operations [][]string `json:"operations,omitempty"`
	// The unique name of the Argo workflow that is created from all the input parameters above.
	ArgoWorkflow       string `json:"argo_workflow,omitempty"`
	ObservedGeneration int    `json:"observedGeneration"`
	Message            string `json:"message,omitempty"`
}

// IufSession
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=iufsessions,scope=Namespaced
// +kubebuilder:storageversion
type IufSession struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              IufSessionSpec `json:"spec"`
	Status            IufStatus      `json:"status,omitempty"`
}

type IufSessionProducts struct {
	// The name of the product
	Name string `json:"name"`
	// The version of the product.
	Version string `json:"version"`
	// The original location of the extracted tar in on the physical storage.
	OriginalLocation string `json:"original_location"`
	// The location of the product manifest uploaded into s3
	IufManifestS3Location string `json:"iuf_manifest_s3_location"`
	// Any before hook scripts for this product. This is an object where the key is operation name, and value is the CR name of the hook script for this product.  Hook scripts are executed either before or after an execution of a operation. They are specified in each product's distribution file, as part of the iuf-manifest.yaml.  The hook scripts are initially taken from the product distribution file and stored in S3, so that they can later be referenced.
	BeforeHookScripts map[string]string `json:"before_hook_scripts,omitempty"`
	// Any after hook scripts for this product. This is an object where the key is operation name, and value is the CR name of the hook script for this product.  Hook scripts are executed either before or after an execution of a operation. They are specified in each product's distribution file, as part of the iuf-manifest.yaml.  The hook scripts are initially taken from the product distribution file and stored in S3, so that they can later be referenced.
	AfterHookScripts map[string]string `json:"after_hook_scripts,omitempty"`
}

// The input parameters supplied by the Admin.
type IufSessionInputParams struct {
	// The pattern to use for all products. Use the following variables in braces {} to specify the pattern:  {product_name} {product_version}  E.g.  {product_name}-{product_version}-test-branch
	VcsWorkingBranchPattern string `json:"vcs_working_branch_pattern,omitempty"`
	// Specify the working branch name per product. This is an object where the key is the product name, and the value is the exact name (not a pattern) of the VCS branch for that product.
	VcsWorkingBranchPerProduct map[string]string `json:"vcs_working_branch_per_product,omitempty"`
}

// +kubebuilder:validation:Enum=install;upgrade
type WorkflowType string

// Node types
const (
	WorkflowTypeInstall WorkflowType = "install"
	WorkflowTypeUpgrade WorkflowType = "upgrade"
)

// An IUF session represents the intent of an Admin to initiate an install-upgrade workflow. It contains both input data, as well as any intermediary data that is needed to generate the final Argo workflow.
type IufSessionSpec struct {
	// What type of workflow are we executing? install or upgrade
	WorkflowType WorkflowType `json:"workflow_type"`
	// The products that need to be installed, as specified by the Admin.
	Products []IufSessionProducts `json:"products"`
	// The stages that need to be executed.
	// This is either explicitly specified by the Admin, or it is computed from the workflow type.
	// An Stage is a group of Operations. Stages represent the overall workflow at a high-level, and executing a stage means executing a bunch of Operations in a predefined manner.  An Admin can specify the stages that must be executed for an install-upgrade workflow. And Product Developers can extend each stage with custom hook scripts that they would like to run before and after the stage's execution.  The high-level stages allow their configuration would revealing too many details to the consumers of IUF.
	// if not specified, we apply all stages
	Stages []string `json:"stages"`

	InputParams *IufSessionInputParams `json:"input_params"`
}

func (s IufSessionSpec) GetProductsName() []string {
	res := []string{}
	for _, product := range s.Products {
		res = append(res, product.Name)
	}
	return res
}
