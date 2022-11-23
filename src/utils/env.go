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
package utils

import (
	"github.com/spf13/viper"
)

// Env has environment stored
type Env struct {
	ServerPort                  string `mapstructure:"SERVER_PORT"`
	Environment                 string `mapstructure:"ENV"`
	ArgoToken                   string `mapstructure:"ARGO_TOKEN"`
	ArgoServerURL               string `mapstructure:"ARGO_SERVER_URL"`
	WorkerRebuildWorkflowFiles  string `mapstructure:"WORKER_REBUILD_WORKFLOW_FILES"`
	StorageRebuildWorkflowFiles string `mapstructure:"STORAGE_REBUILD_WORKFLOW_FILES"`
	IufInstallWorkflowFiles     string `mapstructure:"IUF_INSTALL_WORKFLOW_FILES"`
	MediaDirBase                string `mapstructure:"MEDIA_DIR_BASE"`
}

// NewEnv creates a new environment
func NewEnv(log Logger) Env {

	env := Env{}
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("☠️ cannot read configuration")
	}

	err = viper.Unmarshal(&env)
	if err != nil {
		log.Fatal("☠️ environment can't be loaded: ", err)
	}

	// for the intiial realse, this is hard coded
	if len(env.MediaDirBase) == 0 {
		env.MediaDirBase = "/opt/cray/iuf"
	}

	log.Infof("%+v \n", env)
	return env
}
