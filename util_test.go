package cuckoo

import (
	"crypto/rand"
	"fmt"
	"io"
	"math/bits"
	"testing"
)

func TestIndexAndFP(t *testing.T) {
	data := []byte("seif")
	bucketPow := uint(bits.TrailingZeros(1024))
	i1, i2, fp := getIndicesAndFingerprint(data, bucketPow)
	i11 := getAltIndex(fp, i2, bucketPow)
	i22 := getAltIndex(fp, i1, bucketPow)
	if i1 != i11 {
		t.Errorf("Expected i1 == i11, instead %d != %d", i1, i11)
	}
	if i2 != i22 {
		t.Errorf("Expected i2 == i22, instead %d != %d", i2, i22)
	}
}

func TestCap(t *testing.T) {
	const capacity = 10000
	fmt.Println(getNextPow2(uint64(capacity)) / bucketSize)
}

func TestInsert(t *testing.T) {
	const cap = 10000
	filter := NewFilter(cap)

	var hash [32]byte
	io.ReadFull(rand.Reader, hash[:])

	for i := 0; i < 100; i++ {
		filter.Insert(hash[:])
	}

	fmt.Println(filter.Count())
}

func TestFilter_Lookup(t *testing.T) {
	const cap = 10000

	filter := NewFilter(cap)
	var hash [32]byte
	var m = make(map[[32]byte]struct{})

	for i := 0; i < cap; i++ {
		io.ReadFull(rand.Reader, hash[:])
		m[hash] = struct{}{}
		filter.Insert(hash[:])
	}

	fmt.Println(len(m))

	var lookFail int
	for k := range m {
		if !filter.Lookup(k[:]) {
			lookFail++
		}
	}

	fmt.Println(lookFail)
}

func TestReset(t *testing.T) {
	const cap = 10000

	filter := NewFilter(cap)
	var hash [32]byte

	var insertSuccess int
	var fail int

	for i := 0; i < 10*cap; i++ {
		io.ReadFull(rand.Reader, hash[:])

		if filter.Insert(hash[:]) {
			insertSuccess++
		} else {
			fail++
			filter.Reset()
		}
	}

	fmt.Println("insert success", insertSuccess)

	fmt.Println("insert fail", fail)
}

func BenchmarkFilter_Reset(b *testing.B) {
	const cap = 10000
	filter := NewFilter(cap)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		filter.Reset()
	}
}

func BenchmarkFilter_Insert(b *testing.B) {
	const cap = 10000
	filter := NewFilter(cap)

	b.ResetTimer()

	var hash [32]byte
	for i := 0; i < b.N; i++ {
		io.ReadFull(rand.Reader, hash[:])
		filter.Insert(hash[:])
	}
}

func BenchmarkFilter_Lookup(b *testing.B) {
	const cap = 10000
	filter := NewFilter(cap)

	var hash [32]byte
	for i := 0; i < 10000; i++ {
		io.ReadFull(rand.Reader, hash[:])
		filter.Insert(hash[:])
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		io.ReadFull(rand.Reader, hash[:])
		filter.Lookup(hash[:])
	}
}

func TestBucket_Reset(t *testing.T) {
	var bkt bucket
	for i := byte(0); i < bucketSize; i++ {
		bkt[i] = i
	}
	fmt.Println(bkt)
	bkt.reset()
	fmt.Println(bkt)
}
