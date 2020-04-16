package cuckoo

import (
	metro "github.com/dgryski/go-metro"
)

var (
	altHash = [256]uint{}
	masks   = [65]uint{}
)

func init() {
	for i := 0; i < 256; i++ {
		altHash[i] = (uint(metro.Hash64([]byte{byte(i)}, 1337)))
	}
	for i := uint(0); i <= 64; i++ {
		masks[i] = (1 << i) - 1
	}
}

func getAltIndex(fp byte, i uint, bucketPow uint) uint {
	mask := masks[bucketPow]
	hash := altHash[fp] & mask
	return (i & mask) ^ hash
}

func getFingerprint(data []byte) byte {
	fp := byte(metro.Hash64(data, 1335)%255 + 1)
	return fp
}

// getIndicesAndFingerprint returns the 2 bucket indices and fingerprint to be used
func getIndicesAndFingerprint(data []byte, bucketPow uint) (uint, uint, byte) {
	hash := metro.Hash64(data, 1337)
	f := getFingerprint(data)
	i1 := uint(hash) & masks[bucketPow]
	i2 := getAltIndex(f, i1, bucketPow)
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
