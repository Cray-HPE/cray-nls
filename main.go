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

// @license.name  MIT
// @license.url   https://github.com/Cray-HPE/cray-nls/blob/master/License

// @BasePath  /apis/nls

// @tag.name  NCN Lifecycle Events
// @tag.description.markdown

// @tag.name  Workflow Management
// @tag.description.markdown

// @tag.name  NCN Lifecycle Hooks
// @tag.description.markdown

func main() {
	godotenv.Load()
	logger := utils.GetLogger().GetFxLogger()
	fx.New(bootstrap.Module, fx.Logger(logger)).Run()
}
