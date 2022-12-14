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
	"fmt"
	"github.com/Cray-HPE/cray-nls/src/utils"
	"github.com/alecthomas/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestGetFixedValue(t *testing.T) {
	clientId, clientSecret, fakeOidcToken := "test_client_id", "test_client_secret", "fake_oidc_token"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected to have a POST request, got %s", r.Method)
		}
		if r.URL.Path != "/keycloak/realms/shasta/protocol/openid-connect/token" {
			t.Errorf("Expected to request '/keycloak/realms/shasta/protocol/openid-connect/token', got: %s", r.URL.Path)
		}
		if r.FormValue("grant_type") != "client_credentials" {
			t.Errorf("Expected to have a grant_type:client_credentials form value, got %s", r.PostForm.Get("grant_type"))
		}
		if r.FormValue("client_id") != clientId {
			t.Errorf("Expected to have a grant_type:%s form value, got %s", clientId, r.PostForm.Get("grant_type"))
		}
		if r.FormValue("client_secret") != clientSecret {
			t.Errorf("Expected to have a grant_type:%s form value, got %s", clientSecret, r.PostForm.Get("grant_type"))
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`{"access_token":"%s"}`, fakeOidcToken)))
	}))
	defer server.Close()

	serverUrl, _ := url.Parse(server.URL)
	serverHost := "http://" + serverUrl.Host

	t.Run("It calls with the right parameters", func(t *testing.T) {
		keycloakService := keycloakService{
			logger: utils.GetLogger(),
			env: utils.Env{
				Environment:   "production",
				ApiGatewayURL: serverHost,
			},
			adminClientAuthClientId:     clientId,
			adminClientAuthClientSecret: clientSecret,
		}

		token, err := keycloakService.NewKeycloakAccessToken()
		assert.NoError(t, err)
		assert.Equal(t, fakeOidcToken, token)
	})
}
