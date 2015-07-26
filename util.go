package cuckoofilter

import (
	"hash"
	"hash/fnv"

	"code.google.com/p/gofarmhash"
)

var hashera hash.Hash64 = fnv.New64()

func getAltIndex(fp []byte, i uint, numBuckets uint) uint {
	hash := farmhash.Hash64(fp)
	if hash == 0 {
		hash += 1
	}
	return uint(uint64(i)^farmhash.Hash64(fp)) % numBuckets
}

func getFingerprint(data []byte) []byte {
	hashera.Reset()
	hashera.Write(data)
	hash := hashera.Sum(nil)
	return hash[:fingerprintSize]
}

// getIndicesAndFingerprint returns the 2 bucket indices and fingerprint to be used
func getIndicesAndFingerprint(data []byte, numBuckets uint) (uint, uint, []byte) {
	hash := farmhash.Hash64(data)
	f := getFingerprint(data)
	i1 := uint(hash) % numBuckets
	i2 := getAltIndex(f, i1, numBuckets)
	return i1, i2, f
}

func getNextPow2(n uint) uint {
	n--
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	n |= n >> 32
	n++
	return n
}
