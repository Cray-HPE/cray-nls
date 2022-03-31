package services

import (
	"context"

	"github.com/argoproj/argo-workflows/v3/pkg/apiclient"
)

type ArgoService struct {
	Context context.Context
	Client  apiclient.Client
}

func NewArgoService() ArgoService {
	var argoOps apiclient.Opts = apiclient.Opts{
		ArgoServerOpts: apiclient.ArgoServerOpts{
			URL:                "localhost:2746",
			InsecureSkipVerify: true,
			Secure:             true,
			HTTP1:              true,
		},
		AuthSupplier: func() string {
			return `ZXlKaGJHY2lPaUpTVXpJMU5pSXNJbXRwWkNJNklrRkZNMFkwUjNjMWRWZE5iSFJqZDBoeE4wSmtXRkJJVTNSTWNYQjNia1pHY2tRNU5IUnpTbFpVZEdjaWZRLmV5SnBjM01pT2lKcmRXSmxjbTVsZEdWekwzTmxjblpwWTJWaFkyTnZkVzUwSWl3aWEzVmlaWEp1WlhSbGN5NXBieTl6WlhKMmFXTmxZV05qYjNWdWRDOXVZVzFsYzNCaFkyVWlPaUpoY21kdklpd2lhM1ZpWlhKdVpYUmxjeTVwYnk5elpYSjJhV05sWVdOamIzVnVkQzl6WldOeVpYUXVibUZ0WlNJNkltRnlaMjh0ZEc5clpXNHROamN5ZDNvaUxDSnJkV0psY201bGRHVnpMbWx2TDNObGNuWnBZMlZoWTJOdmRXNTBMM05sY25acFkyVXRZV05qYjNWdWRDNXVZVzFsSWpvaVlYSm5ieUlzSW10MVltVnlibVYwWlhNdWFXOHZjMlZ5ZG1salpXRmpZMjkxYm5RdmMyVnlkbWxqWlMxaFkyTnZkVzUwTG5WcFpDSTZJbU0wTURFNVlXSmtMVE15TkRVdE5EWTBNUzA0WVdOa0xXVmhZbU16T0dVNVpqRXlNU0lzSW5OMVlpSTZJbk41YzNSbGJUcHpaWEoyYVdObFlXTmpiM1Z1ZERwaGNtZHZPbUZ5WjI4aWZRLkNLaTdkMkVRQzNmS0tUdzdqVFE3ajB4VUV4eTM0Z3JVZ2c5elQ3amhTQ0xaZ3d1SngxVXpYbURaOWhUTDRmQ0pkSEw0SElTZlNYNE1ETUxmQ2lWN1U3dXEyMGRWQ3Q2Nlh0WWhObXlpTEVmU05IRE5kSWVuUU96amc0MEZSUkMzcjd4MmptVDNuejJmc2lRTGwwZDJDaUxTcjBYZDY4MC0xc2gwXzZxSlRfeS1yM2F4R1puZHV6SzJFbXBYU0J4VTJ5ZlFDQ2JHWFF2MmRjaXBDSmdEeEtzaXkzMnc5STltWnhCb1FtOHdEU2lEZTIteGxPOThoLS1GM0U2WEJwUHE1OW85UEYyVVhNYXdKSXVjX3kyUERUdzJzSXhwbkpNZXVGaWRzWXN4dVUtdFJIMkNyTFJhUlIxMzJhMVNCazJJcTRfZUM4WmNZRFJDMHF6VzRFbUJfUQ==`
		},
	}
	ctx, client, _ := apiclient.NewClientFromOpts(argoOps)
	return ArgoService{
		Context: ctx,
		Client:  client,
	}
}
