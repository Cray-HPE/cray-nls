/*
 *
 *  MIT License
 *
 *  (C) Copyright 2022-2024 Hewlett Packard Enterprise Development LP
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
package argo_templates

import (
	"bytes"
	_ "embed"
	"fmt"
	"io/fs"
	"text/template"

	models_iuf "github.com/Cray-HPE/cray-nls/src/api/models/iuf"
	models_nls "github.com/Cray-HPE/cray-nls/src/api/models/nls"
	"github.com/Cray-HPE/cray-nls/src/utils"
	"github.com/Masterminds/sprig/v3"
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var validator utils.Validator = utils.NewValidator()

func GetWorkerRebuildWorkflow(workerRebuildWorkflowFS fs.FS, createRebuildWorkflowRequest models_nls.CreateRebuildWorkflowRequest, rebuildHooks models_nls.RebuildHooks) ([]byte, error) {
	err := validator.ValidateWorkerHostnames(createRebuildWorkflowRequest.Hosts)
	if err != nil {
		return nil, err
	}

	tmpl := template.New("worker.rebuild.yaml")

	return GetRebuildWorkflow(tmpl, workerRebuildWorkflowFS, createRebuildWorkflowRequest, rebuildHooks)
}

func GetIufInstallWorkflow(iufInstallWorkflowFS fs.FS, req models_iuf.Session, stageIndex int) ([]byte, error) {
	tmpl := template.New("install.yaml")

	return GetIufWorkflow(tmpl, iufInstallWorkflowFS, req, stageIndex)
}

func GetStorageRebuildWorkflow(storageRebuildWorkflowFS fs.FS, createRebuildWorkflowRequest models_nls.CreateRebuildWorkflowRequest) ([]byte, error) {
	err := validator.ValidateStorageHostnames(createRebuildWorkflowRequest.Hosts)
	if err != nil {
		return nil, err
	}

	tmpl := template.New("storage.rebuild.yaml")

	return GetRebuildWorkflow(tmpl, storageRebuildWorkflowFS, createRebuildWorkflowRequest, models_nls.RebuildHooks{})
}

func GetStorageUpgradeWorkflow(storageRebuildWorkflowFS fs.FS, createRebuildWorkflowRequest models_nls.CreateRebuildWorkflowRequest) ([]byte, error) {
	err := validator.ValidateStorageHostnames(createRebuildWorkflowRequest.Hosts)
	if err != nil {
		return nil, err
	}

	tmpl := template.New("storage.upgrade.yaml")

	return GetRebuildWorkflow(tmpl, storageRebuildWorkflowFS, createRebuildWorkflowRequest, models_nls.RebuildHooks{})
}

func GetRebuildWorkflow(tmpl *template.Template, workflowFS fs.FS, createRebuildWorkflowRequest models_nls.CreateRebuildWorkflowRequest, rebuildHooks models_nls.RebuildHooks) ([]byte, error) {
	// add useful helm templating func: include
	var funcMap template.FuncMap = map[string]interface{}{}
	funcMap["include"] = func(name string, data interface{}) (string, error) {
		buf := bytes.NewBuffer(nil)
		if err := tmpl.ExecuteTemplate(buf, name, data); err != nil {
			return "", err
		}
		return buf.String(), nil
	}

	// add templating func: getHooks
	funcMap["getHooks"] = func(name string, data interface{}) (string, error) {
		dag := v1alpha1.DAGTemplate{}
		var unstructuredHooks []unstructured.Unstructured
		switch name {
		case "before-all":
			unstructuredHooks = rebuildHooks.BeforeAll
		case "before-each":
			unstructuredHooks = rebuildHooks.BeforeEach
		case "after-each":
			unstructuredHooks = rebuildHooks.AfterEach
		case "after-all":
			unstructuredHooks = rebuildHooks.AfterAll
		}

		for _, unstrunstructuredHook := range unstructuredHooks {
			dag.Tasks = append(dag.Tasks, v1alpha1.DAGTask{
				Name: unstrunstructuredHook.GetName(),
				TemplateRef: &v1alpha1.TemplateRef{
					Name:     fmt.Sprintf("%v", unstrunstructuredHook.Object["spec"].(map[string]interface{})["templateRefName"]),
					Template: "shell-script",
				},
				Arguments: v1alpha1.Arguments{
					Parameters: []v1alpha1.Parameter{
						{
							Name:  "scriptContent",
							Value: v1alpha1.AnyStringPtr(fmt.Sprintf("%v", unstrunstructuredHook.Object["spec"].(map[string]interface{})["scriptContent"])),
						},
						{
							Name:  "dryRun",
							Value: v1alpha1.AnyStringPtr(createRebuildWorkflowRequest.DryRun),
						},
						{
							Name:  "bootTimeoutInSeconds",
							Value: v1alpha1.AnyStringPtr(fmt.Sprintf("%v", unstrunstructuredHook.Object["spec"].(map[string]interface{})["bootTimeoutInSeconds"])),
						},
					},
				},
			})
		}

		// set minimum timeout if not specified
		if createRebuildWorkflowRequest.BootTimeoutInSeconds == 0 {
			createRebuildWorkflowRequest.BootTimeoutInSeconds = 600
		}

		if len(dag.Tasks) == 0 {
			dag.Tasks = append(dag.Tasks, v1alpha1.DAGTask{
				Name: "dummy-hook",
				TemplateRef: &v1alpha1.TemplateRef{
					Name:     "ssh-template",
					Template: "shell-script",
				},
				Arguments: v1alpha1.Arguments{
					Parameters: []v1alpha1.Parameter{
						{
							Name:  "scriptContent",
							Value: v1alpha1.AnyStringPtr("echo hello"),
						},
						{
							Name:  "dryRun",
							Value: v1alpha1.AnyStringPtr(createRebuildWorkflowRequest.DryRun),
						},
						{
							Name:  "bootTimeoutInSeconds",
							Value: v1alpha1.AnyStringPtr(createRebuildWorkflowRequest.BootTimeoutInSeconds),
						},
					},
				},
			})
		}
		res, _ := yaml.Marshal(dag.Tasks)
		return string(res), nil
	}

	// add sprig templating func
	tmpl, err := tmpl.Funcs(sprig.TxtFuncMap()).Funcs(funcMap).ParseFS(workflowFS, "**/*.yaml")
	if err != nil {
		return nil, err
	}

	var tmpRes bytes.Buffer
	err = tmpl.Execute(&tmpRes, map[string]interface{}{
		"TargetNcns":           createRebuildWorkflowRequest.Hosts,
		"DryRun":               createRebuildWorkflowRequest.DryRun,
		"ZapOsds":              createRebuildWorkflowRequest.ZapOsds,
		"WorkflowType":         createRebuildWorkflowRequest.WorkflowType,
		"ImageId":              createRebuildWorkflowRequest.ImageId,
		"DesiredCfsConfig":     createRebuildWorkflowRequest.DesiredCfsConfig,
		"BootTimeoutInSeconds": createRebuildWorkflowRequest.BootTimeoutInSeconds,
	})
	if err != nil {
		return nil, err
	}
	return tmpRes.Bytes(), nil
}

func GetIufWorkflow(tmpl *template.Template, workflowFS fs.FS, req models_iuf.Session, stageIndex int) ([]byte, error) {
	// add useful helm templating func: include
	var funcMap template.FuncMap = map[string]interface{}{}
	funcMap["include"] = func(name string, data interface{}) (string, error) {
		buf := bytes.NewBuffer(nil)
		if err := tmpl.ExecuteTemplate(buf, name, data); err != nil {
			return "", err
		}
		return buf.String(), nil
	}

	// add sprig templating func
	tmpl, err := tmpl.Funcs(sprig.TxtFuncMap()).Funcs(funcMap).ParseFS(workflowFS, "*.yaml", "stages/*.yaml")
	if err != nil {
		return nil, err
	}

	var tmpRes bytes.Buffer
	err = tmpl.Execute(&tmpRes, struct {
		Products   []models_iuf.Product
		Stages     []string
		DryRun     string
		StageInput string
	}{
		Products:   req.Products,
		Stages:     []string{req.InputParameters.Stages[stageIndex]},
		DryRun:     "true",
		StageInput: "{}", //todo
	})
	if err != nil {
		return nil, err
	}
	return tmpRes.Bytes(), nil
}
