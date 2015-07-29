package cuckoofilter

const fingerprintSize = 1
const bucketSize = 4

type fingerprint [fingerprintSize]byte
type bucket [bucketSize]fingerprint

var nullFp = fingerprint{0}

func (b *bucket) insert(fp fingerprint) bool {
	for i, tfp := range b {
		if tfp == nullFp {
			b[i] = fp
			return true
		}
	}
	return false
}

func (b *bucket) delete(fp fingerprint) bool {
	for i, tfp := range b {
		if tfp == fp {
			b[i] = nullFp
			return true
		}
	}
	return false
}

func (b *bucket) getFingerprintIndex(fp fingerprint) int {
	for i, tfp := range b {
		if tfp == fp {
			return i
		}
	}
	return -1
}
