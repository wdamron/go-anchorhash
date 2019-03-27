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

// Compact, minimal-memory AnchorHash implementation.
//
// Buckets will be stored as unsigned 16-bit integers, so the maximum size of the working set
// will be limited to 65,536 buckets. This implementation offers improved cache-locality
// relative to Anchor.
type CompactAnchor struct {
	// We use an integer array A of size a to represent the Anchor.
	//
	// Each bucket b ∈ {0, 1, ..., a−1} is represented by A[b] that either equals 0 if b
	// is a working bucket (i.e., A[b] = 0 if b ∈ W), or else equals the size of the working
	// set just after its removal (i.e., A[b] = |Wb| if b ∈ R).
	A []uint16
	// K stores the successor for each removed bucket b (i.e. the bucket that replaced it in W).
	K []uint16
	// W always contains the current set of working buckets in their desired order.
	W []uint16
	// L stores the most recent location for each bucket within W.
	L []uint16
	// R saves removed buckets in a LIFO order for possible future bucket additions.
	R []uint16
	// N is the current length of W
	N uint16
}

// Create a new anchor with a given capacity and initial size.
//
// 	INITANCHOR(a, w)
// 	A[b] ← 0 for b = 0, 1, ..., a−1    ◃ |Wb| ← 0 for b ∈ A
// 	R ← ∅                              ◃ Empty stack
// 	N ← w                              ◃ Number of initially working buckets
// 	K[b] ← L[b] ← W[b] ← b for b = 0, 1, ..., a−1
// 	for b = a−1 downto w do            ◃ Remove initially unused buckets
// 	  REMOVEBUCKET(b)
func NewCompactAnchor(buckets, used uint16) *CompactAnchor {
	a := &CompactAnchor{
		A: make([]uint16, buckets),
		K: make([]uint16, buckets),
		W: make([]uint16, buckets),
		L: make([]uint16, buckets),
		R: make([]uint16, buckets-used, buckets),
		N: uint16(used),
	}
	for b := uint16(0); b < uint16(used); b++ {
		a.K[b], a.W[b], a.L[b] = b, b, b
	}
	for b, r := uint16(buckets)-1, 0; b >= uint16(used); b, r = b-1, r+1 {
		a.A[b], a.R[r] = b, b
	}
	return a
}

// Get the bucket which a hash-key is assigned to.
//
// If the path for a given key contains any non-working buckets, the path (and in turn,
// the assigned bucket for the key) will be determined by the order in which the non-working
// buckets were removed. To maintain consistency in a distributed system, all agents must
// reach consensus on the ordering of changes to the working set. For more information,
// see Section III, Theorem 1 in the paper.
//
// 	GETBUCKET(k)
// 	b ← hash(k) mod a
// 	while A[b] > 0 do          ◃ b is removed
// 	  h ← hb(k)                ◃ hb(k) ≡ hash(k) mod A[b]
// 	  while A[h] ≥ A[b] do     ◃ Wb[h] != h, b removed prior to h
// 	    h ← K[h]               ◃ search for Wb[h]
// 	  b ← h
// 	return b
func (a *CompactAnchor) GetBucket(key uint64) uint16 {
	A, K := a.A, a.K
	ha, hb, hc, hd := fleaInit(key)
	b := uint16(fastMod(uint64(hd), uint64(len(A))))
	for A[b] > 0 {
		ha, hb, hc, hd = fleaRound(ha, hb, hc, hd)
		h := uint16(fastMod(uint64(hd), uint64(A[b])))
		for A[h] >= A[b] {
			h = K[h]
		}
		b = h
	}
	return b
}

// Get the path to the bucket which a hash-key is assigned to.
//
// The returned path will contain all buckets traversed while searching for a
// working bucket. The final bucket in the path will be the assigned bucket for
// the given key.
//
// Buckets will be appended to the provided buffer, though a different slice will
// be returned if the length of the path exceeds the capacity of the buffer.
//
// If the path for a given key contains any non-working buckets, the path (and in turn,
// the assigned bucket for the key) will be determined by the order in which the non-working
// buckets were removed. To maintain consistency in a distributed system, all agents must
// reach consensus on the ordering of changes to the working set. For more information,
// see Section III, Theorem 1 in the paper.
//
// 	GETPATH(k, P)
// 	b ← hash(k) mod a
// 	P.push(b)
// 	while A[b] > 0 do          ◃ b is removed
// 	  h ← hb(k)                ◃ hb(k) ≡ hash(k) mod A[b]
// 	  P.push(h)
// 	  while A[h] ≥ A[b] do     ◃ Wb[h] != h, b removed prior to h
// 	    h ← K[h]               ◃ search for Wb[h]
// 	    P.push(h)
// 	  b ← h
// 	return P
func (a *CompactAnchor) GetPath(key uint64, pathBuffer []uint16) []uint16 {
	A, K := a.A, a.K
	ha, hb, hc, hd := fleaInit(key)
	b := uint16(fastMod(uint64(hd), uint64(len(A))))
	pathBuffer = append(pathBuffer, b)
	for A[b] > 0 {
		ha, hb, hc, hd = fleaRound(ha, hb, hc, hd)
		h := uint16(fastMod(uint64(hd), uint64(A[b])))
		pathBuffer = append(pathBuffer, h)
		for A[h] >= A[b] {
			h = K[h]
			pathBuffer = append(pathBuffer, h)
		}
		b = h
	}
	return pathBuffer
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
func (a *CompactAnchor) AddBucket() uint16 {
	A, K, W, L, R, N := a.A, a.K, a.W, a.L, a.R, a.N
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
func (a *CompactAnchor) RemoveBucket(b uint16) {
	if a.A[b] != 0 {
		return
	}
	a.N--
	A, K, W, L, N := a.A, a.K, a.W, a.L, a.N
	a.R = append(a.R, b)
	A[b] = N
	W[L[b]], K[b] = W[N], W[N]
	L[W[N]] = L[b]
}
