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
	"encoding/base64"
	"os"
	"testing"

	"github.com/Cray-HPE/cray-nls/src/api/models/iuf"
	"github.com/Cray-HPE/cray-nls/src/utils"
	"github.com/alecthomas/assert"
	"github.com/google/uuid"
	fake "k8s.io/client-go/kubernetes/fake"
)

func TestProcessCreateActivityRequest(t *testing.T) {
	mySvc := iufService{logger: utils.GetLogger()}
	type wanted struct {
		err    bool
		errMsg string
	}
	var tests = []struct {
		name    string
		path    string
		tarFile string
		wanted  wanted
	}{
		{
			name: "empty",
			wanted: wanted{
				err:    true,
				errMsg: "no tarball files found: /tmp",
			},
		},
		{
			name: "invalid media dir",
			path: "asdf",
			wanted: wanted{
				err:    true,
				errMsg: "no tarball files found: /tmpasdf",
			},
		},
		{
			name:    "valid media dir with valid tar files",
			path:    "/" + uuid.NewString(),
			tarFile: "H4sIAAAAAAAAA+2VTXObMBCGfeZX6Jb2AJYEguCZ9pJLb+lMei5DsBwzBcRIIjX/vitiOw6J69QZu810nwuwH1oJ9l2C6dd89UXmc6nNtFBmcgIoEEfRcAXGV8pCPmFRGDGWMC7CCWUMjBOyOsVmxnTG5hq2co5a/yAhJbUta/mJxXEai4SLNEhpTAUPWeyBN3/qjYM0uaSXgseJ8xb7c//2yZDXAJKfnrqG03gixF79u/uR/uMI9C9OvTHHRv9aKfu7uEP+8eHeCe777/wBym7ht1rNu8L6dd6UC2ls0Od19ZYafzT/XRwTseA4/8/B8/nPA0ZjFlGaJi/M/5EX5/87x+n/JKLf4YD+QfXhSP/QTNGEnEWT//n8h2+f3cPkL1UzI99pwALqNXktZwQ6w5tLU+iytYP3s0fIt6UkVzrvyXUrNYyG5o7c9MbKmny4ur75GHjbtXgggjDyOeWcUcZ5SFMRel6hGisbO4O15qr4IbW780mb2+VsbZkWUMAD81JW9a7bPTt7pRbWQKPu+jZ9a8BGSGdktrHcSdiN1Z0cPHPZVqpfG8DSyFVnsttK3WbGKi3NbAhz3Z89rHwxhPgu5CFikMbFY7KWrTIluMr92bsxj/m6rc3uIdzzcHrfVNIw4Zs29JuiGRbdLtFn2w/kv/iWXbbPBCSPE23fQqLOf3oHq0p7TGWXtq3um6WLo8fuIvKPPzrkqrrtrDy+uDvL297AdhOvehP3xZNmgEdnLet1U2lZlO2mwzZB4J2uHfjDRRAEQRAEQRAEQRAEQRAEQRDkgV8bNSO9ACgAAA==",
			wanted: wanted{
				err: false,
			},
		},
		{
			name:    "valid media dir with invalid tar files",
			path:    "/" + uuid.NewString(),
			tarFile: "f75uUZwYT2m0IbUMHmtYRSekgS4s79cNJoN5yjkHAXx4eEJnEcExrG0T7ZkOTVKlpJp02hsJTyu+ie4iGtqp7Axz7m840IW+bPu8qXSpQdM2hqeQXabPuiqf82B77Hfrd7dP7hIYx/nn9Eo4huCP0XDT4n73z+OPT8BEpz0abkC92GW+q1RQMpwc3wStBMcWkc+skj5PMRyLUqenIjQRWGt/fktsftaciH65vbj1tvjpWNWuZ5TLQGWpOifynYAyir+UQW5piOloBhUA/NR6ibNWzv1l6LymjcyjU2LakO7BZrUCfOQLuVRRYhnYZ84twD1mRUBw4pQdaiHw1oaeFbp/O7Wtzl2ggFOnU0Gykfcl05im8pA8OluVqcFUihOUL8ce81Z/FXstHrtuzdvYeftazNhcaXzHxUXdQ5Rp/Pc/Lnx55cOlVnv9JNL5Guiq/e88nAvDShpa6S+vqIKH7yef5nlrFU8cKOl9xMNLIz8OcZbTNntP1b5jO7P7FxVYaS8GqtvBmtChiX0zZNdSMajECQDYorMPP+96/kIhe5yBuTH7jBZWsADAAA",
			wanted: wanted{
				err:    true,
				errMsg: "gzip: invalid header",
			},
		},
	}
	mySvc.env.MediaDirBase = "/tmp"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.tarFile != "" {
				data, _ := base64.StdEncoding.DecodeString(tt.tarFile)
				path := mySvc.env.MediaDirBase + tt.path
				os.MkdirAll(path, os.ModePerm)
				os.WriteFile(path+"/"+uuid.NewString()+".tar.gz", data, 0644)
			}
			activity := iuf.Activity{InputParameters: iuf.InputParameters{MediaDir: tt.path}}
			err := mySvc.processActivityInputParameters(&activity)
			if (err != nil) != tt.wanted.err {
				t.Errorf("got %v, wantErr %v", err, tt.wanted.err)
				return
			}
			if err != nil {
				assert.Equal(t, tt.wanted.errMsg, err.Error())
			}
		})
	}
}

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
