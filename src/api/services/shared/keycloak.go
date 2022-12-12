// MIT License
//
// (C) Copyright 2022 Hewlett Packard Enterprise Development LP
//
// Permission is hereby granted, free of charge, to any person obtaining a
// copy of this software and associated documentation files (the "Software"),
// to deal in the Software without restriction, including without limitation
// the rights to use, copy, modify, merge, publish, distribute, sublicense,
// and/or sell copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included
// in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
// THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR
// OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
// ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.
package services_shared

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/Cray-HPE/cray-nls/src/utils"
	k8sMetaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type KeycloakService struct {
	logger                      utils.Logger
	env                         utils.Env
	k8sService                  K8sService
	adminClientAuthClientId     string
	adminClientAuthClientSecret string
}

type oidcAuthToken struct {
	AccessToken string `json:"access_token"`
}

type NewKeycloakAccessTokenError struct {
	body string
}

func (e NewKeycloakAccessTokenError) Error() string {
	return "Could not retrieve OIDC token: " + e.body
}

func NewKeycloakService(logger utils.Logger, env utils.Env, k8sService K8sService) KeycloakService {
	var adminClientAuthClientId string
	var adminClientAuthClientSecret string

	if env.Environment != "development" {
		// only production have access to admin-client-auth
		secret, err := k8sService.Client.CoreV1().Secrets("services").Get(context.TODO(), "admin-client-auth", k8sMetaV1.GetOptions{})
		if err != nil {
			panic(err.Error())
		}

		adminClientAuthClientId = string(secret.Data["client-id"])
		adminClientAuthClientSecret = string(secret.Data["client-secret"])
	}

	return KeycloakService{
		logger:                      logger,
		env:                         env,
		k8sService:                  k8sService,
		adminClientAuthClientId:     adminClientAuthClientId,
		adminClientAuthClientSecret: adminClientAuthClientSecret,
	}
}

func (ks KeycloakService) NewKeycloakAccessToken() (string, error) {
	if ks.env.Environment == "development" {
		return "fake_dev_access_token", nil
	}

	resp, err := http.PostForm("https://"+ks.env.ApiGatewayHost+"/keycloak/realms/shasta/protocol/openid-connect/token", url.Values{
		"grant_type":    {"client_credentials"},
		"client_id":     {ks.adminClientAuthClientId},
		"client_secret": {ks.adminClientAuthClientSecret},
	})

	if err != nil {
		return "", NewKeycloakAccessTokenError{body: err.Error()}
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", NewKeycloakAccessTokenError{body: err.Error()}
	}

	if resp.StatusCode != 200 {
		return "", NewKeycloakAccessTokenError{body: fmt.Sprintf("Expected 200 response but instead got %v %v", resp.StatusCode, string(body))}
	}

	var token oidcAuthToken
	err = json.Unmarshal(body, &token)
	if err != nil {
		return "", NewKeycloakAccessTokenError{body: err.Error()}
	}

	utils.GetLogger().Infof("The access token %v", token.AccessToken)

	return token.AccessToken, nil
}
