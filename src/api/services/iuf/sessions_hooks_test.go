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
package services_iuf

import (
	"fmt"
	"github.com/Cray-HPE/cray-nls/src/api/models/iuf"
	"github.com/Cray-HPE/cray-nls/src/utils"
	"github.com/alecthomas/assert"
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetProductHookTasks(t *testing.T) {
	iufService := iufService{
		logger: utils.GetLogger(),
		env:    utils.Env{WorkerRebuildWorkflowFiles: "badname", IufInstallWorkflowFiles: "./_test_data_"},
	}

	stages, _ := iufService.GetStages()

	allTemplatesByName := map[string]bool{
		"master-host-hook-script": true,
		"worker-host-hook-script": true,
	}

	globalParamsPerProduct := map[string][]byte{
		"cos": []byte("cos_test"),
		"sdu": []byte("sdu_test"),
	}

	authToken := "fake_auth_token"

	// we only want to check the length of the preSteps and postSteps.
	// The actual validation of these is done in individual unit tests below.

	t.Run("creates correct number of product hook tasks across multiple products", func(t *testing.T) {
		preSteps, postSteps := iufService.getProductHookTasks(session, iuf.Stage{
			Name: "pre-install-check",
		}, stages, allTemplatesByName, globalParamsPerProduct, authToken)
		assert.Equal(t, 1, len(preSteps))
		assert.Equal(t, 2, len(postSteps))
	})
	t.Run("correctly ignores hooks for stages with NoHooks defined", func(t *testing.T) {
		preSteps, postSteps := iufService.getProductHookTasks(session, iuf.Stage{
			Name:    "pre-install-check",
			NoHooks: true,
		}, stages, allTemplatesByName, globalParamsPerProduct, authToken)
		assert.Equal(t, 0, len(preSteps))
		assert.Equal(t, 0, len(postSteps))
	})
	t.Run("correctly ignores hooks not populated with script_path and with missing hook templates", func(t *testing.T) {
		preSteps, postSteps := iufService.getProductHookTasks(session, iuf.Stage{
			Name: "deliver-product",
		}, stages, allTemplatesByName, globalParamsPerProduct, authToken)
		assert.Equal(t, 0, len(preSteps))
		assert.Equal(t, 0, len(postSteps))
	})
	t.Run("correctly ignores hooks with script_path that is empty", func(t *testing.T) {
		preSteps, postSteps := iufService.getProductHookTasks(session, iuf.Stage{
			Name: "prepare-images",
		}, stages, allTemplatesByName, globalParamsPerProduct, authToken)
		assert.Equal(t, 0, len(preSteps))
		assert.Equal(t, 0, len(postSteps))
	})
	t.Run("creates correct number of product hook tasks across a single product", func(t *testing.T) {
		preSteps, postSteps := iufService.getProductHookTasks(session, iuf.Stage{
			Name: "update-vcs-config",
		}, stages, allTemplatesByName, globalParamsPerProduct, authToken)
		assert.Equal(t, 0, len(preSteps))
		assert.Equal(t, 1, len(postSteps))
	})
	t.Run("correctly ignores hooks with invalid schema", func(t *testing.T) {
		preSteps, postSteps := iufService.getProductHookTasks(session, iuf.Stage{
			Name: "deploy-product",
		}, stages, allTemplatesByName, globalParamsPerProduct, authToken)
		assert.Equal(t, 0, len(preSteps))
		assert.Equal(t, 0, len(postSteps))
	})

}

func TestGetProductHooks(t *testing.T) {
	iufService := iufService{
		logger: utils.GetLogger(),
	}

	session := iuf.Session{
		Products: []iuf.Product{
			iuf.Product{
				Name:     "cos",
				Version:  "1.2.3",
				Manifest: cosManifest,
			},
			iuf.Product{
				Name:     "sdu",
				Version:  "2.3.4",
				Manifest: sduManifest,
			},
			iuf.Product{
				Name:     "incorrectSchema",
				Version:  "9.9.9",
				Manifest: incorrectSchemaManifest,
			},
			iuf.Product{
				Name:     "incorrectYamlSyntax",
				Manifest: incorrectYamlSyntaxManifest,
			},
		},
	}

	tests := map[string]map[string]iuf.ManifestStageHooks{
		iufService.getProductVersionKeyFromNameAndVersion("cos", "1.2.3"): map[string]iuf.ManifestStageHooks{
			"pre-install-check": iuf.ManifestStageHooks{
				PreHook: iuf.ManifestHookScript{
					ScriptPath:       "hooks/pre-pre-install-check.sh",
					ExecutionContext: "master_host",
				},
				PostHook: iuf.ManifestHookScript{
					ScriptPath:       "hooks/post-pre-install-check.sh",
					ExecutionContext: "worker_host",
				},
			},
			"deliver-product": iuf.ManifestStageHooks{
				PreHook: iuf.ManifestHookScript{},
				PostHook: iuf.ManifestHookScript{
					ScriptPath:       "hooks/post-deliver-product.sh",
					ExecutionContext: "storage_host",
				},
			},
		},
		iufService.getProductVersionKeyFromNameAndVersion("sdu", "2.3.4"): map[string]iuf.ManifestStageHooks{
			"pre-install-check": iuf.ManifestStageHooks{
				PreHook: iuf.ManifestHookScript{},
				PostHook: iuf.ManifestHookScript{
					ScriptPath:       "hooks/post-pre-install-check.sh",
					ExecutionContext: "worker_host",
				},
			},
			"prepare-images": iuf.ManifestStageHooks{},
			"update-vcs-config": iuf.ManifestStageHooks{
				PostHook: iuf.ManifestHookScript{
					ScriptPath:       "hooks/post-update-vcs-config.sh",
					ExecutionContext: "master_host",
				},
			},
		},
		iufService.getProductVersionKeyFromNameAndVersion("incorrectSchema", "9.9.9"): map[string]iuf.ManifestStageHooks{
			"deploy-product": iuf.ManifestStageHooks{},
		},
	}

	for productKey, productTests := range tests {
		for stageName, hooksToVerify := range productTests {
			t.Run(fmt.Sprintf("getProductHooks can parse hook script for fake manifest of %s in stage %s", productKey, stageName),
				func(t *testing.T) {

					hooks := iufService.getProductHooks(session, iuf.Stage{
						Name: stageName,
					})

					assert.NotNil(t, hooks[productKey])
					assert.Equal(t, hooksToVerify.PreHook.ScriptPath, hooks[productKey].PreHook.ScriptPath)
					assert.Equal(t, hooksToVerify.PreHook.ExecutionContext, hooks[productKey].PreHook.ExecutionContext)
					assert.Equal(t, hooksToVerify.PostHook.ScriptPath, hooks[productKey].PostHook.ScriptPath)
					assert.Equal(t, hooksToVerify.PostHook.ExecutionContext, hooks[productKey].PostHook.ExecutionContext)
				})
		}
	}

	// since incorrectSchema is never iterated over, we do a final negative test
	t.Run("getProductHooks can correctly ignore bad YAML syntax in manifest", func(t *testing.T) {
		hooks := iufService.getProductHooks(session, iuf.Stage{
			Name: "pre-install-check",
		})
		_, exists := hooks["incorrectYamlSyntax"]
		assert.False(t, exists)
	})
}

func TestExtractPathAndExecutionContext(t *testing.T) {
	iufService := iufService{
		logger: utils.GetLogger(),
	}

	cosProductManifest, _ := iufService.getProductManifestAsInterface(iuf.Product{
		Manifest: cosManifest, Name: "cos",
	})
	incorrectSchemaProduct, _ := iufService.getProductManifestAsInterface(iuf.Product{
		Manifest: incorrectSchemaManifest, Name: "incorrectSchema",
	})
	incorrectYamlSyntaxProduct, _ := iufService.getProductManifestAsInterface(iuf.Product{
		Manifest: incorrectYamlSyntaxManifest, Name: "incorrectYamlSyntax",
	})

	t.Run("extractPathAndExecutionContext can correctly parse pre and post hook stages", func(t *testing.T) {
		preHookScript := iufService.extractPathAndExecutionContext("pre-install-check", &cosProductManifest, true)
		assert.Equal(t, "hooks/pre-pre-install-check.sh", preHookScript.ScriptPath)
		assert.Equal(t, "master_host", preHookScript.ExecutionContext)

		postHookScript := iufService.extractPathAndExecutionContext("pre-install-check", &cosProductManifest, false)
		assert.Equal(t, "hooks/post-pre-install-check.sh", postHookScript.ScriptPath)
		assert.Equal(t, "worker_host", postHookScript.ExecutionContext)

		preHookScript = iufService.extractPathAndExecutionContext("deliver-product", &cosProductManifest, true)
		assert.Equal(t, "", preHookScript.ScriptPath)
		assert.Equal(t, "", preHookScript.ExecutionContext)

		postHookScript = iufService.extractPathAndExecutionContext("deliver-product", &cosProductManifest, false)
		assert.Equal(t, "hooks/post-deliver-product.sh", postHookScript.ScriptPath)
		assert.Equal(t, "storage_host", postHookScript.ExecutionContext)

		preHookScript = iufService.extractPathAndExecutionContext("non_existent_stage", &cosProductManifest, true)
		assert.Equal(t, "", preHookScript.ScriptPath)
		assert.Equal(t, "", preHookScript.ExecutionContext)

		postHookScript = iufService.extractPathAndExecutionContext("non_existent_stage", &cosProductManifest, false)
		assert.Equal(t, "", postHookScript.ScriptPath)
		assert.Equal(t, "", postHookScript.ExecutionContext)
	})

	t.Run("extractPathAndExecutionContext can correctly ignore manifests with invalid schemas", func(t *testing.T) {
		preHookScript := iufService.extractPathAndExecutionContext("deliver-product", &incorrectSchemaProduct, true)
		assert.Equal(t, "", preHookScript.ScriptPath)
		assert.Equal(t, "", preHookScript.ExecutionContext)

		postHookScript := iufService.extractPathAndExecutionContext("deliver-product", &incorrectSchemaProduct, false)
		assert.Equal(t, "", postHookScript.ScriptPath)
		assert.Equal(t, "", postHookScript.ExecutionContext)
	})

	t.Run("extractPathAndExecutionContext can correctly ignore manifests with bad YAML syntax", func(t *testing.T) {
		preHookScript := iufService.extractPathAndExecutionContext("deliver-product", &incorrectYamlSyntaxProduct, true)
		assert.Equal(t, "", preHookScript.ScriptPath)
		assert.Equal(t, "", preHookScript.ExecutionContext)

		postHookScript := iufService.extractPathAndExecutionContext("deliver-product", &incorrectYamlSyntaxProduct, false)
		assert.Equal(t, "", postHookScript.ScriptPath)
		assert.Equal(t, "", postHookScript.ExecutionContext)
	})
}

func TestCreateHookDAGTask(t *testing.T) {
	iufService := iufService{
		logger: utils.GetLogger(),
	}

	globalParamsPerProduct := map[string][]byte{
		"cos-1.2.3": []byte("cos_test"),
		"sdu-3.4.5": []byte("sdu_test"),
	}

	authToken := "fake_auth_token"

	stage := iuf.Stage{
		Name: "deliver-product",
	}

	allTemplatesByName := map[string]bool{
		"master-host-hook-script": true,
		"worker-host-hook-script": true,
	}

	hookTemplateMap := map[string]string{
		"master_host":  "master-host-hook-script",
		"worker_host":  "worker-host-hook-script",
		"storage_host": "storage-host-template",
	}

	t.Run("creates a pre-stage task properly", func(t *testing.T) {
		task, err := iufService.createHookDAGTask(true, iuf.ManifestHookScript{
			ScriptPath:       "/something/something/something/darkside",
			ExecutionContext: "master_host",
		}, iufService.getProductVersionKeyFromNameAndVersion("cos", "1.2.3"),
			session, stage, hookTemplateMap, allTemplatesByName, globalParamsPerProduct, authToken)
		assert.NoError(t, err)
		assert.True(t, strings.HasPrefix(task.Name, fmt.Sprintf("cos-1.2.3-pre-hook-%s", stage.Name)))
		assert.Equal(t, "master-host-hook-script", task.TemplateRef.Name)
		assert.Equal(t, "main", task.TemplateRef.Template)
		assert.Equal(t, v1alpha1.AnyStringPtr(authToken), task.Arguments.GetParameterByName("auth_token").Value)
		assert.Equal(t, v1alpha1.AnyStringPtr(string(globalParamsPerProduct["cos-1.2.3"])), task.Arguments.GetParameterByName("global_params").Value)
		assert.Equal(t, v1alpha1.AnyStringPtr(filepath.Join(cosOriginalLocation, "/something/something/something/darkside")), task.Arguments.GetParameterByName("script_path").Value)
	})

	t.Run("creates a post-stage task properly", func(t *testing.T) {
		task, err := iufService.createHookDAGTask(false, iuf.ManifestHookScript{
			ScriptPath:       "/something/something/something/darkside",
			ExecutionContext: "worker_host",
		}, iufService.getProductVersionKeyFromNameAndVersion("cos", "1.2.3"),
			session, stage, hookTemplateMap, allTemplatesByName, globalParamsPerProduct, authToken)
		assert.NoError(t, err)
		assert.True(t, strings.HasPrefix(task.Name, fmt.Sprintf("cos-1.2.3-post-hook-%s", stage.Name)))
		assert.Equal(t, "worker-host-hook-script", task.TemplateRef.Name)
		assert.Equal(t, "main", task.TemplateRef.Template)
		assert.Equal(t, v1alpha1.AnyStringPtr(authToken), task.Arguments.GetParameterByName("auth_token").Value)
		assert.Equal(t, v1alpha1.AnyStringPtr(string(globalParamsPerProduct["cos-1.2.3"])), task.Arguments.GetParameterByName("global_params").Value)
		assert.Equal(t, v1alpha1.AnyStringPtr(filepath.Join(cosOriginalLocation, "/something/something/something/darkside")), task.Arguments.GetParameterByName("script_path").Value)
	})

	t.Run("doesn't create a task when there is no script path", func(t *testing.T) {
		_, err := iufService.createHookDAGTask(true, iuf.ManifestHookScript{
			ScriptPath:       "",
			ExecutionContext: "master_host",
		}, "cos", session, stage, hookTemplateMap, allTemplatesByName, globalParamsPerProduct, authToken)
		assert.Error(t, err)
	})

	t.Run("doesn't create a task when there is no execution context", func(t *testing.T) {
		_, err := iufService.createHookDAGTask(true, iuf.ManifestHookScript{
			ScriptPath:       "/something/something/something/darkside",
			ExecutionContext: "",
		}, "cos", session, stage, hookTemplateMap, allTemplatesByName, globalParamsPerProduct, authToken)
		assert.Error(t, err)
	})

	t.Run("doesn't create a task when there is no Argo template in the system", func(t *testing.T) {
		_, err := iufService.createHookDAGTask(true, iuf.ManifestHookScript{
			ScriptPath:       "/something/something/something/darkside",
			ExecutionContext: "storage_host",
		}, "cos", session, stage, hookTemplateMap, allTemplatesByName, globalParamsPerProduct, authToken)
		assert.Error(t, err)
	})

	t.Run("doesn't create a task when there is no hook template defined in stages.yaml", func(t *testing.T) {
		_, err := iufService.createHookDAGTask(true, iuf.ManifestHookScript{
			ScriptPath:       "/something/something/something/darkside",
			ExecutionContext: "non_existent_execution_context",
		}, "cos", session, stage, hookTemplateMap, allTemplatesByName, globalParamsPerProduct, authToken)
		assert.Error(t, err)
	})
}

// don't change these -- you will break the tests...

const cosOriginalLocation = "/opt/cray/iuf/test-activity/cos-123"
const cosManifest = `
---
iuf_version: ^0.5.0
name: cos

hooks:
  pre_install_check:
    pre:
      script_path: hooks/pre-pre-install-check.sh
    post:
      script_path: hooks/post-pre-install-check.sh
      execution_context: worker_host
  deliver_product:
    pre:
      execution_context: worker_host
    post:
      script_path: hooks/post-deliver-product.sh
      execution_context: storage_host
`

const sduOriginalLocation = "/opt/cray/iuf/test-activity/sdu-345"
const sduManifest = `
---
iuf_version: ^0.5.0
name: sdu

hooks:
  pre_install_check:
    pre:
      execution_context: worker_host
    post:
      script_path: hooks/post-pre-install-check.sh
      execution_context: worker_host
  prepare_images:
    post:
      script_path: ""
  update_vcs_config:
    post:
      script_path: hooks/post-update-vcs-config.sh
`

const sduManifest_alt = `
---
iuf_version: ^0.5.0
name: sdu

hooks:
  pre_install_check:
    pre:
      script_path: hooks/pre-pre-install-check.sh
      execution_context: worker_host
    post:
      script_path: hooks/post-pre-install-check.sh
      execution_context: worker_host
  prepare_images:
    post:
      script_path: ""
  update_vcs_config:
    post:
      script_path: hooks/post-update-vcs-config.sh
`

const incorrectSchemaOriginalLocation = "/opt/cray/iuf/test-activity/incorrectSchema-345"
const incorrectSchemaManifest = `
---
iuf_version: ^0.5.0
name: incorrectSchema

hooks:
  deploy_product: ""
`

const incorrectYamlSyntaxOriginalLocation = "/opt/cray/iuf/test-activity/incorrectYamlSyntax-345"
const incorrectYamlSyntaxManifest = `
---
iuf_version: ^0.5.0
name: incorrectYamlSyntax

hooks:}
`

var session = iuf.Session{
	Products: []iuf.Product{
		iuf.Product{
			Name:             "cos",
			Version:          "1.2.3",
			Manifest:         cosManifest,
			OriginalLocation: cosOriginalLocation,
		},
		iuf.Product{
			Name:             "sdu",
			Version:          "3.4.5",
			Manifest:         sduManifest,
			OriginalLocation: sduOriginalLocation,
		},
		iuf.Product{
			Name:             "incorrectSchema",
			Version:          "9.9.9",
			Manifest:         incorrectSchemaManifest,
			OriginalLocation: incorrectSchemaOriginalLocation,
		},
		iuf.Product{
			Name:             "incorrectYamlSyntax",
			Manifest:         incorrectYamlSyntaxManifest,
			OriginalLocation: incorrectYamlSyntaxOriginalLocation,
		},
	},
}
