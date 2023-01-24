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
	"encoding/json"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"regexp"
	"testing"

	"github.com/Cray-HPE/cray-nls/src/api/models/iuf"
	"github.com/Cray-HPE/cray-nls/src/utils"
	fake "k8s.io/client-go/kubernetes/fake"
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
					BootprepConfigManagement: "BootprepConfigManagement",
					SiteParameters:           "deprecated_field",
				},
			},
			expectActivity: iuf.Activity{
				Name: "test",
				InputParameters: iuf.InputParameters{
					MediaDir:                 "/a/b/c",
					BootprepConfigManagement: "BootprepConfigManagement",
					BootprepConfigManaged:    "BootprepConfigManaged",
					Force:                    true,
				},
			},
			req:     toPatchRequest(`{"input_parameters": {"media_dir": "/a/b/c", "force": true, "bootprep_config_management": "BootprepConfigManagement", "bootprep_config_managed": "BootprepConfigManaged"}}`),
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
