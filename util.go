package cuckoofilter

import (
	"encoding/binary"

	"code.google.com/p/gofarmhash"
)

// getHash returns a 64-bit hash value for the given data.
func getHash(data []byte) uint64 {
	return farmhash.Hash64(data)
}

func getAltIndex(fp []byte, i uint, numBuckets uint) uint {
	hash := getHash(fp)
	return (i ^ uint(hash)) % numBuckets
}

func getFingerprint(hash64 uint64) []byte {
	hash := make([]byte, 8)
	binary.BigEndian.PutUint64(hash, hash64)
	return hash[:fingerprintSize]
}

// getIndicesAndFingerprint returns the 2 bucket indices and fingerprint to be used
func getIndicesAndFingerprint(data []byte, numBuckets uint) (uint, uint, []byte) {
	hash := getHash(data)
	f := getFingerprint(hash)
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
