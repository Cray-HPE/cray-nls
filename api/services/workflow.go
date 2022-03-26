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
package services

import (
	"context"
	_ "embed"

	"github.com/argoproj/pkg/json"
	"sigs.k8s.io/yaml"

	"github.com/Cray-HPE/cray-nls/utils"
	"github.com/argoproj/argo-workflows/v3/pkg/apiclient"
	"github.com/argoproj/argo-workflows/v3/pkg/apiclient/workflow"
	"github.com/argoproj/argo-workflows/v3/pkg/apiclient/workflowtemplate"
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/gin-gonic/gin"
)

//go:embed template.argo.yaml
var ArgoWorkflowTemplate []byte

//go:embed workflow.argo.yaml
var ArgoWorkflow []byte

// WorkflowService service layer
type WorkflowService struct {
	logger                utils.Logger
	ctx                   context.Context
	workflowCient         workflow.WorkflowServiceClient
	workflowTemplateCient workflowtemplate.WorkflowTemplateServiceClient
}

// NewWorkflowService creates a new Workflowservice
func NewWorkflowService(logger utils.Logger) WorkflowService {
	var argoOps apiclient.Opts = apiclient.Opts{
		ArgoServerOpts: apiclient.ArgoServerOpts{
			URL:                "localhost:2746",
			InsecureSkipVerify: true,
			Secure:             true,
			HTTP1:              true,
		},
		AuthSupplier: func() string {
			return `ZXlKaGJHY2lPaUpTVXpJMU5pSXNJbXRwWkNJNklrRkZNMFkwUjNjMWRWZE5iSFJqZDBoeE4wSmtXRkJJVTNSTWNYQjNia1pHY2tRNU5IUnpTbFpVZEdjaWZRLmV5SnBjM01pT2lKcmRXSmxjbTVsZEdWekwzTmxjblpwWTJWaFkyTnZkVzUwSWl3aWEzVmlaWEp1WlhSbGN5NXBieTl6WlhKMmFXTmxZV05qYjNWdWRDOXVZVzFsYzNCaFkyVWlPaUpoY21kdklpd2lhM1ZpWlhKdVpYUmxjeTVwYnk5elpYSjJhV05sWVdOamIzVnVkQzl6WldOeVpYUXVibUZ0WlNJNkltRnlaMjh0ZEc5clpXNHROamN5ZDNvaUxDSnJkV0psY201bGRHVnpMbWx2TDNObGNuWnBZMlZoWTJOdmRXNTBMM05sY25acFkyVXRZV05qYjNWdWRDNXVZVzFsSWpvaVlYSm5ieUlzSW10MVltVnlibVYwWlhNdWFXOHZjMlZ5ZG1salpXRmpZMjkxYm5RdmMyVnlkbWxqWlMxaFkyTnZkVzUwTG5WcFpDSTZJbU0wTURFNVlXSmtMVE15TkRVdE5EWTBNUzA0WVdOa0xXVmhZbU16T0dVNVpqRXlNU0lzSW5OMVlpSTZJbk41YzNSbGJUcHpaWEoyYVdObFlXTmpiM1Z1ZERwaGNtZHZPbUZ5WjI4aWZRLkNLaTdkMkVRQzNmS0tUdzdqVFE3ajB4VUV4eTM0Z3JVZ2c5elQ3amhTQ0xaZ3d1SngxVXpYbURaOWhUTDRmQ0pkSEw0SElTZlNYNE1ETUxmQ2lWN1U3dXEyMGRWQ3Q2Nlh0WWhObXlpTEVmU05IRE5kSWVuUU96amc0MEZSUkMzcjd4MmptVDNuejJmc2lRTGwwZDJDaUxTcjBYZDY4MC0xc2gwXzZxSlRfeS1yM2F4R1puZHV6SzJFbXBYU0J4VTJ5ZlFDQ2JHWFF2MmRjaXBDSmdEeEtzaXkzMnc5STltWnhCb1FtOHdEU2lEZTIteGxPOThoLS1GM0U2WEJwUHE1OW85UEYyVVhNYXdKSXVjX3kyUERUdzJzSXhwbkpNZXVGaWRzWXN4dVUtdFJIMkNyTFJhUlIxMzJhMVNCazJJcTRfZUM4WmNZRFJDMHF6VzRFbUJfUQ==`
		},
	}
	ctx, serviceClient, _ := apiclient.NewClientFromOpts(argoOps)
	workflowTemplateCient, _ := serviceClient.NewWorkflowTemplateServiceClient()

	defer logger.Infof("%s", "Workflow(Template) service initialized")

	res := WorkflowService{
		logger:                logger,
		ctx:                   ctx,
		workflowCient:         serviceClient.NewWorkflowServiceClient(),
		workflowTemplateCient: workflowTemplateCient,
	}
	res.initializeWorkflowTemplate(ArgoWorkflowTemplate)
	return res
}

func (s WorkflowService) GetWorkflows(ctx *gin.Context) (*v1alpha1.WorkflowList, error) {
	return s.workflowCient.ListWorkflows(s.ctx, &workflow.WorkflowListRequest{Namespace: "argo"})
}

func (s WorkflowService) initializeWorkflowTemplate(template []byte) error {
	s.logger.Infof("initialize workflow template")

	var myWorkflowTemplate v1alpha1.WorkflowTemplate
	tmpBytes, _ := yaml.YAMLToJSON(template)
	err := json.Unmarshal(tmpBytes, &myWorkflowTemplate)
	if err != nil {
		s.logger.Error(err)
	}

	workflowTemplateList, err := s.workflowTemplateCient.ListWorkflowTemplates(s.ctx, &workflowtemplate.WorkflowTemplateListRequest{Namespace: "argo"})

	initialized := false
	for _, workflowTemplate := range workflowTemplateList.Items {
		if workflowTemplate.Name == myWorkflowTemplate.Name {
			initialized = true
			s.logger.Info("workflow template has already been initialized")
			break
		}
	}

	if !initialized {
		_, err = s.workflowTemplateCient.CreateWorkflowTemplate(
			s.ctx,
			&workflowtemplate.WorkflowTemplateCreateRequest{
				Namespace: "argo",
				Template:  &myWorkflowTemplate,
			})
		if err != nil {
			s.logger.Error(err)
		}
	}

	return err
}
