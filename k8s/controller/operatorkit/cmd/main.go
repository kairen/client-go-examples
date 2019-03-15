/*
Copyright © 2018 inwinSTACK.inc

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	goflag "flag"
	"fmt"
	"os"

	"github.com/kairen/simple-operator/cmd/app"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var viperWhiteList = []string{
	"v",
	"alsologtostderr",
	"log_dir",
}

const usageTemplate = `Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`

var rootCmd = &cobra.Command{
	Use:   "operator",
	Short: "operator is the core component tool.",
}

func main() {
	addcommands()
	goflag.Parse()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func addcommands() {
	rootCmd.AddCommand(app.OperatorCmd)
	rootCmd.AddCommand(app.VersionCmd)

	rootCmd.SetHelpCommand(&cobra.Command{})
	pflag.CommandLine.AddGoFlagSet(goflag.CommandLine)
	rootCmd.SetUsageTemplate(usageTemplate)
	viper.BindPFlags(rootCmd.PersistentFlags())
	setFlagsUsingViper()
}

func setFlagsUsingViper() {
	for _, config := range viperWhiteList {
		var a = pflag.Lookup(config)
		viper.SetDefault(a.Name, a.DefValue)
		if a.Changed {
			viper.Set(a.Name, a.Value.String())
		}
		a.Value.Set(viper.GetString(a.Name))
		a.Changed = true
	}
}
