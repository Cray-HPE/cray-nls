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
package services_shared

import (
	"context"

	"github.com/Cray-HPE/cray-nls/src/utils"
	"github.com/argoproj/argo-workflows/v3/pkg/apiclient"
)

type ArgoService struct {
	Context context.Context
	Client  apiclient.Client
}

func NewArgoService(env utils.Env) ArgoService {
	var argoOps apiclient.Opts = apiclient.Opts{
		ArgoServerOpts: apiclient.ArgoServerOpts{
			URL:                env.ArgoServerURL,
			InsecureSkipVerify: true,
			Secure:             false,
			HTTP1:              true,
		},
		AuthSupplier: func() string {
			return env.ArgoToken
		},
	}
	ctx, client, _ := apiclient.NewClientFromOpts(argoOps)
	return ArgoService{
		Context: ctx,
		Client:  client,
	}
}
