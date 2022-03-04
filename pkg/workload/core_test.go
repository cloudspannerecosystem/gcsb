package workload

import "testing"

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
