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

//go:generate mockgen -destination=../mocks/services/iuf.go -package=mocks -source=iuf.go

import (
	_ "embed"
	"encoding/json"
	"os"

	iuf "github.com/Cray-HPE/cray-nls/src/api/models/iuf"
	services_shared "github.com/Cray-HPE/cray-nls/src/api/services/shared"
	"github.com/Cray-HPE/cray-nls/src/utils"
	"github.com/argoproj/argo-workflows/v3/pkg/apiclient/workflow"
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"gopkg.in/yaml.v2"
	core_v1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	DEFAULT_NAMESPACE  = "argo"
	LABEL_ACTIVITY     = "iuf_activity"
	LABEL_HISTORY      = "iuf_history"
	LABEL_SESSION      = "iuf_session"
	LABEL_ACTIVITY_REF = "iuf_activity_ref"
)

type IufService interface {
	CreateActivity(req iuf.CreateActivityRequest) (iuf.Activity, error)
	ListActivities() ([]iuf.Activity, error)
	GetActivity(name string) (iuf.Activity, error)
	// history
	ListActivityHistory(activityName string) ([]iuf.History, error)
	HistoryRunAction(activityName string, req iuf.HistoryRunActionRequest) (iuf.Session, error)
	GetActivityHistory(activityName string, startTime int32) (iuf.History, error)
	ReplaceHistoryComment(activityName string, startTime int32, req iuf.ReplaceHistoryCommentRequest) (iuf.History, error)
	// session
	ListSessions(activityName string) ([]iuf.Session, error)
	GetSession(sessionName string) (iuf.Session, error)
	// session operator
	ConfigMapDataToSession(data string) (iuf.Session, error)
	UpdateActivityStateFromSessionState(session iuf.Session) error
	UpdateSession(session iuf.Session) error
	CreateIufWorkflow(req iuf.Session) (*v1alpha1.Workflow, error)
	RunNextStage(session *iuf.Session) (iuf.SyncResponse, error)
	ProcessOutput(session iuf.Session, workflow *v1alpha1.Workflow) error
}

// IufService service layer
type iufService struct {
	logger           utils.Logger
	workflowCient    workflow.WorkflowServiceClient
	k8sRestClientSet kubernetes.Interface
	env              utils.Env
}

// NewIufService creates a new Iufservice
func NewIufService(logger utils.Logger, argoService services_shared.ArgoService, k8sSvc services_shared.K8sService, env utils.Env) IufService {

	iufSvc := iufService{
		logger:           logger,
		workflowCient:    argoService.Client.NewWorkflowServiceClient(),
		k8sRestClientSet: k8sSvc.Client,
		env:              env,
	}
	return iufSvc
}

func (s iufService) iufObjectToConfigMapData(activity interface{}, name string, iufType string) (core_v1.ConfigMap, error) {
	reqBytes, err := json.Marshal(activity)
	if err != nil {
		s.logger.Error(err)
		return core_v1.ConfigMap{}, err
	}
	res := core_v1.ConfigMap{
		ObjectMeta: v1.ObjectMeta{
			Name: name,
			Labels: map[string]string{
				"type": iufType,
			},
		},
		Data: map[string]string{iufType: string(reqBytes)},
	}
	return res, nil
}

func (s iufService) getStages() (iuf.Stages, error) {
	stagesBytes, _ := os.ReadFile(s.env.IufInstallWorkflowFiles + "/stages.yaml")
	var stages iuf.Stages
	err := yaml.Unmarshal(stagesBytes, &stages)
	if err != nil {
		s.logger.Error(err)
	}
	return stages, err
}
