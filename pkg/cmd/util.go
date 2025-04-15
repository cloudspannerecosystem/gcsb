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
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/cloudspannerecosystem/gcsb/pkg/config"
	"github.com/olekukonko/tablewriter"
	"github.com/rcrowley/go-metrics"
)

// graceful wraps a context cancel func with a listener for OS interrupt signals
func graceful(cancelFn context.CancelFunc) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		oscall := <-c
		log.Printf("System call received. Exiting! (%s)", oscall)
		cancelFn()
	}()
}

func logTable(str *strings.Builder) {
	scanner := bufio.NewScanner(strings.NewReader(str.String()))
	for scanner.Scan() {
		// Write to stderr for console output
		log.Println(scanner.Text())
		// Write to stdout for file capture
		fmt.Println(scanner.Text())
	}
}

func summarizeMetricsAsciiTable(registry metrics.Registry) {
	// tableString := &strings.Builder{}
	// t := tablewriter.NewWriter(tableString)

	summarizeTimings(registry)
}

func summarizeTimings(registry metrics.Registry) {
	tableString := &strings.Builder{}
	t := tablewriter.NewWriter(tableString)

	t.SetHeader([]string{
		"metric",
		"count",
		"min",
		"max",
		"mean",
		"stddev",
		"median",
		"95%",
		"99%",
	})

	mtrcs := []string{
		"schema.inference",
		"run",
		"operations.read.data",
		"operations.read.time",
		"operations.write.data",
		"operations.write.time",
	}

	for _, mtrc := range mtrcs {
		mr := registry.Get(mtrc)
		if mr == nil {
			log.Println("Encountered missing metric:", mtrc)
		}

		tmr, ok := mr.(metrics.Timer)
		if !ok {
			log.Println("Encountered non-timer metric: ", mtrc)
		}

		ps := tmr.Percentiles([]float64{0.5, 0.95, 0.99})

		t.Append([]string{
			mtrc,
			fmt.Sprintf("%d", tmr.Count()),
			time.Duration(tmr.Min()).String(),
			time.Duration(tmr.Max()).String(),
			time.Duration(tmr.Mean()).String(),
			time.Duration(tmr.StdDev()).String(),
			time.Duration(ps[0]).String(),
			time.Duration(ps[1]).String(),
			time.Duration(ps[2]).String(),
		})
	}

	t.Render()
	logTable(tableString)
}

func logConfig(cfg *config.Config) {
	// Write to stderr for console output
	log.Println("Configuration:")
	log.Printf("\tProject: %s", cfg.Project)
	log.Printf("\tInstance: %s", cfg.Instance)
	log.Printf("\tDatabase: %s", cfg.Database)
	log.Printf("\tThreads: %d", cfg.Threads)
	log.Printf("\tNumConns: %d", cfg.NumConns)
	log.Printf("\tOperations:")
	log.Printf("\t\tTotal: %d", cfg.Operations.Total)
	log.Printf("\t\tRead: %d", cfg.Operations.Read)
	log.Printf("\t\tWrite: %d", cfg.Operations.Write)

	// Write to stdout for file capture
	fmt.Println("Configuration:")
	fmt.Printf("\tProject: %s\n", cfg.Project)
	fmt.Printf("\tInstance: %s\n", cfg.Instance)
	fmt.Printf("\tDatabase: %s\n", cfg.Database)
	fmt.Printf("\tThreads: %d\n", cfg.Threads)
	fmt.Printf("\tNumConns: %d\n", cfg.NumConns)
	fmt.Printf("\tOperations:\n")
	fmt.Printf("\t\tTotal: %d\n", cfg.Operations.Total)
	fmt.Printf("\t\tRead: %d\n", cfg.Operations.Read)
	fmt.Printf("\t\tWrite: %d\n", cfg.Operations.Write)
}
