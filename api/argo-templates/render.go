package argo_templates

import (
	_ "embed"
)

//go:embed worker.rebuild.argo.yaml
var ArgoWorkflow []byte

//go:embed base/template.argo.yaml
var ArgoWorkflowTemplate []byte
