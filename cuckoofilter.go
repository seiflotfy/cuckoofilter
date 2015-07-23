package cuckoofilter

import (
	"fmt"
	"math/rand"
)

const maxCuckooCount = 500

/*
CuckooFilter represents a probabalistic counter
*/
type CuckooFilter struct {
	buckets []bucket
	count   uint64
	indices [][]bool
}

/*
NewCuckooFilter returns a new cuckoofilter with a given capacity
*/
func NewCuckooFilter(capacity uint64) *CuckooFilter {
	capacity = getNextPow2(capacity) / bucketSize
	if capacity == 0 {
		capacity = 2
	}
	buckets := make([]bucket, capacity)
	indices := make([][]bool, capacity)
	for i := range buckets {
		buckets[i] = make([]fingerprint, bucketSize)
		indices[i] = make([]bool, bucketSize)
	}
	return &CuckooFilter{buckets, 0, indices}
}

/*
NewDefaultCuckooFilter returns a new cuckoofilter with the default capacity of 1000000
*/
func NewDefaultCuckooFilter() *CuckooFilter {
	return NewCuckooFilter(1000000)
}

/*
Lookup returns true if data is in the counter
*/
func (cf *CuckooFilter) Lookup(data []byte) bool {
	i1, i2, fp := getIndicesAndFingerprint(data, uint64(len(cf.buckets)))
	b1, b2 := cf.buckets[i1], cf.buckets[i2]
	return b1.getFingerprintIndex(fp) > -1 || b2.getFingerprintIndex(fp) > -1
}

/*
Inserts inserts data into the counter and returns true upon success
*/
func (cf *CuckooFilter) Insert(data []byte) bool {

	limit := uint64(90*len(cf.buckets)*4) / 100

	if cf.count >= limit {
		cf.expand()
	}

	i1, i2, fp := getIndicesAndFingerprint(data, uint64(len(cf.buckets)))
	if cf.insert(fp, i1, true) || cf.insert(fp, i2, false) {
		return true
	}
	return cf.reinsert(fp, i2)
}

/*
InsertUnique inserts data into the counter if not exists and returns true upon success
*/
func (cf *CuckooFilter) InsertUnique(data []byte) bool {
	if cf.Lookup(data) {
		return false
	}
	return cf.Insert(data)
}

func (cf *CuckooFilter) insert(fp fingerprint, i uint64, prim bool) bool {
	if j := cf.buckets[i].insert(fp); j > -1 {
		cf.count++
		cf.indices[i][j] = prim
		return true
	}
	return false
}

func (cf *CuckooFilter) reinsert(fp fingerprint, i uint64) bool {
	prim := false
	for k := 0; k < maxCuckooCount; k++ {
		j := rand.Intn(bucketSize)
		oldfp := fp
		oldprim := prim
		fp = cf.buckets[i][j]
		prim = !cf.indices[i][j]
		cf.buckets[i][j] = oldfp
		cf.indices[i][j] = oldprim

		// look in the alternate location for that random element
		i = getAltIndex(fp, i, uint64(len(cf.buckets)))
		if cf.insert(fp, i, prim) {
			return true
		}
	}
	return false
}

/*
Delete data from counter if exists and return if deleted or not
*/
func (cf *CuckooFilter) Delete(data []byte) bool {
	i1, i2, fp := getIndicesAndFingerprint(data, uint64(len(cf.buckets)))
	return cf.delete(fp, i1) || cf.delete(fp, i2)
}

func (cf *CuckooFilter) delete(fp fingerprint, i uint64) bool {
	if cf.buckets[i].delete(fp) {
		cf.count--
		return true
	}
	return false
}

/*
GetCount returns the number of items in the counter
*/
func (cf *CuckooFilter) GetCount() uint64 {
	return cf.count
}

func (cf *CuckooFilter) expand() {
	//N := len(cf.buckets)
	//M := N * 2

	//TODO: Finish expansion implementation
	fmt.Println("EXPANDING")
	fmt.Println(cf.buckets)

	cf.buckets = append(cf.buckets, cf.buckets...)
	cf.indices = append(cf.indices, cf.indices...)

	for i, bucket := range cf.buckets {
		for j, fp := range bucket {
			if fp == nil {
				continue
			}
			if cf.indices[i][j] {

			}
		}
	}

}
