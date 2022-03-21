package main

import (
	"github.com/Cray-HPE/cray-nls/bootstrap"
	_ "github.com/Cray-HPE/cray-nls/docs"
	"github.com/Cray-HPE/cray-nls/utils"
	"github.com/joho/godotenv"
	"go.uber.org/fx"
)

// @title    NCN Lifecycle Management API
// @version  1.0
// @description.markdown

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath  /api/v1

func main() {
	godotenv.Load()
	logger := utils.GetLogger().GetFxLogger()
	fx.New(bootstrap.Module, fx.Logger(logger)).Run()
}
