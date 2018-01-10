package cuckoofilter

import (
	"fmt"
	"math/rand"
)

const maxCuckooCount = 500

//CuckooFilter represents a probabalistic counter
type CuckooFilter struct {
	buckets []bucket
	count   uint
}

//NewCuckooFilter returns a new cuckoofilter with a given capacity
func NewCuckooFilter(capacity uint) *CuckooFilter {
	capacity = getNextPow2(uint64(capacity)) / bucketSize
	if capacity == 0 {
		capacity = 1
	}
	buckets := make([]bucket, capacity)
	for i := range buckets {
		buckets[i] = [bucketSize]byte{}
	}
	return &CuckooFilter{buckets, 0}
}

//NewDefaultCuckooFilter returns a new cuckoofilter with the default capacity of 1000000
func NewDefaultCuckooFilter() *CuckooFilter {
	return NewCuckooFilter(1000000)
}

//Lookup returns true if data is in the counter
func (cf *CuckooFilter) Lookup(data []byte) bool {
	i1, i2, fp := getIndicesAndFingerprint(data, uint(len(cf.buckets)))
	b1, b2 := cf.buckets[i1], cf.buckets[i2]
	return b1.getFingerprintIndex(fp) > -1 || b2.getFingerprintIndex(fp) > -1
}

func randi(i1, i2 uint) uint {
	if rand.Intn(2) == 0 {
		return i1
	}
	return i2
}

//Insert inserts data into the counter and returns true upon success
func (cf *CuckooFilter) Insert(data []byte) bool {
	i1, i2, fp := getIndicesAndFingerprint(data, uint(len(cf.buckets)))
	if cf.insert(fp, i1) || cf.insert(fp, i2) {
		return true
	}
	return cf.reinsert(fp, randi(i1, i2))
}

//InsertUnique inserts data into the counter if not exists and returns true upon success
func (cf *CuckooFilter) InsertUnique(data []byte) bool {
	if cf.Lookup(data) {
		return false
	}
	return cf.Insert(data)
}

func (cf *CuckooFilter) insert(fp byte, i uint) bool {
	if cf.buckets[i].insert(fp) {
		cf.count++
		return true
	}
	return false
}

func (cf *CuckooFilter) reinsert(fp byte, i uint) bool {
	for k := 0; k < maxCuckooCount; k++ {
		j := rand.Intn(bucketSize)
		oldfp := fp
		fp = cf.buckets[i][j]
		cf.buckets[i][j] = oldfp

		// look in the alternate location for that random element
		i = getAltIndex(fp, i, uint(len(cf.buckets)))
		if cf.insert(fp, i) {
			return true
		}
	}
	return false
}

//Delete data from counter if exists and return if deleted or not
func (cf *CuckooFilter) Delete(data []byte) bool {
	i1, i2, fp := getIndicesAndFingerprint(data, uint(len(cf.buckets)))
	return cf.delete(fp, i1) || cf.delete(fp, i2)
}

func (cf *CuckooFilter) delete(fp byte, i uint) bool {
	if cf.buckets[i].delete(fp) {
		cf.count--
		return true
	}
	return false
}

//Count returns the number of items in the counter
func (cf *CuckooFilter) Count() uint {
	return cf.count
}

// Encode returns a byte slice representing a Cuckoofilter
func (cf *CuckooFilter) Encode() []byte {
	bytes := make([]byte, len(cf.buckets)*bucketSize)
	for i, b := range cf.buckets {
		for j, f := range b {
			index := (i * len(b)) + j
			bytes[index] = f
		}
	}
	return bytes
}

// Decode returns a Cuckoofilter from a byte slice
func Decode(bytes []byte) (*CuckooFilter, error) {
	var count uint
	if len(bytes)%4 != 0 {
		return nil, fmt.Errorf("expected bytes to be multiuple of 4, got %d", len(bytes))
	}
	buckets := make([]bucket, len(bytes)/4)
	for i, b := range buckets {
		for j := range b {
			index := (i * len(b)) + j
			if bytes[index] != 0 {
				buckets[i][j] = bytes[index]
				count++
			}
		}
	}
	return &CuckooFilter{
		buckets: buckets,
		count:   count,
	}, nil
}
