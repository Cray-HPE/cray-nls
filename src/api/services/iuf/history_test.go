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
	"reflect"
	"testing"
	"time"

	iuf "github.com/Cray-HPE/cray-nls/src/api/models/iuf"
	"github.com/Cray-HPE/cray-nls/src/utils"
	"github.com/alecthomas/assert"
	"github.com/google/uuid"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	fake "k8s.io/client-go/kubernetes/fake"
)

func TestGetActivityHistory(t *testing.T) {
	name := uuid.NewString()
	time := int32(time.Now().UnixMilli())
	iufHistory := iuf.History{
		ActivityState: "",
		SessionName:   name,
		StartTime:     time,
		Comment:       "",
		Name:          name,
	}
	reqBytes, _ := json.Marshal(iufHistory)
	configmap := v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: DEFAULT_NAMESPACE,
			Labels: map[string]string{
				"type":             LABEL_HISTORY,
				LABEL_ACTIVITY_REF: name,
			},
		},
		Data: map[string]string{LABEL_HISTORY: string(reqBytes)},
	}
	fakeClient := fake.NewSimpleClientset(&configmap)
	mySvc := iufService{logger: utils.GetLogger(), k8sRestClientSet: fakeClient}
	var tests = []struct {
		name         string
		activityName string
		startTime    int32
		wantErr      bool
		wantedRes    iuf.History
	}{
		{
			name:         "activity doesn't exist",
			activityName: "asdf",
			startTime:    0,
			wantErr:      false,
			wantedRes:    iuf.History{},
		},
		{
			name:         "activity exists but history doesn't",
			activityName: name,
			startTime:    0,
			wantErr:      false,
			wantedRes:    iuf.History{},
		},
		{
			name:         "activity exists and history exists",
			activityName: name,
			startTime:    time,
			wantErr:      false,
			wantedRes:    iufHistory,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			history, err := mySvc.GetActivityHistory(tt.activityName, tt.startTime)
			if (err != nil) != tt.wantErr {
				t.Errorf("got %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.True(t, reflect.DeepEqual(history, tt.wantedRes))
		})
	}
}
