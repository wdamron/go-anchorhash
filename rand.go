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

const (
	fleaSeed       = uint32(0xf1ea5eed)
	fleaRot1       = 27
	fleaRot2       = 17
	fleaInitRounds = 5 // initializing with 5 rounds works well enough in practice
)

// "A small noncryptographic PRNG" (Jenkins, 2007)
// http://burtleburtle.net/bob/rand/smallprng.html
//
// Also known as FLEA
func fleaInit(key uint64) (a, b, c, d uint32) {
	seed := uint32((key >> 32) ^ key)
	return fleaSeed, seed, seed, seed
}

func fleaRound(a, b, c, d uint32) (uint32, uint32, uint32, uint32) {
	e := a - fleaRot(b, fleaRot1)
	a = b ^ fleaRot(c, fleaRot2)
	b = c + d
	c = d + e
	d = e + a
	return a, b, c, d
}

func fleaRot(n, shift uint32) uint32 { return (n << shift) | (n >> (32 - shift)) }
