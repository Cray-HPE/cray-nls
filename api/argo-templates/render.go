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
package argo_templates

import (
	"bytes"
	"embed"
	_ "embed"
	"fmt"
	"text/template"

	"github.com/Cray-HPE/cray-nls/utils"
	"github.com/Masterminds/sprig/v3"
)

//go:embed ncn/*
var workerRebuildWorkflow embed.FS

//go:embed base/template.argo.yaml
var argoWorkflowTemplate []byte

var validator utils.Validator = utils.NewValidator()

func GetWorkflowTemplate() []byte {
	return argoWorkflowTemplate
}

func GetWorkerRebuildWorkflow(hostnames []string) ([]byte, error) {
	err := validator.ValidateWorkerHostnames(hostnames)
	if err != nil {
		return nil, err
	}

	tmpl := template.New("worker.rebuild.yaml")

	// introduce useful helm-like templating func
	var funcMap template.FuncMap = map[string]interface{}{}
	funcMap["include"] = func(name string, data interface{}) (string, error) {
		buf := bytes.NewBuffer(nil)
		if err := tmpl.ExecuteTemplate(buf, name, data); err != nil {
			return "", err
		}
		return buf.String(), nil
	}

	// add sprig templating func
	tmpl, err = tmpl.Funcs(sprig.TxtFuncMap()).Funcs(funcMap).ParseFS(workerRebuildWorkflow, "ncn/*.yaml")
	if err != nil {
		return nil, err
	}

	var tmpRes bytes.Buffer
	err = tmpl.Execute(&tmpRes, map[string]interface{}{
		"TargetNcns": hostnames})
	if err != nil {
		return nil, err
	}
	fmt.Println(tmpRes.String())
	return tmpRes.Bytes(), nil
}
