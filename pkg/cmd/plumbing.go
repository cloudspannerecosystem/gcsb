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
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/davecgh/go-spew/spew"
	"github.com/cloudspannerecosystem/gcsb/pkg/config"
	"github.com/cloudspannerecosystem/gcsb/pkg/generator"
	"github.com/cloudspannerecosystem/gcsb/pkg/schema"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	plumbingSchemaInferCmd.Flags().StringVarP(&plumbingSchemaInferTableName, "table", "t", "", "table name")

	plumbingConfigDumpCmd.Flags().BoolVarP(&plumbingConfigDumpCmdValidate, "validate", "v", false, "Validate the configuration")

	plumbingTableSampleCmd.Flags().StringVarP(&plumbingSampleTable, "table", "t", "", "Table to sample")
	plumbingTableSampleCmd.Flags().IntVarP(&plumbingSampleSamples, "samples", "s", 20, "number of samples to print")

	plumbingConfigCmd.AddCommand(plumbingConfigDumpCmd)
	plumbingSchemaCmd.AddCommand(plumbingSchemaInferCmd)
	plumbingTableCmd.AddCommand(plumbingTableSampleCmd)
	plumbingCmd.AddCommand(plumbingConfigCmd, plumbingSchemaCmd, plumbingTableCmd)
	rootCmd.AddCommand(plumbingCmd)
}

var (
	plumbingSampleTable           string
	plumbingSampleSamples         int
	plumbingConfigDumpCmdValidate bool   // Validate the configuration
	plumbingSchemaInferTableName  string // table name to inspect

	plumbingCmd = &cobra.Command{
		Use:    "plumbing",
		Short:  "Plumbing commands used during development",
		Long:   `These commands are not a part of --help messages. Test things here. `,
		Hidden: true,
	}

	plumbingConfigCmd = &cobra.Command{
		Use:   "config",
		Short: "Configuration related commands",
		Long:  ``,
	}

	plumbingConfigDumpCmd = &cobra.Command{
		Use:   "dump",
		Short: "Dump the configuration",
		Long:  `Used to help test the configuration package to make sure values and flags are parsed correclty`,
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := config.NewConfig(viper.GetViper())
			if err != nil {
				log.Fatalf("unable to parse configuration: %s", err.Error())
			}

			if plumbingConfigDumpCmdValidate {
				err = cfg.Validate()
				if err != nil {
					log.Fatalf("unable to validate configuration %s", err.Error())
				}
			}

			prettyPrint(cfg)
		},
	}

	plumbingSchemaCmd = &cobra.Command{
		Use:   "schema",
		Short: "Schema related commands",
		Long:  ``,
	}

	plumbingSchemaInferCmd = &cobra.Command{
		Use:   "infer",
		Short: "Connect to the database and infer the schema",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			cfg, err := config.NewConfig(viper.GetViper())
			if err != nil {
				log.Fatalf("unable to parse configuration: %s", err.Error())
			}

			if plumbingConfigDumpCmdValidate {
				err = cfg.Validate()
				if err != nil {
					log.Fatalf("unable to validate configuration %s", err.Error())
				}
			}

			var s schema.Schema
			if plumbingSchemaInferTableName == "" {
				s, err = schema.LoadSchema(ctx, cfg)
			} else {
				s, err = schema.LoadSingleTableSchema(ctx, cfg, plumbingSchemaInferTableName)
			}
			if err != nil {
				log.Fatalf("unable to load schema: %s", err.Error())
			}
			spew.Dump(s)
		},
	}

	plumbingTableCmd = &cobra.Command{
		Use:   "table",
		Short: "Table related commands",
		Long:  ``,
	}

	plumbingTableSampleCmd = &cobra.Command{
		Use:   "sample",
		Short: "Sample rows from a table",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			if plumbingSampleTable == "" {
				log.Fatal("table name must be set (-t)")
			}

			ctx := context.Background()
			cfg, err := config.NewConfig(viper.GetViper())
			if err != nil {
				log.Fatalf("unable to parse configuration: %s", err.Error())
			}

			client, err := cfg.Client(ctx)
			if err != nil {
				log.Fatalf("error creating spanner client: %s", err.Error())
			}

			s, err := schema.LoadSchema(ctx, cfg)
			if err != nil {
				log.Fatalf("unable to load schema: %s", err.Error())
			}

			table := s.GetTable(plumbingSampleTable)
			if table == nil {
				log.Fatalf("could not find table '%s'", plumbingSampleTable)
			}

			samples, err := generator.SampleTable(cfg, ctx, client, table)
			if err != nil {
				log.Fatalf("error sampling table: %s", err.Error())
			}

			gen, err := generator.GetReadGeneratorMap(samples, table.PrimaryKeyNames())
			if err != nil {
				log.Fatalf("error getting read generator: %s", err.Error())
			}

			for i := 0; i <= plumbingSampleSamples; i++ {
				log.Printf("%+v", gen.Next())
			}
		},
	}
)

// I'm too lazy to format output of plumbing commands
func prettyPrint(i interface{}) {
	s, _ := json.MarshalIndent(i, "", "\t")
	fmt.Println(string(s))
}
