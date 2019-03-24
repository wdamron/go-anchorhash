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
	"testing"
)

var _benchIgnore uint32 = 0

func BenchmarkGetBucket_9_10(b *testing.B) {
	const (
		buckets = 10
		used    = 9
	)

	a := NewCompactAnchor(buckets, used)
	n := b.N
	b.ResetTimer()
	for i := 0; i < n; i++ {
		_benchIgnore += uint32(a.GetBucket(uint64(i)))
	}
}

func BenchmarkGetBucket_5_10(b *testing.B) {
	const (
		buckets = 10
		used    = 5
	)

	a := NewCompactAnchor(buckets, used)
	n := b.N
	b.ResetTimer()
	for i := 0; i < n; i++ {
		_benchIgnore += uint32(a.GetBucket(uint64(i)))
	}
}

func BenchmarkGetBucket_900k_1m(b *testing.B) {
	const (
		buckets = 1000000
		used    = 900000
	)

	a := NewAnchor(buckets, used)
	n := b.N
	b.ResetTimer()
	for i := 0; i < n; i++ {
		_benchIgnore += a.GetBucket(uint64(i))
	}
}

func BenchmarkGetBucket_500k_1m(b *testing.B) {
	const (
		buckets = 1000000
		used    = 500000
	)

	a := NewAnchor(buckets, used)
	n := b.N
	b.ResetTimer()
	for i := 0; i < n; i++ {
		_benchIgnore += a.GetBucket(uint64(i))
	}
}
