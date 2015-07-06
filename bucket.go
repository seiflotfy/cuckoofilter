package cuckoofilter

import "bytes"

const fingerprintSize = 1
const bucketSize = 4

type fingerprint []byte
type bucket []fingerprint

func (b bucket) insert(fp fingerprint) bool {
	for i, tfp := range b {
		if tfp == nil {
			b[i] = fp
			return true
		}
	}
	return false
}

func (b bucket) delete(fp fingerprint) bool {
	for i, tfp := range b {
		if bytes.Equal(fp, tfp) {
			b[i] = nil
			return true
		}
	}
	return false
}

func (b bucket) getFingerprintIndex(fp []byte) int {
	for i, tfp := range b {
		if bytes.Equal(tfp, fp) {
			return i
		}
	}
	return -1
}
