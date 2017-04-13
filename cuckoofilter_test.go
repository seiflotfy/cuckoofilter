package cuckoofilter

import (
	"bufio"
	"fmt"
	"os"
	"testing"
)

func TestInsertion(t *testing.T) {

	cf := NewCuckooFilter(1000000)

	fd, err := os.Open("/usr/share/dict/words")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	scanner := bufio.NewScanner(fd)

	var values [][]byte
	var lineCount uint
	for scanner.Scan() {
		s := []byte(scanner.Text())
		if cf.InsertUnique(s) {
			lineCount++
		}
		values = append(values, s)
	}

	count := cf.Count()
	if count != lineCount {
		t.Errorf("Expected count = %d, instead count = %d", lineCount, count)
	}

	for _, v := range values {
		cf.Delete(v)
	}

	count = cf.Count()
	if count != 0 {
		t.Errorf("Expected count = 0, instead count == %d", count)
	}
}

func TestSaveLoad(t *testing.T) {
	defer os.Remove("./test.gob")
	cf := NewCuckooFilter(1000000)
	cf.InsertUnique([]byte("test1"))
	cf.Save("./test.gob")
	cf.InsertUnique([]byte("test2"))
	cf.Load("./test.gob")
	if !cf.Lookup([]byte("test1")) {
		t.Errorf("Expected 'test1' in filter")
	}
	if cf.Lookup([]byte("test2")) {
		t.Errorf("Not expected 'test2' in filter")
	}
	count := cf.Count()
	if count != 1 {
		t.Errorf("Expected count = 1, instead count == %d", count)
	}
}
