package main

import (
	"github.com/Cray-HPE/cray-nls/bootstrap"
	_ "github.com/Cray-HPE/cray-nls/docs"
	"github.com/Cray-HPE/cray-nls/utils"
	"github.com/joho/godotenv"
	"go.uber.org/fx"
)

func main() {
	godotenv.Load()
	logger := utils.GetLogger().GetFxLogger()
	fx.New(bootstrap.Module, fx.Logger(logger)).Run()
}
