// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

var (
	cfgFile  string
	project  string
	instance string
	database string

	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:   "gcsb",
		Short: "Like YCSB but for spanner",
		Long:  ``,
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Config file flag
	flags := rootCmd.PersistentFlags()
	flags.StringVar(&cfgFile, "config", "", "config file (default is ./gcsb.yaml)")

	flags.StringVarP(&project, "project", "p", "", "GCP Project")
	viper.BindPFlag("project", flags.Lookup("project")) // bind flag to config

	flags.StringVarP(&instance, "instance", "i", "", "Spanner Instance")
	viper.BindPFlag("instance", flags.Lookup("instance"))

	flags.StringVarP(&database, "database", "d", "", "Spanner Database")
	viper.BindPFlag("database", flags.Lookup("database"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in "." directory with name gcsb.yaml
		viper.AddConfigPath(".")
		viper.SetConfigName("gcsb")
	}

	viper.AutomaticEnv() // read in environment variables that match
	viper.SetEnvPrefix("GCSB")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// err := viper.ReadInConfig()
	// if err != nil {
	// 	log.Fatalf("error reading config: %s", err.Error())
	// }

	viper.ReadInConfig() // Ignore errors here. We don't want to exit if no config file is found
}
