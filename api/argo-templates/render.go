package argo_templates

import (
	"bytes"
	_ "embed"
	"text/template"

	"github.com/Cray-HPE/cray-nls/utils"
)

//go:embed worker.rebuild.argo.yaml
var argoWorkflow []byte

//go:embed base/template.argo.yaml
var argoWorkflowTemplate []byte

var validator utils.Validator = utils.NewValidator()

func GetWorkflowTemplate() []byte {
	return argoWorkflowTemplate
}

func GetWrokerRebuildWorkflow(hostname string, xName string) ([]byte, error) {
	err := validator.ValidateWorkerHostname(hostname)
	if err != nil {
		return nil, err
	}

	tmpl := template.New("render")
	tmpl, _ = tmpl.Parse(string(argoWorkflow))
	var tmpRes bytes.Buffer
	err = tmpl.Execute(&tmpRes, map[string]interface{}{
		"TargetNcn":  hostname,
		"TargetNcns": []string{hostname}})
	if err != nil {
		return nil, err
	}
	return tmpRes.Bytes(), nil
}

func GetMasterRebuildWorkflow(hostname string, xName string) []byte {
	return nil
}

func GetStorageRebuildWorkflow(hostname string, xName string) []byte {
	return nil
}
