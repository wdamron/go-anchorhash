// The MIT License (MIT)
//
// Copyright (c) 2019 West Damron
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package anchor

import (
	"math"
	"reflect"
	"testing"
)

func TestAnchor(t *testing.T) {
	const (
		buckets = 7
		used    = 7
	)

	a := NewAnchor(buckets, used)

	// Fig. 2(d)
	a.RemoveBucket(6)
	a.RemoveBucket(5)

	// Ex. 13
	a.RemoveBucket(1)

	if !reflect.DeepEqual(a.W, []uint32{0, 4, 2, 3, 4, 5, 6}) {
		t.Fatalf("W = %#+v", a.W)
	}
	if !reflect.DeepEqual(a.L, []uint32{0, 1, 2, 3, 1, 5, 6}) {
		t.Fatalf("L = %#+v", a.L)
	}

	// Ex. 14/15
	a.RemoveBucket(0)

	if !reflect.DeepEqual(a.K, []uint32{3, 4, 2, 3, 4, 5, 6}) {
		t.Fatalf("K = %#+v", a.K)
	}

	// Restore
	a.AddBucket() // 0
	a.AddBucket() // 1
	a.AddBucket() // 5
	a.AddBucket() // 6

	if !reflect.DeepEqual(a.W, []uint32{0, 1, 2, 3, 4, 5, 6}) {
		t.Fatalf("W = %#+v", a.W)
	}
	if !reflect.DeepEqual(a.L, []uint32{0, 1, 2, 3, 4, 5, 6}) {
		t.Fatalf("L = %#+v", a.L)
	}
	if !reflect.DeepEqual(a.K, []uint32{0, 1, 2, 3, 4, 5, 6}) {
		t.Fatalf("K = %#+v", a.K)
	}
}

func TestPaths(t *testing.T) {
	const (
		buckets = 10
		used    = 5
	)
	a := NewAnchor(buckets, used)
	path := make([]uint32, 0, 64)

	const count = 1e6
	sum := 0
	for i := uint64(0); i < count; i++ {
		path = a.GetPath(i, path)
		sum += len(path)
		path = path[:0]
	}
	t.Logf("avg trace   = %v\n", float64(sum)/count)
	t.Logf("1 + ln(a/w) = %v\n", float64(1)+math.Log(float64(buckets)/float64(used)))
}

func TestOrdering(t *testing.T) {
	const (
		buckets = 5
		used    = 5
	)
	a := NewAnchor(buckets, used)

	counts := make([]int, buckets)
	for i := uint64(0); i < 1e6; i++ {
		counts[a.GetBucket(i)]++
	}
	t.Logf("%#+v\n", counts)

	a.RemoveBucket(4)
	a.RemoveBucket(3)
	a.RemoveBucket(2)
	t.Logf("removed b=4,3,2\n")

	counts = make([]int, buckets)
	for i := uint64(0); i < 1e6; i++ {
		counts[a.GetBucket(i)]++
	}
	t.Logf("%#+v\n", counts)

	a.AddBucket()
	a.AddBucket()
	a.AddBucket()
	t.Logf("added b=2,3,4\n")

	counts = make([]int, buckets)
	for i := uint64(0); i < 1e6; i++ {
		counts[a.GetBucket(i)]++
	}
	t.Logf("%#+v\n", counts)

	a.RemoveBucket(2)
	a.RemoveBucket(3)
	a.RemoveBucket(4)
	t.Logf("removed b=2,3,4\n")

	counts = make([]int, buckets)
	for i := uint64(0); i < 1e6; i++ {
		counts[a.GetBucket(i)]++
	}
	t.Logf("%#+v\n", counts)

	a.AddBucket()
	a.AddBucket()
	a.AddBucket()
	t.Logf("added b=4,3,2\n")

	counts = make([]int, buckets)
	for i := uint64(0); i < 1e6; i++ {
		counts[a.GetBucket(i)]++
	}
	t.Logf("%#+v\n", counts)
}

func TestDistributionSimple(t *testing.T) {
	const (
		buckets = 10
		used    = 5
	)
	a := NewAnchor(buckets, used)

	counts := make([]int, buckets)
	for i := uint64(0); i < 1e6; i++ {
		counts[a.GetBucket(i)]++
	}
	t.Logf("%#+v\n", counts)

	a.AddBucket()
	a.AddBucket()
	a.AddBucket()
	a.AddBucket()
	a.AddBucket()

	counts = make([]int, buckets)
	for i := uint64(0); i < 1e6; i++ {
		counts[a.GetBucket(i)]++
	}
	t.Logf("%#+v\n", counts)

	a.RemoveBucket(9)
	a.RemoveBucket(8)
	a.RemoveBucket(7)
	a.RemoveBucket(6)
	a.RemoveBucket(5)

	counts = make([]int, buckets)
	for i := uint64(0); i < 1e6; i++ {
		counts[a.GetBucket(i)]++
	}
	t.Logf("%#+v\n", counts)
}

func TestDistributionExtended(t *testing.T) {
	const (
		buckets = 10
		used    = 10
	)
	a := NewAnchor(buckets, used)

	counts := make([]int, buckets)
	for i := uint64(0); i < 1e6; i++ {
		counts[a.GetBucket(i)]++
	}
	t.Logf("%#+v\n", counts)

	a.RemoveBucket(9)
	t.Logf("removed b=9\n")
	counts = make([]int, buckets)
	for i := uint64(0); i < 1e6; i++ {
		counts[a.GetBucket(i)]++
	}
	t.Logf("%#+v\n", counts)

	a.RemoveBucket(5)
	t.Logf("removed b=5\n")
	counts = make([]int, buckets)
	for i := uint64(0); i < 1e6; i++ {
		counts[a.GetBucket(i)]++
	}
	t.Logf("%#+v\n", counts)

	a.RemoveBucket(3)
	t.Logf("removed b=3\n")
	counts = make([]int, buckets)
	for i := uint64(0); i < 1e6; i++ {
		counts[a.GetBucket(i)]++
	}
	t.Logf("%#+v\n", counts)

	a.AddBucket()
	t.Logf("re-added b=3\n")
	counts = make([]int, buckets)
	for i := uint64(0); i < 1e6; i++ {
		counts[a.GetBucket(i)]++
	}
	t.Logf("%#+v\n", counts)

	a.AddBucket()
	t.Logf("re-added b=5\n")
	counts = make([]int, buckets)
	for i := uint64(0); i < 1e6; i++ {
		counts[a.GetBucket(i)]++
	}
	t.Logf("%#+v\n", counts)

	a.AddBucket()
	t.Logf("re-added b=9\n")
	counts = make([]int, buckets)
	for i := uint64(0); i < 1e6; i++ {
		counts[a.GetBucket(i)]++
	}
	t.Logf("%#+v\n", counts)
}
