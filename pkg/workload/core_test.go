package workload

import (
	"testing"

	"github.com/cloudspannerecosystem/gcsb/pkg/config"
	"github.com/cloudspannerecosystem/gcsb/pkg/schema"
)

func TestBucketOps(t *testing.T) {
	tests := []struct {
		desc       string
		operations int
		buckets    int
		want       []int
	}{
		{
			desc:       "operations are evenly distributed",
			operations: 10,
			buckets:    5,
			want:       []int{2, 2, 2, 2, 2},
		},
		{
			desc:       "operations are distributed unequally",
			operations: 10,
			buckets:    4,
			want:       []int{3, 3, 2, 2},
		},
		{
			desc:       "some buckets have empty operations",
			operations: 5,
			buckets:    10,
			want:       []int{1, 1, 1, 1, 1, 0, 0, 0, 0, 0},
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			got := (&CoreWorkload{}).bucketOps(test.operations, test.buckets)
			if !isSameSlice(got, test.want) {
				t.Errorf("bucketOps(%v, %v) = %v, but want = %v", test.operations, test.buckets, got, test.want)
			}
		})
	}
}

func TestPlan(t *testing.T) {
	testSchema := schema.NewSchema()
	t1 := schema.NewTable()
	t1.SetName("Singers")
	testSchema.Tables().AddTable(t1)

	t2 := schema.NewTable()
	t2.SetName("Albums")
	t2.SetParent(t1)
	t2.SetParentName(t1.Name())
	t1.SetChild(t2)
	t1.SetChildName(t2.Name())
	testSchema.Tables().AddTable(t2)

	tests := []struct {
		desc           string
		initialTargets []string
		wantTargets    []string
	}{
		{
			desc:           "no tables are planned",
			initialTargets: []string{},
			wantTargets:    []string{},
		},
		{
			desc:           "only parent table is planned",
			initialTargets: []string{"Singers"},
			wantTargets:    []string{"Singers"},
		},
		{
			desc:           "parent table is also planned",
			initialTargets: []string{"Albums"},
			wantTargets:    []string{"Singers", "Albums"},
		},
		{
			desc:           "all tables are planned",
			initialTargets: []string{"Singers", "Albums"},
			wantTargets:    []string{"Singers", "Albums"},
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			workload := CoreWorkload{
				Schema: testSchema,
				Config: &config.Config{},
			}

			err := workload.Plan(JobLoad, test.initialTargets)
			if err != nil {
				t.Fatalf("workload.Plan got error: %v", err)
			}

			got := extractTableNamesFromTargets(workload.plan)
			if !isSameStringSet(got, test.wantTargets) {
				t.Errorf("workload.Plan(%v) = %v, but want = %v", test.initialTargets, got, test.wantTargets)
			}
		})
	}
}

func isSameSlice(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func isSameStringSet(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	counters := make(map[string]int)
	for _, v := range a {
		counters[v]++
	}
	for _, v := range b {
		counters[v]--
	}

	for _, counter := range counters {
		if counter != 0 {
			return false
		}
	}
	return true
}

func extractTableNamesFromTargets(targets []*Target) []string {
	var names []string
	for _, t := range targets {
		names = append(names, t.TableName)
	}
	return names
}
