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
	_ "embed"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"regexp"
	"testing"

	"github.com/Cray-HPE/cray-nls/src/api/models/iuf"
	"github.com/Cray-HPE/cray-nls/src/utils"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	fake "k8s.io/client-go/kubernetes/fake"
	workflowmocks "github.com/argoproj/argo-workflows/v3/pkg/apiclient/workflow/mocks"
    "github.com/stretchr/testify/mock"
    wfv1 "github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
)

func TestCreateActivity(t *testing.T) {
	fakeClient := fake.NewSimpleClientset()
	mySvc := iufService{logger: utils.GetLogger(), k8sRestClientSet: fakeClient}
	var tests = []struct {
		name    string
		req     iuf.CreateActivityRequest
		wantErr bool
	}{
		{
			name:    "no name",
			req:     iuf.CreateActivityRequest{},
			wantErr: true,
		},
		{
			name:    "has name",
			req:     iuf.CreateActivityRequest{Name: "this-is-a-name"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := mySvc.CreateActivity(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("got %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestPatchActivity(t *testing.T) {
	fakeClient := fake.NewSimpleClientset()
	mySvc := iufService{logger: utils.GetLogger(), k8sRestClientSet: fakeClient}

	toPatchRequest := func(jsonStr string) iuf.PatchActivityRequest {
		var req iuf.PatchActivityRequest
		err := json.Unmarshal([]byte(jsonStr), &req)
		if err != nil {
			t.Errorf("Cannot convert to iuf.PatchActivityRequest: %s", jsonStr)
		}
		return req
	}

	type PatchTest struct {
		testName       string
		activity       iuf.Activity
		activities     []iuf.Activity // used when there are multiple starting states to test against
		expectActivity iuf.Activity
		req            iuf.PatchActivityRequest
		wantErr        bool
	}

	var tests = []PatchTest{
		{
			testName:       "Nothing to patch",
			activity:       iuf.Activity{Name: "test"},
			expectActivity: iuf.Activity{Name: "test"},
			req:            iuf.PatchActivityRequest{},
			wantErr:        false,
		},
		{
			testName: "Site parameters was previously empty, now is populated",
			activity: iuf.Activity{Name: "test"},
			expectActivity: iuf.Activity{
				Name: "test",
				SiteParameters: iuf.SiteParameters{
					Global: map[string]interface{}{
						"a": "1",
					},
					Products: map[string]map[string]interface{}{
						"cos": map[string]interface{}{
							"branch": "9.9.9-integration",
						},
					},
				},
			},
			req:     toPatchRequest(`{"site_parameters": {"global": {"a": "1"}, "products": {"cos": {"branch": "9.9.9-integration"}}}}`),
			wantErr: false,
		},
		{
			testName: "Site parameters was previously non-empty, now is populated",
			activity: iuf.Activity{
				Name: "test",
				SiteParameters: iuf.SiteParameters{
					Global: map[string]interface{}{
						"x": "99",
					},
					Products: map[string]map[string]interface{}{
						"sma": map[string]interface{}{
							"version": "sma-1.2.3",
							"branch":  "sma-1.2.3-integration",
						},
						"cos": map[string]interface{}{
							"version": "cos-2.3.4",
						},
					},
				},
			},
			expectActivity: iuf.Activity{
				Name: "test",
				SiteParameters: iuf.SiteParameters{
					Global: map[string]interface{}{
						"x": "99",
						"a": "1",
					},
					Products: map[string]map[string]interface{}{
						"sdu": map[string]interface{}{
							"version": "sdu-1.2.3",
							"branch":  "sdu-1.2.3-integration",
						},
						"cos": map[string]interface{}{
							"version": "cos-2.3.4",
							"branch":  "cos-9.9.9-integration",
						},
					},
				},
			},
			req:     toPatchRequest(`{"site_parameters": {"global": {"x": "99", "a": "1"}, "products": {"sdu": {"version": "sdu-1.2.3", "branch": "sdu-1.2.3-integration"}, "cos": {"version": "cos-2.3.4", "branch": "cos-9.9.9-integration"}}}}`),
			wantErr: false,
		},
		{
			testName: "Input parameters was previously empty, now is populated",
			activity: iuf.Activity{Name: "test"},
			expectActivity: iuf.Activity{
				Name: "test",
				InputParameters: iuf.InputParameters{
					MediaDir: "/a/b/c",
				},
			},
			req:     toPatchRequest(`{"input_parameters": {"media_dir": "/a/b/c"}}`),
			wantErr: false,
		},
		{
			testName: "Input parameters was previously non-empty, now is updated",
			activity: iuf.Activity{
				Name: "test",
				InputParameters: iuf.InputParameters{
					MediaDir:                 "/a/b/c",
					SiteParameters:           "deprecated_field",
					LimitManagementNodes:     []string{"Management_Master"},
					BootprepConfigManagement: "BootprepConfigManagement",
				},
			},
			expectActivity: iuf.Activity{
				Name: "test",
				InputParameters: iuf.InputParameters{
					MediaDir:                 "/a/b/c",
					SiteParameters:           "deprecated_field",
					LimitManagementNodes:     []string{"Management_Worker"},
					BootprepConfigManagement: "BootprepConfigManagement",
					BootprepConfigManaged:    "BootprepConfigManaged",
					Force:                    true,
				},
			},
			req:     toPatchRequest(`{"input_parameters": {"limit_management_nodes": ["Management_Worker"], "force": true, "bootprep_config_managed": "BootprepConfigManaged"}}`),
			wantErr: false,
		},
		{
			testName:       "Should be allowed to put activity state from in_progress to paused",
			activity:       iuf.Activity{Name: "test", ActivityState: "in_progress"},
			expectActivity: iuf.Activity{Name: "test", ActivityState: iuf.ActivityStatePaused},
			req:            toPatchRequest(fmt.Sprintf(`{"activity_state": "%s"}`, iuf.ActivityStatePaused)),
			wantErr:        false,
		},
		{
			testName: "Should not be allowed to put activity state from any other state to paused",
			activities: []iuf.Activity{
				iuf.Activity{Name: "test", ActivityState: iuf.ActivityStateWaitForAdmin},
				iuf.Activity{Name: "test", ActivityState: iuf.ActivityStateDebug},
				iuf.Activity{Name: "test", ActivityState: iuf.ActivityStateBlocked},
			},
			req:     toPatchRequest(fmt.Sprintf(`{"activity_state": "%s"}`, iuf.ActivityStatePaused)),
			wantErr: true,
		},
		{
			testName:       "Should be allowed to put activity state from debug to blocked",
			activity:       iuf.Activity{Name: "test", ActivityState: iuf.ActivityStateDebug},
			expectActivity: iuf.Activity{Name: "test", ActivityState: iuf.ActivityStateBlocked},
			req:            toPatchRequest(fmt.Sprintf(`{"activity_state": "%s"}`, iuf.ActivityStateBlocked)),
			wantErr:        false,
		},
		{
			testName:       "Should be allowed to put activity state from wait_for_admin to blocked",
			activity:       iuf.Activity{Name: "test", ActivityState: iuf.ActivityStateWaitForAdmin},
			expectActivity: iuf.Activity{Name: "test", ActivityState: iuf.ActivityStateBlocked},
			req:            toPatchRequest(fmt.Sprintf(`{"activity_state": "%s"}`, iuf.ActivityStateBlocked)),
			wantErr:        false,
		},
		{
			testName:       "Should be allowed to put activity state from paused to blocked",
			activity:       iuf.Activity{Name: "test", ActivityState: iuf.ActivityStatePaused},
			expectActivity: iuf.Activity{Name: "test", ActivityState: iuf.ActivityStateBlocked},
			req:            toPatchRequest(fmt.Sprintf(`{"activity_state": "%s"}`, iuf.ActivityStateBlocked)),
			wantErr:        false,
		},
		{
			testName: "Should not be allowed to put activity state from in_progress to blocked",
			activities: []iuf.Activity{
				iuf.Activity{Name: "test", ActivityState: "in_progress"},
			},
			req:     toPatchRequest(fmt.Sprintf(`{"activity_state": "%s"}`, iuf.ActivityStateBlocked)),
			wantErr: true,
		},
		{
			testName: "Should not be allowed to change the name of the activity",
			activity: iuf.Activity{
				Name: "test",
			},
			expectActivity: iuf.Activity{
				Name: "test",
			},
			req:     toPatchRequest(`{"name": "test2"}`),
			wantErr: false,
		},
		{
			testName: "Should not be allowed to change operation outputs",
			activity: iuf.Activity{
				Name:             "test",
				OperationOutputs: map[string]interface{}{"a": "b"},
			},
			expectActivity: iuf.Activity{
				Name:             "test",
				OperationOutputs: map[string]interface{}{"a": "b"},
			},
			req:     toPatchRequest(`{"operation_outputs": {"a": "c"}}`),
			wantErr: false,
		},
		{
			testName: "Should not be allowed to change products info",
			activity: iuf.Activity{
				Name: "test",
				Products: []iuf.Product{
					iuf.Product{
						Name:    "cos",
						Version: "1.2.3",
					},
					iuf.Product{
						Name:    "sdu",
						Version: "2.3.4",
					},
				},
			},
			expectActivity: iuf.Activity{
				Name: "test",
				Products: []iuf.Product{
					iuf.Product{
						Name:    "cos",
						Version: "1.2.3",
					},
					iuf.Product{
						Name:    "sdu",
						Version: "2.3.4",
					},
				},
			},
			req:     toPatchRequest(`{"products": {"cos": {"name": "cos", "version": "9.9.9"}}}`),
			wantErr: false,
		},
	}

	m1 := regexp.MustCompile(`[^a-zA-Z]`)

	runTest := func(startActivity iuf.Activity, test PatchTest, t *testing.T) bool {

		activityName := utils.GenerateName(startActivity.Name + m1.ReplaceAllString(test.testName, "-"))
		startActivity.Name = activityName
		if test.expectActivity.Name != "" {
			test.expectActivity.Name = activityName
		}

		_, err := mySvc.CreateActivity(iuf.CreateActivityRequest{
			Name: activityName,
		})
		if err != nil {
			t.Errorf("got %v, wantErr %v", err, test.wantErr)
			return false
		}

		patchedActivity, err := mySvc.PatchActivity(startActivity, test.req)
		if (err != nil) != test.wantErr {
			t.Errorf("got %v, wantErr %v", err, test.wantErr)
			return false
		}
		if test.expectActivity.Name != "" && !cmp.Equal(patchedActivity, test.expectActivity) {
			t.Errorf("Wrong object received: %s", cmp.Diff(test.expectActivity, patchedActivity))
			return false
		}
		return true
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			if len(test.activities) > 0 {
				for i, activity := range test.activities {
					activity.Name = fmt.Sprintf("%s-%v", activity.Name, i)
					runTest(activity, test, t)
				}
			} else {
				runTest(test.activity, test, t)
			}
		})
	}
}

// TestDeleteActivity verifies the basic deletion functionality of activities.
func TestDeleteActivity(t *testing.T) {
    fakeClient := fake.NewSimpleClientset()

    // Create mock workflow client
    wfServiceClientMock := &workflowmocks.WorkflowServiceClient{}
    wfServiceClientMock.On("ListWorkflows", mock.Anything, mock.Anything).Return(&wfv1.WorkflowList{}, nil)
    wfServiceClientMock.On("DeleteWorkflow", mock.Anything, mock.Anything).Return(nil, nil)
    
    // Add workflow client to service (MODIFY THIS)
    mySvc := iufService{
        logger:           utils.GetLogger(),
        k8sRestClientSet: fakeClient,
        workflowClient:   wfServiceClientMock,
    }

    var tests = []struct {
        name    string
        setup   func() string // Returns activity name to delete
        wantErr bool
    }{
        {
            name: "Successfully delete existing activity",
            setup: func() string {
                activityName := "test-delete-activity"
                _, err := mySvc.CreateActivity(iuf.CreateActivityRequest{Name: activityName})
                if err != nil {
                    t.Fatalf("Setup failed: %v", err)
                }
                return activityName
            },
            wantErr: false,
        },
        {
            name: "Delete non-existent activity should not error",
            setup: func() string {
                return "non-existent-activity"
            },
            wantErr: false,
        },
        {
            name: "Delete activity with history entries",
            setup: func() string {
                activityName := "test-with-history"
                _, err := mySvc.CreateActivity(iuf.CreateActivityRequest{Name: activityName})
                if err != nil {
                    t.Fatalf("Setup failed: %v", err)
                }
                // History is automatically created in CreateActivity
                return activityName
            },
            wantErr: false,
        },
        {
            name: "Delete activity with sessions",
            setup: func() string {
                activityName := "test-with-sessions"
                activity, err := mySvc.CreateActivity(iuf.CreateActivityRequest{Name: activityName})
                if err != nil {
                    t.Fatalf("Setup failed: %v", err)
                }
                // Create a session
                session := iuf.Session{
                    Name:            utils.GenerateName(activityName),
                    InputParameters: iuf.InputParameters{},
                }
                _, err = mySvc.CreateSession(session, "Test session", activity)
                if err != nil {
                    t.Fatalf("Setup failed creating session: %v", err)
                }
                return activityName
            },
            wantErr: false,
        },
        {
            name: "Delete activity with multiple history and sessions",
            setup: func() string {
                activityName := "test-multiple-resources"
                activity, err := mySvc.CreateActivity(iuf.CreateActivityRequest{Name: activityName})
                if err != nil {
                    t.Fatalf("Setup failed: %v", err)
                }
                // Create multiple history entries
                for i := 0; i < 3; i++ {
                    err = mySvc.CreateHistoryEntry(activityName, iuf.ActivityStateWaitForAdmin, fmt.Sprintf("Test history %d", i))
                    if err != nil {
                        t.Fatalf("Setup failed creating history: %v", err)
                    }
                }
                // Create multiple sessions
                for i := 0; i < 2; i++ {
                    session := iuf.Session{
                        Name:            utils.GenerateName(fmt.Sprintf("%s-session-%d", activityName, i)),
                        InputParameters: iuf.InputParameters{},
                    }
                    _, err = mySvc.CreateSession(session, fmt.Sprintf("Test session %d", i), activity)
                    if err != nil {
                        t.Fatalf("Setup failed creating session: %v", err)
                    }
                }
                return activityName
            },
            wantErr: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            activityName := tt.setup()
            
            success, err := mySvc.DeleteActivity(activityName)
            if (err != nil) != tt.wantErr {
                t.Errorf("DeleteActivity() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            
            if !tt.wantErr && !success {
                t.Errorf("DeleteActivity() success = %v, want true", success)
                return
            }
            
            // Verify activity was deleted
            if !tt.wantErr {
                _, err := mySvc.GetActivity(activityName)
                if err == nil {
                    t.Errorf("Activity %s still exists after deletion", activityName)
                }
                
                // Verify history entries were deleted
                historyList, err := mySvc.k8sRestClientSet.
                    CoreV1().
                    ConfigMaps(DEFAULT_NAMESPACE).
                    List(
                        context.TODO(),
                        v1.ListOptions{
                            LabelSelector: fmt.Sprintf("type=%s,%s=%s", LABEL_HISTORY, LABEL_ACTIVITY_REF, activityName),
                        },
                    )
                if err == nil && len(historyList.Items) > 0 {
                    t.Errorf("History entries still exist after deletion: found %d items", len(historyList.Items))
                }
                
                // Verify sessions were deleted
                sessionList, err := mySvc.k8sRestClientSet.
                    CoreV1().
                    ConfigMaps(DEFAULT_NAMESPACE).
                    List(
                        context.TODO(),
                        v1.ListOptions{
                            LabelSelector: fmt.Sprintf("type=%s,%s=%s", LABEL_SESSION, LABEL_ACTIVITY_REF, activityName),
                        },
                    )
                if err == nil && len(sessionList.Items) > 0 {
                    t.Errorf("Session entries still exist after deletion: found %d items", len(sessionList.Items))
                }
            }
        })
    }

	wfServiceClientMock.AssertCalled(t, "ListWorkflows", mock.Anything, mock.Anything)
}

// TestDeleteActivityRetryLogic ensures that the deletion operation handles retry scenarios gracefully.
func TestDeleteActivityRetryLogic(t *testing.T) {
    fakeClient := fake.NewSimpleClientset()

    // Add workflow client mock
    wfServiceClientMock := &workflowmocks.WorkflowServiceClient{}
    wfServiceClientMock.On("ListWorkflows", mock.Anything, mock.Anything).Return(&wfv1.WorkflowList{}, nil)
    
    mySvc := iufService{
        logger:           utils.GetLogger(),
        k8sRestClientSet: fakeClient,
        workflowClient:   wfServiceClientMock,
    }

    t.Run("Retry on not found should succeed", func(t *testing.T) {
        // Try to delete an activity that doesn't exist
        // Should succeed with warning, not error
        success, err := mySvc.DeleteActivity("non-existent")
        if err != nil {
            t.Errorf("DeleteActivity() should not error on not found, got: %v", err)
        }
        if !success {
            t.Errorf("DeleteActivity() should return success=true even for non-existent activity")
        }
    })
}

// TestDeleteActivityVerifyResourceCleanup confirms that when an activity is deleted, all associated
// resources are properly removed from the system.
func TestDeleteActivityVerifyResourceCleanup(t *testing.T) {
    fakeClient := fake.NewSimpleClientset()

    // Add workflow client mock
    wfServiceClientMock := &workflowmocks.WorkflowServiceClient{}
    wfServiceClientMock.On("ListWorkflows", mock.Anything, mock.Anything).Return(&wfv1.WorkflowList{}, nil)
    wfServiceClientMock.On("DeleteWorkflow", mock.Anything, mock.Anything).Return(nil, nil)
    
    mySvc := iufService{
        logger:           utils.GetLogger(),
        k8sRestClientSet: fakeClient,
        workflowClient:   wfServiceClientMock,
    }

    activityName := "test-cleanup-verification"
    
    // Create activity with all associated resources
    activity, err := mySvc.CreateActivity(iuf.CreateActivityRequest{Name: activityName})
    if err != nil {
        t.Fatalf("Setup failed: %v", err)
    }
    
    // Add extra history
    for i := 0; i < 3; i++ {
        err = mySvc.CreateHistoryEntry(activityName, iuf.ActivityStateInProgress, fmt.Sprintf("Entry %d", i))
        if err != nil {
            t.Fatalf("Setup failed: %v", err)
        }
    }
    
    // Add sessions
    for i := 0; i < 2; i++ {
        session := iuf.Session{
            Name:            utils.GenerateName(fmt.Sprintf("%s-session-%d", activityName, i)),
            InputParameters: iuf.InputParameters{},
        }
        _, err = mySvc.CreateSession(session, fmt.Sprintf("Test session %d", i), activity)
        if err != nil {
            t.Fatalf("Setup failed: %v", err)
        }
    }
    
    // Verify resources exist before deletion
    historyBefore, _ := mySvc.k8sRestClientSet.
        CoreV1().
        ConfigMaps(DEFAULT_NAMESPACE).
        List(
            context.TODO(),
            v1.ListOptions{
                LabelSelector: fmt.Sprintf("type=%s,%s=%s", LABEL_HISTORY, LABEL_ACTIVITY_REF, activityName),
            },
        )
    
    sessionsBefore, _ := mySvc.k8sRestClientSet.
        CoreV1().
        ConfigMaps(DEFAULT_NAMESPACE).
        List(
            context.TODO(),
            v1.ListOptions{
                LabelSelector: fmt.Sprintf("type=%s,%s=%s", LABEL_SESSION, LABEL_ACTIVITY_REF, activityName),
            },
        )
    
    if len(historyBefore.Items) < 3 {
        t.Errorf("Expected at least 3 history items before deletion, got %d", len(historyBefore.Items))
    }
    
    if len(sessionsBefore.Items) < 2 {
        t.Errorf("Expected at least 2 session items before deletion, got %d", len(sessionsBefore.Items))
    }
    
    // Delete activity
    success, err := mySvc.DeleteActivity(activityName)
    if err != nil {
        t.Fatalf("DeleteActivity failed: %v", err)
    }
    if !success {
        t.Fatalf("DeleteActivity returned success=false")
    }
    
    // Verify all resources are gone
    historyAfter, _ := mySvc.k8sRestClientSet.
        CoreV1().
        ConfigMaps(DEFAULT_NAMESPACE).
        List(
            context.TODO(),
            v1.ListOptions{
                LabelSelector: fmt.Sprintf("type=%s,%s=%s", LABEL_HISTORY, LABEL_ACTIVITY_REF, activityName),
            },
        )
    
    sessionsAfter, _ := mySvc.k8sRestClientSet.
        CoreV1().
        ConfigMaps(DEFAULT_NAMESPACE).
        List(
            context.TODO(),
            v1.ListOptions{
                LabelSelector: fmt.Sprintf("type=%s,%s=%s", LABEL_SESSION, LABEL_ACTIVITY_REF, activityName),
            },
        )
    
    if len(historyAfter.Items) > 0 {
        t.Errorf("Expected 0 history items after deletion, got %d", len(historyAfter.Items))
    }
    
    if len(sessionsAfter.Items) > 0 {
        t.Errorf("Expected 0 session items after deletion, got %d", len(sessionsAfter.Items))
    }
    
    // Verify activity itself is gone
    _, err = mySvc.GetActivity(activityName)
    if err == nil {
        t.Errorf("Activity should not exist after deletion")
    }
}
