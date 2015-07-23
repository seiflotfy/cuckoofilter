package cuckoofilter

import (
	"bytes"
	"fmt"
)

const fingerprintSize = 1
const bucketSize = 2

type fingerprint []byte
type bucket []fingerprint

func (b bucket) insert(fp fingerprint) int {
	for i, tfp := range b {
		if tfp == nil {
			b[i] = fp
			return i
		}
	}
	return -1
}

func (b bucket) delete(fp fingerprint) int {
	for i, tfp := range b {
		if bytes.Equal(fp, tfp) {
			b[i] = nil
			return i
		}
	}
	return -1
}

func (b bucket) getFingerprintIndex(fp []byte) int {
	for i, tfp := range b {
		if bytes.Equal(tfp, fp) {
			return i
		}
	}
	return -1
}

func printBucket(buckets []bucket) {
	for _, b := range buckets {
		fmt.Print(b, "   ")
	}
	fmt.Println("")
}

func printIndicators(indicators [][]uint) {
	for _, i := range indicators {
		fmt.Print(i, "   ")
	}
	fmt.Println("")
}
