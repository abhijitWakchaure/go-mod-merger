// Package cmd ...
/*
Copyright © 2023 Abhijit Wakchaure<abhijitwakchaure.2014@gmail.com>

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
package cmd

import (
	"fmt"
	"os"

	"github.com/abhijitWakchaure/go-mod-merger/modparser"
	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go-mod-merger",
	Short: "Check mergability of multiple go.mod files and produce single merged go.mod file",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		p, _ := cmd.Flags().GetString("package")
		outputDir, _ := cmd.Flags().GetString("outputDir")
		err := modparser.Parse(p, outputDir, args)
		if err != nil {
			fmt.Printf("\nError! %s\n", err.Error())
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.Flags().StringP("package", "p", "dummy", "Package name for generated go.mod and imports.go file")
	rootCmd.Flags().StringP("outputDir", "o", "./output", "Output directory path to store generated artifacts")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		pwd, err := os.Getwd()
		cobra.CheckErr(err)

		// Search config in pwd with name "go-mod-merger.json".
		viper.AddConfigPath(pwd)
		viper.SetConfigType("json")
		viper.SetConfigName("go-mod-merger")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
