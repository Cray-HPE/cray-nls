package argo_templates

import (
	_ "embed"
)

//go:embed worker.rebuild.argo.yaml
var argoWorkflow []byte

//go:embed base/template.argo.yaml
var argoWorkflowTemplate []byte

func GetWorkflowTemplate() []byte {
	return argoWorkflowTemplate
}

func GetWrokerRebuildWorkflow(hostname string, xName string) []byte {
	return argoWorkflow
}

func GetMasterRebuildWorkflow(hostname string, xName string) []byte {
	return nil
}

func GetStorageRebuildWorkflow(hostname string, xName string) []byte {
	return nil
}
