package cuckoofilter

import (
	"encoding/gob"
	"math/rand"
	"os"
)

const maxCuckooCount = 500

/*
CuckooFilter represents a probabalistic counter
*/
type CuckooFilter struct {
	Buckets   []bucket
	ItemCount uint
}

/*
NewCuckooFilter returns a new cuckoofilter with a given capacity
*/
func NewCuckooFilter(capacity uint) *CuckooFilter {
	capacity = getNextPow2(uint64(capacity)) / bucketSize
	if capacity == 0 {
		capacity = 1
	}
	buckets := make([]bucket, capacity, capacity)
	for i := range buckets {
		buckets[i] = [bucketSize]fingerprint{}
	}
	return &CuckooFilter{buckets, 0}
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
	i1, i2, fp := getIndicesAndFingerprint(data, uint(len(cf.Buckets)))
	b1, b2 := cf.Buckets[i1], cf.Buckets[i2]
	return b1.getFingerprintIndex(fp) > -1 || b2.getFingerprintIndex(fp) > -1
}

/*
Insert inserts data into the counter and returns true upon success
*/
func (cf *CuckooFilter) Insert(data []byte) bool {
	i1, i2, fp := getIndicesAndFingerprint(data, uint(len(cf.Buckets)))
	if cf.insert(fp, i1) || cf.insert(fp, i2) {
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

func (cf *CuckooFilter) insert(fp fingerprint, i uint) bool {
	if cf.Buckets[i].insert(fp) {
		cf.ItemCount++
		return true
	}
	return false
}

func (cf *CuckooFilter) reinsert(fp fingerprint, i uint) bool {
	for k := 0; k < maxCuckooCount; k++ {
		j := rand.Intn(bucketSize)
		oldfp := fp
		fp = cf.Buckets[i][j]
		cf.Buckets[i][j] = oldfp

		// look in the alternate location for that random element
		i = getAltIndex(fp, i, uint(len(cf.Buckets)))
		if cf.insert(fp, i) {
			return true
		}
	}
	return false
}

/*
Delete data from counter if exists and return if deleted or not
*/
func (cf *CuckooFilter) Delete(data []byte) bool {
	i1, i2, fp := getIndicesAndFingerprint(data, uint(len(cf.Buckets)))
	return cf.delete(fp, i1) || cf.delete(fp, i2)
}

func (cf *CuckooFilter) delete(fp fingerprint, i uint) bool {
	if cf.Buckets[i].delete(fp) {
		cf.ItemCount--
		return true
	}
	return false
}

/*
Count returns the number of items in the counter
*/
func (cf *CuckooFilter) Count() uint {
	return cf.ItemCount
}

/*
Save the filter to a file
*/
func (cf *CuckooFilter) Save(path string) error {
	file, err := os.Create(path)
	defer file.Close()
	if err != nil {
		return err
	}
	encoder := gob.NewEncoder(file)
	err = encoder.Encode(cf)
	if err != nil {
		return err
	}
	return nil
}

/*
Load a filter from a file
*/
func (cf *CuckooFilter) Load(path string) error {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return err
	}
	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&cf)
	if err != nil {
		return err
	}
	return nil
}
