package cuckoofilter

import "math/rand"

const maxCuckooCount = 500

/*
CuckooFilter ...
*/
type CuckooFilter struct {
	buckets []bucket
	count   uint
}

/*
NewCuckooFilter ...
*/
func NewCuckooFilter(capacity uint) *CuckooFilter {
	capacity = getNextPow2(capacity)
	buckets := make([]bucket, capacity/bucketSize)
	for i := range buckets {
		buckets[i] = make([]fingerprint, bucketSize)
	}
	return &CuckooFilter{buckets, 0}
}

/*
NewDefaultCuckooFilter ...
*/
func NewDefaultCuckooFilter() *CuckooFilter {
	return NewCuckooFilter(1000000)
}

/*
Lookup ...
*/
func (cf *CuckooFilter) Lookup(data []byte) bool {
	i1, i2, fp := getIndicesAndFingerprint(data)
	b1, b2 := cf.buckets[i1%uint(len(cf.buckets))],
		cf.buckets[i2%uint(len(cf.buckets))]
	return b1.getFingerprintIndex(fp) > -1 || b2.getFingerprintIndex(fp) > -1
}

/*
Insert ...
*/
func (cf *CuckooFilter) Insert(data []byte) bool {
	i1, i2, fp := getIndicesAndFingerprint(data)
	if cf.insert(fp, i1) || cf.insert(fp, i2) {
		return true
	}
	return cf.reinsert(fp, i2)
}

/*
InsertUnique ...
*/
func (cf *CuckooFilter) InsertUnique(data []byte) bool {
	if cf.Lookup(data) {
		return false
	}
	i1, i2, fp := getIndicesAndFingerprint(data)
	if cf.insert(fp, i1) || cf.insert(fp, i2) {
		return true
	}
	return cf.reinsert(fp, i2)
}

/*
insert ...
*/
func (cf *CuckooFilter) insert(fp fingerprint, i uint) bool {
	b := cf.buckets[i%uint(len(cf.buckets))]
	if b.insert(fp) {
		cf.count++
		return true
	}
	return false
}

/*
reinsert ...
*/
func (cf *CuckooFilter) reinsert(fp fingerprint, i uint) bool {
	for k := 0; k < maxCuckooCount; k++ {
		j := rand.Intn(bucketSize)
		fp, newfp := cf.buckets[i%uint(len(cf.buckets))][j], fp
		cf.buckets[i%uint(len(cf.buckets))][j] = newfp

		// look in the alternate location for that random element
		i = getAltIndex(fp, i)

		if cf.insert(fp, i) {
			return true
		}
	}
	return false
}

/*
Delete ...
*/
func (cf *CuckooFilter) Delete(data []byte) bool {
	i1, i2, fp := getIndicesAndFingerprint(data)
	return cf.delete(fp, i1) || cf.delete(fp, i2)
}

/*
delete ...
*/
func (cf *CuckooFilter) delete(fp fingerprint, i uint) bool {
	b := cf.buckets[i%uint(len(cf.buckets))]
	if b.delete(fp) {
		cf.count--
		return true
	}
	return false
}
