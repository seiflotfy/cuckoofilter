package cuckoo

import (
	"crypto/rand"
	"io"
	"math/bits"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndexAndFP(t *testing.T) {
	data := []byte("seif")
	bucketPow := uint(bits.TrailingZeros(1024))
	i1, fp := getIndexAndFingerprint(data, bucketPow)
	i2 := getAltIndex(fp, i1, bucketPow)
	i11 := getAltIndex(fp, i2, bucketPow)
	i22 := getAltIndex(fp, i11, bucketPow)
	assert.EqualValues(t, i11, i1)
	assert.EqualValues(t, i22, i2)
}

func TestCap(t *testing.T) {
	const capacity = 10000
	res := getNextPow2(uint64(capacity)) / bucketSize
	assert.EqualValues(t, res, 4096)
}

func TestInsert(t *testing.T) {
	const cap = 10000
	filter := NewFilter(cap)

	var hash [32]byte
	io.ReadFull(rand.Reader, hash[:])

	for i := 0; i < 100; i++ {
		filter.Insert(hash[:])
	}

	assert.EqualValues(t, filter.Count(), 8)
}

func TestFilter_Lookup(t *testing.T) {
	const cap = 10000

	var (
		m      = make(map[[32]byte]struct{})
		filter = NewFilter(cap)
		hash   [32]byte
	)

	for i := 0; i < cap; i++ {
		io.ReadFull(rand.Reader, hash[:])
		m[hash] = struct{}{}
		filter.Insert(hash[:])
	}

	assert.EqualValues(t, len(m), 10000)

	var lookFail int
	for k := range m {
		if !filter.Lookup(k[:]) {
			lookFail++
		}
	}

	assert.EqualValues(t, lookFail, 0)
}

func TestReset(t *testing.T) {
	const cap = 10000

	var (
		filter        = NewFilter(cap)
		hash          [32]byte
		insertSuccess int
		insertFails   int
	)

	for i := 0; i < 10*cap; i++ {
		io.ReadFull(rand.Reader, hash[:])

		if filter.Insert(hash[:]) {
			insertSuccess++
		} else {
			insertFails++
			filter.Reset()
		}
	}

	assert.EqualValues(t, insertSuccess, 99994)
	assert.EqualValues(t, insertFails, 6)
}

func TestBucket_Reset(t *testing.T) {
	var bkt bucket
	for i := byte(0); i < bucketSize; i++ {
		bkt[i] = fingerprint(i)
	}
	bkt.reset()
	for _, val := range bkt {
		assert.EqualValues(t, 0, val)
	}
}
