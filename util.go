package cuckoofilter

import (
	"encoding/binary"
	"hash/fnv"

	"code.google.com/p/gofarmhash"
)

var hashera = fnv.New64()

func getAltIndex(fp fingerprint, i uint, numBuckets uint) uint {
	bytes := make([]byte, 64, 64)
	for i, b := range fp {
		bytes[i] = b
	}

	hash := binary.LittleEndian.Uint64(bytes)
	return uint(uint64(i)^(hash*0x5bd1e995)) % numBuckets
}

func getFingerprint(data []byte) fingerprint {
	hashera.Reset()
	hashera.Write(data)
	hash := hashera.Sum(nil)

	fp := fingerprint{}
	for i := 0; i < fingerprintSize; i++ {
		fp[i] = hash[i]
	}
	if fp == nullFp {
		fp[0] += 7
	}
	return fp
}

// getIndicesAndFingerprint returns the 2 bucket indices and fingerprint to be used
func getIndicesAndFingerprint(data []byte, numBuckets uint) (uint, uint, fingerprint) {
	hash := farmhash.Hash64(data)
	f := getFingerprint(data)
	i1 := uint(hash) % numBuckets
	i2 := getAltIndex(f, i1, numBuckets)
	return i1, i2, f
}

func getNextPow2(n uint64) uint {
	n--
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	n |= n >> 32
	n++
	return uint(n)
}
