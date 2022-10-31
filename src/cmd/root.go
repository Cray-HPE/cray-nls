/*
 *
 *  MIT License
 *
 *  (C) Copyright 2022 Hewlett Packard Enterprise Development LP
 *
 *  Permission is hereby granted, free of charge, to any person obtaining a
 *  copy of this software and associated documentation files (the "Software"),
 *  to deal in the Software without restriction, including without limitation
 *  the rights to use, copy, modify, merge, publish, distribute, sublicense,
 *  and/or sell copies of the Software, and to permit persons to whom the
 *  Software is furnished to do so, subject to the following conditions:
 *
 *  The above copyright notice and this permission notice shall be included
 *  in all copies or substantial portions of the Software.
 *
 *  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *  IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *  FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
 *  THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR
 *  OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
 *  ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
 *  OTHER DEALINGS IN THE SOFTWARE.
 *
 */

package cmd

import (
	"os"

	"github.com/Cray-HPE/cray-nls/src/bootstrap"
	"github.com/Cray-HPE/cray-nls/src/utils"
	"github.com/joho/godotenv"
	"github.com/lithammer/dedent"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "craynls",
	Short: "Cray NCN Lifecycle Service (NLS) and Installation and Upgrade Framework (IUF) operator",
	Long: dedent.Dedent(`
		This is the entry point for the Cray NCN Lifecycle Service (NLS)
		and the Installation and Upgrade Framework (IUF) operator.

		When called without any arguments, this will start the service that
		implements the NLS and IUF APIs.

		This entry point can also be used to validate IUF Product Manifest files
		using the "validate" subcommand.
		`),
	Run: func(cmd *cobra.Command, args []string) {
		godotenv.Load()
		logger := utils.GetLogger().GetFxLogger()
		fx.New(bootstrap.Module, fx.Logger(logger)).Run()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
