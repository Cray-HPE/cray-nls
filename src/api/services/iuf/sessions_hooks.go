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

// Package services_iuf
// Various session-focused functionality for dealing with IUF hooks
package services_iuf

import (
	"encoding/json"
	"fmt"
	"github.com/Cray-HPE/cray-nls/src/api/models/iuf"
	"github.com/Cray-HPE/cray-nls/src/utils"
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/oliveagle/jsonpath"
	"path/filepath"
	"sigs.k8s.io/yaml"
	"strings"
)

// gets a map of productName vs hook task for pre-stage and post-stage hooks.
func (s iufService) getProductHookTasks(session iuf.Session, stage iuf.Stage, stages iuf.Stages,
	prevStepsCompleted map[string]map[string]string,
	allTemplatesByName map[string]bool,
	workflowParamNamesGlobalParamsPerProduct map[string]string, workflowParamNameAuthToken string) (preSteps map[string]v1alpha1.DAGTask,
	postSteps map[string]v1alpha1.DAGTask) {

	if stage.NoHooks {
		return preSteps, postSteps
	}

	hooks := s.getProductHooks(session, stage)

	preSteps = map[string]v1alpha1.DAGTask{}
	postSteps = map[string]v1alpha1.DAGTask{}

	for productKey, productHooks := range hooks {
		hook := productHooks.PreHook
		if hook.ScriptPath != "" {
			task, err := s.createHookDAGTask(true, hook, productKey, session, stage, prevStepsCompleted, stages.Hooks, allTemplatesByName, workflowParamNamesGlobalParamsPerProduct, workflowParamNameAuthToken)
			if err == nil {
				if task.Name != "" { // empty name means we are skipping this task
					preSteps[productKey] = task
				}
			} else {
				s.logger.Error(err)
			}
		}

		hook = productHooks.PostHook
		if hook.ScriptPath != "" {
			task, err := s.createHookDAGTask(false, hook, productKey, session, stage, prevStepsCompleted, stages.Hooks, allTemplatesByName, workflowParamNamesGlobalParamsPerProduct, workflowParamNameAuthToken)
			if err == nil {
				if task.Name != "" { // empty name means we are skipping this task
					postSteps[productKey] = task
				}
			} else {
				s.logger.Error(err)
			}
		}
	}

	return preSteps, postSteps
}

// Returns a map of the name of the product vs its hooks for the given stage
func (s iufService) getProductHooks(session iuf.Session, stage iuf.Stage) map[string]iuf.ManifestStageHooks {
	ret := map[string]iuf.ManifestStageHooks{}
	for _, product := range session.Products {
		manifest, err := s.getProductManifestAsInterface(product)
		if err != nil {
			continue
		}

		preHook := s.extractPathAndExecutionContext(stage.Name, &manifest, true)
		postHook := s.extractPathAndExecutionContext(stage.Name, &manifest, false)

		if preHook.ScriptPath != "" || postHook.ScriptPath != "" {
			ret[s.getProductVersionKey(product)] = iuf.ManifestStageHooks{
				PreHook:  preHook,
				PostHook: postHook,
			}
		}
	}

	return ret
}

// Returns ManifestHookScript for the given stageName and parsed IUF manifest
func (s iufService) extractPathAndExecutionContext(stageName string, manifest *interface{}, pre bool) iuf.ManifestHookScript {
	stageName = strings.Replace(stageName, "-", "_", -1)

	preOrPost := "pre"
	if !pre {
		preOrPost = "post"
	}

	ret := iuf.ManifestHookScript{
		ScriptPath:       "",
		ExecutionContext: "master_host",
	}

	jsonPath := fmt.Sprintf("$.hooks.%s.%s.script_path", stageName, preOrPost)
	scriptPathInterface, err := jsonpath.JsonPathLookup(*manifest, jsonPath)
	if err != nil || scriptPathInterface == nil {
		return iuf.ManifestHookScript{}
	}
	ret.ScriptPath = fmt.Sprintf("%s", scriptPathInterface)

	jsonPath = fmt.Sprintf("$.hooks.%s.%s.execution_context", stageName, preOrPost)
	executionContextInterface, err := jsonpath.JsonPathLookup(*manifest, jsonPath)
	if err == nil && executionContextInterface != nil {
		ret.ExecutionContext = fmt.Sprintf("%s", executionContextInterface)
	}

	return ret
}

