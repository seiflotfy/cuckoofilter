package cuckoofilter

import (
	"encoding/binary"
	"hash/fnv"
)

// getHash returns a 32-bit hash value for the given data.
func getHash(data []byte) []byte {
	hasher := fnv.New64()
	hasher.Write(data)
	hash64 := hasher.Sum64()
	hash := make([]byte, 8)
	binary.BigEndian.PutUint64(hash, hash64)
	return hash[4:]
}

func getAltIndex(fp []byte, i uint) uint {
	hash := getHash(fp)
	return i ^ uint(binary.BigEndian.Uint32(hash))
}

// getIndicesAndFingerprint returns the 2 bucket indices and fingerprint to be used
func getIndicesAndFingerprint(data []byte) (uint, uint, []byte) {
	hash := getHash(data)
	f := hash[0:fingerprintSize]
	i1 := uint(binary.BigEndian.Uint32(hash))
	i2 := getAltIndex(f, i1)
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
