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
	"testing"

	iuf "github.com/Cray-HPE/cray-nls/src/api/models/iuf"
	"github.com/Cray-HPE/cray-nls/src/utils"
	"github.com/alecthomas/assert"
	workflowmocks "github.com/argoproj/argo-workflows/v3/pkg/apiclient/workflow/mocks"
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/stretchr/testify/mock"
)

func TestCreateIufWorkflow(t *testing.T) {

	t.Run("It can create a new iuf workflow", func(t *testing.T) {
		// setup mocks
		wfServiceClientMock := &workflowmocks.WorkflowServiceClient{}
		wfServiceClientMock.On(
			"CreateWorkflow",
			mock.Anything,
			mock.Anything,
		).Return(new(v1alpha1.Workflow), nil)

		workflowSvc := iufService{
			logger:        utils.GetLogger(),
			workflowCient: wfServiceClientMock,
			env:           utils.Env{WorkerRebuildWorkflowFiles: "badname"},
		}
		_, err := workflowSvc.CreateIufWorkflow(iuf.Session{}, 0)

		// we don't actually test the template render/upload
		// this is tested in the render package
		assert.Contains(t, err.Error(), "template: pattern matches no files: `*.yaml`")
	})
}
