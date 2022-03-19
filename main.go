package main

import (
	"github.com/Cray-HPE/cray-nls/bootstrap"
	_ "github.com/Cray-HPE/cray-nls/docs"
	"github.com/Cray-HPE/cray-nls/utils"
	"github.com/joho/godotenv"
	"go.uber.org/fx"
)

// @title           NCN Lifecycle Management API
// @version         1.0
// @description     <h1>This is a description</h1>
// @description 	<h3>Multiple line description</h3>
// @description 	##### Markdown is here
// @description 	`cmd`

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth

// @securityDefinitions.apikey  ApiKeyAuth
// @in                          header
// @name                        Authorization

// @securitydefinitions.oauth2.application  OAuth2Application
// @tokenUrl                                https://example.com/oauth/token
// @scope.write                             Grants write access
// @scope.admin                             Grants read and write access to administrative information

// @securitydefinitions.oauth2.implicit  OAuth2Implicit
// @authorizationUrl                     https://example.com/oauth/authorize
// @scope.write                          Grants write access
// @scope.admin                          Grants read and write access to administrative information

// @securitydefinitions.oauth2.password  OAuth2Password
// @tokenUrl                             https://example.com/oauth/token
// @scope.read                           Grants read access
// @scope.write                          Grants write access
// @scope.admin                          Grants read and write access to administrative information

// @securitydefinitions.oauth2.accessCode  OAuth2AccessCode
// @tokenUrl                               https://example.com/oauth/token
// @authorizationUrl                       https://example.com/oauth/authorize
// @scope.admin                            Grants read and write access to administrative information

func main() {
	godotenv.Load()
	logger := utils.GetLogger().GetFxLogger()
	fx.New(bootstrap.Module, fx.Logger(logger)).Run()
}
