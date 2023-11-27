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
	"github.com/argoproj/argo-workflows/v3/pkg/apiclient/workflowtemplate"

	iuf "github.com/Cray-HPE/cray-nls/src/api/models/iuf"
	services_shared "github.com/Cray-HPE/cray-nls/src/api/services/shared"
	"github.com/Cray-HPE/cray-nls/src/utils"
	"github.com/argoproj/argo-workflows/v3/pkg/apiclient/workflow"
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	core_v1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	DEFAULT_NAMESPACE  = "argo"
	LABEL_ACTIVITY     = "iuf_activity"
	LABEL_HISTORY      = "iuf_history"
	LABEL_SESSION      = "iuf_session"
	LABEL_SESSION_LOCK = "iuf_session_lock"
	LABEL_ACTIVITY_REF = "iuf_activity_ref"
)

type IufService interface {
	CreateActivity(req iuf.CreateActivityRequest) (iuf.Activity, error)
	PatchActivity(activity iuf.Activity, req iuf.PatchActivityRequest) (iuf.Activity, error)
	ListActivities() ([]iuf.Activity, error)
	GetActivity(name string) (iuf.Activity, error)
	// history
	ListActivityHistory(activityName string) ([]iuf.History, error)
	HistoryRunAction(activityName string, req iuf.HistoryRunActionRequest) (iuf.Session, error)
	HistoryAbortAction(activityName string, req iuf.HistoryAbortRequest) (iuf.Session, error)
	HistoryPausedAction(activityName string, req iuf.HistoryActionRequest) (iuf.Session, error)
	HistoryResumeAction(activityName string, req iuf.HistoryActionRequest) (iuf.Session, error)
	HistoryRestartAction(activityName string, req iuf.HistoryRestartRequest) (iuf.Session, error)
	HistoryBlockedAction(activityName string, req iuf.HistoryActionRequest) (iuf.Session, error)
	GetActivityHistory(activityName string, startTime int32) (iuf.History, error)
	ReplaceHistoryComment(activityName string, startTime int32, req iuf.ReplaceHistoryCommentRequest) (iuf.History, error)
	// session
	ListSessions(activityName string) ([]iuf.Session, error)
	GetSession(sessionName string) (iuf.Session, error)
	SyncWorkflowsToSession(session *iuf.Session) error
	FindLastWorkflowForCurrentStage(session *iuf.Session) *v1alpha1.Workflow
	RestartCurrentStage(session *iuf.Session, comment string) error
	// session operator
	ConfigMapDataToSession(data string) (iuf.Session, error)
	UpdateActivityStateFromSessionState(session iuf.Session, comment string) error
	UpdateSession(session iuf.Session) error
	UpdateSessionAndActivity(session iuf.Session, comment string) error
	IsSessionLocked(session iuf.Session) bool
	LockSession(session iuf.Session) bool
	UnlockSession(session iuf.Session)
	CreateIufWorkflow(req *iuf.Session) (retWorkflow *v1alpha1.Workflow, err error, skipStage bool)
	RunNextPartialWorkflow(session *iuf.Session) (response iuf.SyncResponse, err error, sessionCompleted bool)
	RunNextStage(session *iuf.Session) (response iuf.SyncResponse, err error, sessionCompleted bool)
	ProcessOutput(session *iuf.Session, workflow *v1alpha1.Workflow) error
	GetStages() (iuf.Stages, error)
}

// IufService service layer
type iufService struct {
	logger                 utils.Logger
	workflowClient         workflow.WorkflowServiceClient
	workflowTemplateClient workflowtemplate.WorkflowTemplateServiceClient
	k8sRestClientSet       kubernetes.Interface
	keycloakService        services_shared.KeycloakService
	env                    utils.Env
}

// NewIufService creates a new Iufservice
func NewIufService(logger utils.Logger, argoService services_shared.ArgoService, k8sSvc services_shared.K8sService, keycloakService services_shared.KeycloakService, env utils.Env) IufService {

	workflowTemplateClient, err := argoService.Client.NewWorkflowTemplateServiceClient()
	if err != nil {
		panic(err.Error)
	}

	iufSvc := iufService{
		logger:                 logger,
		workflowClient:         argoService.Client.NewWorkflowServiceClient(),
		workflowTemplateClient: workflowTemplateClient,
		k8sRestClientSet:       k8sSvc.Client,
		keycloakService:        keycloakService,
		env:                    env,
	}
	return iufSvc
}

func (s iufService) iufObjectToConfigMapData(activity interface{}, name string, iufType string) (core_v1.ConfigMap, error) {
	reqBytes, err := json.Marshal(activity)
	if err != nil {
		s.logger.Errorf("iufObjectToConfigMapData.1: an error occurred while trying to convert %s of type %s to ConfigMap json %#v: %v", name, iufType, activity, err)
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
