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

// package anchor provides a minimal-memory AnchorHash consistent-hash implementation for Go.
//
// [AnchorHash: A Scalable Consistent Hash]: https://arxiv.org/abs/1812.09674
package anchor

// Minimal-memory AnchorHash implementation.
type Anchor struct {
	// We use an integer array A of size a to represent the Anchor.
	//
	// Each bucket b ∈ {0, 1, ..., a−1} is represented by A[b] that either equals 0 if b
	// is a working bucket (i.e., A[b] = 0 if b ∈ W), or else equals the size of the working
	// set just after its removal (i.e., A[b] = |Wb| if b ∈ R).
	A []uint
	// R saves removed buckets in a LIFO order for possible future bucket additions.
	R []uint
	// W always contains the current set of working buckets in their desired order.
	W []uint
	// L stores the most recent location for each bucket within W.
	L []uint
	// K stores the successor for each removed bucket b (i.e. the bucket that replaced it in W).
	K []uint
	// N is the current length of W
	N uint
}

// Create a new anchor with a given capacity and initial size.
//
// 	INITANCHOR(a, w)
// 	A[b] ← 0 for b = 0, 1, ..., a−1    ◃ |Wb| ← 0 for b ∈ A
// 	R ← ∅                              ◃ Empty stack
// 	N ← w                              ◃ Number of initially working buckets
// 	K[b] ← L[b] ← W[b] ← b for b = 0, 1, ..., a−1
// 	for b = a−1 downto w do            ◃ Remove initially unused buckets
// 	  R.push(b)
func NewAnchor(buckets, used int) *Anchor {
	a := &Anchor{
		A: make([]uint, buckets),
		R: make([]uint, 0, buckets),
		W: make([]uint, buckets),
		L: make([]uint, buckets),
		K: make([]uint, buckets),
		N: uint(buckets),
	}
	for b := uint(0); b < uint(buckets); b++ {
		a.A[b], a.W[b], a.L[b], a.K[b] = 0, b, b, b
	}
	for b := uint(buckets) - 1; b >= uint(used); b-- {
		a.RemoveBucket(b)
	}
	return a
}

// Get the bucket which a hash is assigned to.
//
// 	GETBUCKET(k)
// 	b ← hash(k) mod a
// 	while A[b] > 0 do          ◃ b is removed
// 	  h ← hb(k)                ◃ hb(k) ≡ hash(k) mod A[b]
// 	  while A[h] ≥ A[b] do     ◃ Wb[h] != h, b removed prior to h
// 	    h ← K[h]               ◃ search for Wb[h]
// 	  b ← h
// 	return b
func (a *Anchor) GetBucket(hash uint) uint {
	A, K := a.A, a.K
	b := hash % uint(len(A))
	for A[b] > 0 {
		h := hash % A[b]
		for A[h] >= A[b] {
			h = K[h]
		}
		b = h
	}
	return b
}

// Add a bucket to the anchor.
//
// 	ADDBUCKET()
// 	b ← R.pop()
// 	A[b] ← 0       ◃ W ← W ∪ {b}, delete Wb
// 	L[W[N]] ← N
// 	W[L[b]] ← K[b] ← b
// 	N ← N + 1
// 	return b
func (a *Anchor) AddBucket() uint {
	A, R, W, L, K, N := a.A, a.R, a.W, a.L, a.K, a.N
	b := R[len(R)-1]
	a.R = R[:len(R)-1]
	A[b] = 0
	L[W[N]] = N
	W[L[b]], K[b] = b, b
	a.N++
	return b
}

// Remove a bucket from the anchor.
//
// 	REMOVEBUCKET(b)
// 	R.push(b)
// 	N ← N − 1
// 	A[b] ← N       ◃ Wb ← W \ b, A[b] ← |Wb|
// 	W[L[b]] ← K[b] ← W[N]
// 	L[W[N]] ← L[b]
func (a *Anchor) RemoveBucket(b uint) {
	a.N--
	A, W, L, K, N := a.A, a.W, a.L, a.K, a.N
	a.R = append(a.R, b)
	A[b] = N
	W[L[b]], K[b] = W[N], W[N]
	L[W[N]] = L[b]
}
