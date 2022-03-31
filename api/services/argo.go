package services

import (
	"context"

	"github.com/Cray-HPE/cray-nls/utils"
	"github.com/argoproj/argo-workflows/v3/pkg/apiclient"
)

type ArgoService struct {
	Context context.Context
	Client  apiclient.Client
}

func NewArgoService(env utils.Env) ArgoService {
	var argoOps apiclient.Opts = apiclient.Opts{
		ArgoServerOpts: apiclient.ArgoServerOpts{
			URL:                "localhost:2746",
			InsecureSkipVerify: true,
			Secure:             true,
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
