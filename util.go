package cuckoofilter

import (
	"encoding/binary"

	"code.google.com/p/gofarmhash"
)

// getHash returns a 32-bit hash value for the given data.
func getHash(data []byte) []byte {
	hash32 := farmhash.FingerPrint32(data)
	hash := make([]byte, 4)
	binary.BigEndian.PutUint32(hash, hash32)
	return hash
}

func getAltIndex(fp []byte, i uint) uint {
	return i ^ uint(binary.BigEndian.Uint32(getHash(fp)))
}

// getIndicesAndFingerprint returns the 2 bucket indices and fingerprint to be used
func getIndicesAndFingerprint(data []byte) (uint, uint, []byte) {
	hash := getHash(data)
	f := hash[0:fingerprintSize]
	i1 := uint(binary.BigEndian.Uint32(hash))
	i2 := getAltIndex(f, i1)
	return i1, i2, f
}

func getNextPow2(x uint) uint {
	x--
	x |= x >> 1
	x |= x >> 2
	x |= x >> 4
	x |= x >> 8
	x |= x >> 16
	x |= x >> 32
	x++
	return x
}
