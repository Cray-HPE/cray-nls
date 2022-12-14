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
