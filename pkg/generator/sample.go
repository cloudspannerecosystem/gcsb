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

package generator

import (
	"context"
	"fmt"
	"math/big"
	"math/rand"
	"time"

	"cloud.google.com/go/civil"
	"cloud.google.com/go/spanner"
	"cloud.google.com/go/spanner/spansql"
	"github.com/cloudspannerecosystem/gcsb/pkg/config"
	"github.com/cloudspannerecosystem/gcsb/pkg/generator/sample"
	"github.com/cloudspannerecosystem/gcsb/pkg/schema"
)

func SampleTable(cfg *config.Config, ctx context.Context, client *spanner.Client, table schema.Table) (map[string]interface{}, error) {
	// Get primary keys for table
	pkeys := table.PrimaryKeys()
	if pkeys.Len() <= 0 {
		return nil, fmt.Errorf("cannot find primary key(s) for table '%s'", table.Name())
	}

	ret := make(map[string]interface{}, pkeys.Len())
	for pkeys.HasNext() {
		pkey := pkeys.GetNext()
		pt := pkey.Type().Base

		switch pt {
		case spansql.Bool:
			ret[pkey.Name()] = make([]bool, 0)
		case spansql.String:
			ret[pkey.Name()] = make([]string, 0)
		case spansql.Int64:
			ret[pkey.Name()] = make([]int64, 0)
		case spansql.Float64:
			ret[pkey.Name()] = make([]float64, 0)
		case spansql.Bytes:
			ret[pkey.Name()] = make([][]byte, 0)
		case spansql.Timestamp:
			ret[pkey.Name()] = make([]time.Time, 0)
		case spansql.Date:
			ret[pkey.Name()] = make([]civil.Date, 0)
		case spansql.Numeric:
			ret[pkey.Name()] = make([]*big.Rat, 0)
		case spansql.JSON:
			ret[pkey.Name()] = make([]map[string]interface{}, 0) // TODO: This needs to be spanner.NullJSON
		}
	}

	pkeys.ResetIterator()

	stmt, err := table.TableSample(cfg.Operations.SampleSize)
	if err != nil {
		return nil, err
	}

	iter := client.Single().Query(ctx, spanner.NewStatement(stmt))
	err = iter.Do(func(r *spanner.Row) error {
		defer pkeys.ResetIterator()

		for pkeys.HasNext() {
			pkey := pkeys.GetNext()
			pt := pkey.Type().Base

			switch pt {
			case spansql.Bool:
				var val bool
				err := r.ColumnByName(pkey.Name(), &val)
				if err != nil {
					return fmt.Errorf("unable to read row value: %s", err.Error())
				}

				arr := ret[pkey.Name()].([]bool)
				ret[pkey.Name()] = append(arr, val)
			case spansql.String:
				var val string
				err := r.ColumnByName(pkey.Name(), &val)
				if err != nil {
					return fmt.Errorf("unable to read row value: %s", err.Error())
				}

				arr := ret[pkey.Name()].([]string)
				ret[pkey.Name()] = append(arr, val)
			case spansql.Int64:
				var val int64
				err := r.ColumnByName(pkey.Name(), &val)
				if err != nil {
					return fmt.Errorf("unable to read row value: %s", err.Error())
				}

				arr := ret[pkey.Name()].([]int64)
				ret[pkey.Name()] = append(arr, val)
			case spansql.Float64:
				var val float64
				err := r.ColumnByName(pkey.Name(), &val)
				if err != nil {
					return fmt.Errorf("unable to read row value: %s", err.Error())
				}

				arr := ret[pkey.Name()].([]float64)
				ret[pkey.Name()] = append(arr, val)
			case spansql.Bytes:
				var val []byte
				err := r.ColumnByName(pkey.Name(), &val)
				if err != nil {
					return fmt.Errorf("unable to read row value: %s", err.Error())
				}

				arr := ret[pkey.Name()].([][]byte)
				ret[pkey.Name()] = append(arr, val)
			case spansql.Timestamp:
				var val time.Time
				err := r.ColumnByName(pkey.Name(), &val)
				if err != nil {
					return fmt.Errorf("unable to read row value: %s", err.Error())
				}

				arr := ret[pkey.Name()].([]time.Time)
				ret[pkey.Name()] = append(arr, val)
			case spansql.Date:
				var val civil.Date
				err := r.ColumnByName(pkey.Name(), &val)
				if err != nil {
					return fmt.Errorf("unable to read row value: %s", err.Error())
				}

				arr := ret[pkey.Name()].([]civil.Date)
				ret[pkey.Name()] = append(arr, val)
			case spansql.Numeric:
				var val *big.Rat
				err := r.ColumnByName(pkey.Name(), &val)
				if err != nil {
					return fmt.Errorf("unable to read row value: %s", err.Error())
				}

				arr := ret[pkey.Name()].([]*big.Rat)
				ret[pkey.Name()] = append(arr, val)
			case spansql.JSON:
				val := make(map[string]interface{}) // TODO: This needs to be spanner.NullJSON
				err := r.ColumnByName(pkey.Name(), &val)
				if err != nil {
					return fmt.Errorf("unable to read row value: %s", err.Error())
				}

				arr := ret[pkey.Name()].([]map[string]interface{}) // TODO: This needs to be spanner.NullJSON
				ret[pkey.Name()] = append(arr, val)
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error fetching results from table sample: %s", err.Error())
	}

	return ret, nil
}

func GetReadGeneratorMap(samples map[string]interface{}, cols []string) (*sample.SampleGenerator, error) {
	ret, err := sample.NewSampleGenerator(
		rand.New(rand.NewSource(time.Now().UnixNano())),
		samples,
		cols,
	)

	return ret, err
}
