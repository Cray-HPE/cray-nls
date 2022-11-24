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
			tarFile: "H4sIAEN+f2MAA+1YT2/jRBT3roQQFgcucOA0gsMCkp0Zj8eOU4qUNJG6S91sNslCORC5zqRxif+sZ9xtcoCvwZlPwkdAghOfhhknTYrVJN1s6O7S/OTReN68N2/+vN/4yXrPj5lm6ETHpmZAw0AQGQaGDsHKtgAhtAkBeW1Na2iY01oCiwYyMbYtm1g2BBCZBrEUALc2gxXIGPdSMZX+JM5W6Qm1wWBF/3QpYF6/K3jvk/eVh4riej5otsH3YAYpUz4QxRDlN1Fk+6/bDVntdJ7NXqXFr6J8XVB5sJB/7Meh7iXJiOovMi/1Ih5EVHlRgrCM9ixsOmWLkr2DYRqHdK9axbBuYqQ5yLY1EzfKWq1qWRo0nRoitTKGpL6FTbk/eOpdHlKvT9PSf3cPrOW/VeS/JR4FXG5hfWtxz/mPIQh5ENJ9ZFmO4TjIKOsYOQhD27ZVhBxw9LhWfXZw+Ph5Q7/0OE/1m/i679arptthhjtpkWb9h9Fxh9e9wycXp9917Va3MXEbLehOGrzZaeDjMYQu7JpH3Se14/MT7tZbpHXeQK3zs4lbf95QHQLawuPRySqPG98P6pve87cJS1lf2p6PdfyXfCl8/xEyFUC2N4XluOf8X37+ei/IBlqSxv3M51roRcGAMq6PvXD0ij7Eflimufz8TTS//21omwAayLDQLv+7E9wm//tbyfO/Bx/ebshF/pdb/C7KUUHl4Uz+kaJ8urjdRx7jGaP9vsfp50/bM90/lBV5Iv75F18q/flV9pmsd3njq2E5/xeZ4eveA+v4T6wi/7F4dvnfXaCY/9m2oYtDsBAiDv6f5n/WykUVLqF9F5PSycuqQL3zbdQNqjlU4iybY3GA4hX1VuWfy/m/ra//Wv4jGxnF7z80d/9/7gTilHsXNGVBHFXAj1BHOlQjL6QVICJD7VPmp0HC895vVAA6QwoOUm8MmglNPcHDM9AeM05D8MVBs/2lrs7HujGoVNWPI04jXhFj9WP/J5rKNw0kHh9WZpKSLxyoQjyko/B6t2xL+SgecCZC8nrfVYQyVWYvgn69K8kZFbPhaUbznj5NRvF4JhCSiF5mrHc6ik97jMcpZZVcTcZ5bzryo1xFkypTjZwEjxbGKU1iFoiuYLn1dZ2FfZqE7PoiZDtfvcZGlCGisQRrkR/lg86HGPfmB3QzdaW1hogwLhrycSIMU++lutYr5Zt4lmZz7xobSj246SxMbfOlC9s4TDJON3cu1/J6OzCfxK124sL/VzCIppQG4SyoUuoHyVWEXSmJ3tKs403fJDvssMO7hn8AZuukMAAcAAA=",
			wanted: wanted{
				err:    true,
				errMsg: "yaml: control characters are not allowed",
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
