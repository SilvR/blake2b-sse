/*
 * Copyright 2016 Frank Wessels <fwessels@xs4all.nl>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package blake2b

import (
	"encoding/hex"
	"hash"
	"testing"
)

type blake2bTest struct {
	in     []byte
	out512 string
	// FIXME: add out256.
}

func gen128() []byte {
	in := make([]byte, 128)
	for i := range in {
		in[i] = byte(i)
	}
	return in
}

var golden = []blake2bTest{
	// Not working yet, enable when ready.
	/*
		{
			[]byte(""),
			"786a02f742015903c6c6fd852552d272912f4740e15847618a86e217f71f5419d25e1031afee585313896444934eb04b903a685b1448b755d56f701afe9be2ce",
		},
		{
			[]byte("a"),
			"333fcb4ee1aa7c115355ec66ceac917c8bfd815bf7587d325aec1864edd24e34d5abe2c6b1b5ee3face62fed78dbef802f2a85cb91d455a8f5249d330853cb3c",
		},
	*/
	{
		gen128(),
		"15ab5ef920578d8e181c75dce93f2f84234daa02075e212615a51f46d0f17d7c546fee070adb49845b3a397bfdafb51df2cb7b239a882a24e0e4b8aab25c6676",
	},
}

func testHash(t *testing.T, name string, in []byte, outHex string, oneShotResult []byte, digestFunc hash.Hash) {
	if calculated := hex.EncodeToString(oneShotResult); calculated != outHex {
		t.Errorf("one-shot result for %s(%q) = %q, but expected %q", name, in, calculated, outHex)
		return
	}

	for pass := 0; pass < 3; pass++ {
		if pass < 2 {
			digestFunc.Write(in)
		} else {
			digestFunc.Write(in[:len(in)/2])
			digestFunc.Sum(nil)
			digestFunc.Write(in[len(in)/2:])
		}

		if calculated := hex.EncodeToString(digestFunc.Sum(nil)); calculated != outHex {
			t.Errorf("%s(%q) = %q (in pass #%d), but expected %q", name, in, calculated, pass, outHex)
		}
		digestFunc.Reset()
	}
}

func TestGolden(t *testing.T) {
	for _, test := range golden {
		in := test.in

		blake512 := Sum512(in)
		testHash(t, "blake512", test.in, test.out512, blake512[:], New512())
	}
}

var bench = New512()
var buf = make([]byte, 8192)

func benchmarkSize(b *testing.B, size int) {
	b.SetBytes(int64(size))
	sum := make([]byte, bench.Size())
	for i := 0; i < b.N; i++ {
		bench.Reset()
		bench.Write(buf[:size])
		bench.Sum(sum[:0])
	}
}

func BenchmarkHash8Bytes(b *testing.B) {
	benchmarkSize(b, 8)
}

func BenchmarkHash1K(b *testing.B) {
	benchmarkSize(b, 1024)
}

func BenchmarkHash8K(b *testing.B) {
	benchmarkSize(b, 8192)
}
