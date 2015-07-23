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
	buckets    []bucket
	count      uint64
	indicators [][]uint
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
	indicators := make([][]uint, capacity)
	for i := range buckets {
		buckets[i] = make([]fingerprint, bucketSize)
		indicators[i] = make([]uint, bucketSize)
	}
	return &CuckooFilter{buckets, 0, indicators}
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

	limit := uint64(90*len(cf.buckets)*bucketSize) / 100

	if cf.count >= limit {
		cf.expand()
	}

	i1, i2, fp := getIndicesAndFingerprint(data, uint64(len(cf.buckets)))
	fmt.Println("INSERT >>>>>>", string(data), " fp:", fp, " i1:", i1, " i2:", i2)

	indicator := uint(0)
	if i1 != i2 {
		indicator = 1
	}

	if cf.insert(fp, i1, 0) || cf.insert(fp, i2, indicator) {
		fmt.Print("Buckets: ")
		printBucket(cf.buckets)
		fmt.Print("Indicat: ")
		printIndicators(cf.indicators)
		fmt.Println("")
		return true
	}
	res := cf.reinsert(fp, i2)
	fmt.Print("Buckets: ")
	printBucket(cf.buckets)
	fmt.Print("Indicat: ")
	printIndicators(cf.indicators)
	fmt.Println("")
	return res
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

func (cf *CuckooFilter) insert(fp fingerprint, i uint64, indicator uint) bool {
	if j := cf.buckets[i].insert(fp); j > -1 {
		cf.count++
		cf.indicators[i][j] = indicator
		return true
	}
	return false
}

func (cf *CuckooFilter) reinsert(fp fingerprint, i uint64) bool {
	indicator := uint(1)
	for k := 0; k < maxCuckooCount; k++ {
		j := rand.Intn(bucketSize)
		oldfp := fp
		oldIndicator := indicator
		fp = cf.buckets[i][j]
		if cf.indicators[i][j] == 0 {
			indicator = 1
		} else {
			indicator = 0
		}
		cf.buckets[i][j] = oldfp
		cf.indicators[i][j] = oldIndicator

		// look in the alternate location for that random element
		oldi := i
		i = getAltIndex(fp, i, uint64(len(cf.buckets)))
		if oldi == i {
			indicator = 0
		}
		if cf.insert(fp, i, indicator) {
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
	res := cf.delete(fp, i1) || cf.delete(fp, i2)
	fmt.Println("DELETE >>>>>>", string(data), " fp:", fp, " i1:", i1, " i2:", i2, " found:", res)
	fmt.Print("Buckets: ")
	printBucket(cf.buckets)
	fmt.Print("Indicat: ")
	printIndicators(cf.indicators)
	fmt.Println("")
	return res
}

func (cf *CuckooFilter) delete(fp fingerprint, i uint64) bool {
	if j := cf.buckets[i].delete(fp); j > -1 {
		cf.count--
		cf.indicators[i][j] = 0
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
	N := uint64(len(cf.buckets))
	M := N * 2

	origBucket := cf.buckets
	origInidcators := cf.indicators

	cf.buckets = make([]bucket, M)
	cf.indicators = make([][]uint, M)
	for i := range cf.buckets {
		cf.buckets[i] = make([]fingerprint, bucketSize)
		cf.indicators[i] = make([]uint, bucketSize)
	}

	//TODO: Finish expansion implementation
	fmt.Println("EXPANDING")
	fmt.Print("Original Buckets: ")
	printBucket(origBucket)
	fmt.Print("Original Indicat: ")
	printIndicators(origInidcators)

	for i := uint64(0); i < N; i++ {
		bucket := origBucket[i]
		for j, fp := range bucket {
			if fp == nil {
				continue
			}
			if origInidcators[i][j] == 0 {
				/*
					If fp is in the i1 bucket (i.e., the indicator bit is 0),
					then we still put this fingerprint in the i1 bucket in the new table.
				*/
				cf.buckets[i][j] = fp
				cf.indicators[i][j] = 0
			} else {
				/*
					However, if the it is in the i2 bucket in original table (i.e., indicator bit == 1),
					then we first calculate its i1 bucket index by i1 = i2 ^ hash(fingerprint) % N as in the paper,
					then calculate the new i2 bucket index in the new table by (i1 ^ hash(fingerprint)) % M,
					and move this fingerprint to the i2 bucket in the new table.
				*/
				i2 := i
				i1 := getAltIndex(fp, i, N)
				i2 = getAltIndex(fp, i1, M)
				cf.buckets[i2][j] = fp
				cf.indicators[i2][j] = 1
			}
		}
	}

	fmt.Print("Dedupled Buckets: ")
	printBucket(cf.buckets)
	fmt.Print("Dedupled Indicat: ")
	printIndicators(cf.indicators)
	fmt.Println("\n")
}
