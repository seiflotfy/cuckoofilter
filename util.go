package cuckoofilter

import (
	"encoding/binary"
	"sync"

	"github.com/dgryski/go-metro"
)

var hlock sync.Mutex

func getAltIndex(fp byte, i uint, numBuckets uint) uint {
	bytes := []byte{0, 0, 0, 0, 0, 0, 0, fp}
	hash := binary.LittleEndian.Uint64(bytes)
	return uint(uint64(i)^(hash*0x5bd1e995)) % numBuckets
}

func getFingerprint(data []byte) byte {
	hlock.Lock()
	defer hlock.Unlock()

	fp := byte(metro.Hash64(data, 1337))
	if fp == 0 {
		fp += 7
	}
	return fp
}

// getIndicesAndFingerprint returns the 2 bucket indices and fingerprint to be used
func getIndicesAndFingerprint(data []byte, numBuckets uint) (uint, uint, byte) {
	hash := metro.Hash64(data, 1337)
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