// creates a DAG task  for a hook
func (s iufService) createHookDAGTask(pre bool, hook iuf.ManifestHookScript, productKey string, session iuf.Session, stage iuf.Stage,
	prevStepsCompleted map[string]map[string]string,
	hookTemplateMap map[string]string, allTemplatesByName map[string]bool,
	workflowParamNamesGlobalParamsPerProduct map[string]string, workflowParamNameAuthToken string) (v1alpha1.DAGTask, error) {

	// find the original location
	var originalLocation string
	for _, product := range session.Products {
		if s.getProductVersionKey(product) == productKey {
			originalLocation = product.OriginalLocation
			break
		}
	}

	if hook.ScriptPath == "" || hook.ExecutionContext == "" || originalLocation == "" {
		return v1alpha1.DAGTask{}, utils.GenericError{Message: fmt.Sprintf("No valid hook script found for product %s in stage %s.", productKey, stage.Name)}
	}

	originalLocation = filepath.Clean(originalLocation)
	filePath := filepath.Clean(filepath.Join(originalLocation, hook.ScriptPath))
	if strings.Index(filePath, originalLocation) != 0 {
		// possible hack attempt ... reading a parent directory through relative paths.
		return v1alpha1.DAGTask{}, utils.GenericError{Message: fmt.Sprintf("Bad hook script path %v found for product %s in stage %s.", filePath, productKey, stage.Name)}
	}

	preOrPost := "-pre-hook-"
	if !pre {
		preOrPost = "-post-hook-"
	}

	if prevStepsCompleted[productKey][preOrPost+stage.Name] != "" {
		// this hook script was already completed in a previous run, so skip it.
		return v1alpha1.DAGTask{}, nil
	}

	hookTemplateName := hookTemplateMap[hook.ExecutionContext]
	templateExists := allTemplatesByName[hookTemplateName]
	if hookTemplateName == "" || !templateExists {
		return v1alpha1.DAGTask{}, utils.GenericError{Message: fmt.Sprintf("The template %s is not available in Argo.", hookTemplateName)}
	}

	templateRef := v1alpha1.TemplateRef{
		Name:     hookTemplateName,
		Template: "main",
	}

	task := v1alpha1.DAGTask{
		Name:        utils.GenerateName(productKey + preOrPost + stage.Name),
		TemplateRef: &templateRef,
		Arguments: v1alpha1.Arguments{
			Parameters: []v1alpha1.Parameter{
				{
					Name:  "auth_token",
					Value: v1alpha1.AnyStringPtr(fmt.Sprintf("{{workflow.parameters.%s}}", workflowParamNameAuthToken)),
				},
				{
					Name:  "global_params",
					Value: v1alpha1.AnyStringPtr(fmt.Sprintf("{{workflow.parameters.%s}}", workflowParamNamesGlobalParamsPerProduct[productKey])),
				},
				{
					Name:  "script_path",
					Value: v1alpha1.AnyStringPtr(filePath),
				},
			},
		},
	}
	return task, nil
}

// Gets the product manifest as an interface
func (s iufService) getProductManifestAsInterface(product iuf.Product) (i interface{}, e error) {
	jsonBytes, err := yaml.YAMLToJSON([]byte(product.Manifest))
	if err != nil {
		return i, err
	}

	var manifest interface{}
	err = json.Unmarshal([]byte(jsonBytes), &manifest)
	if err != nil {
		return i, err
	}

	return manifest, nil
}
