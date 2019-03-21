// The MIT License (MIT)
//
// Copyright (c) 2019 West Damron
// Portions Copyright (c) 2007 Bob Jenkins
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

// "A small noncryptographic PRNG" (Jenkins, 2007)
// http://burtleburtle.net/bob/rand/smallprng.html
type randctx struct {
	a, b, c, d uint32
}

func newRand(seed uint32) randctx {
	r := randctx{a: 0xf1ea5eed, b: seed, c: seed, d: seed}
	for i := uint32(0); i < 3; i++ {
		_ = r.next()
	}
	return r
}

func (r *randctx) next() uint32 {
	e := r.a - (((r.b) << (27)) | ((r.b) >> (32 - (27))))
	r.a = r.b ^ (((r.c) << (17)) | ((r.c) >> (32 - (17))))
	r.b = r.c + r.d
	r.c = r.d + e
	r.d = e + r.a
	return r.d
}
