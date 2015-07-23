package cuckoofilter

import (
	"bufio"
	"fmt"
	"os"
	"testing"
)

func TestInsertion(t *testing.T) {

	cf := NewCuckooFilter(4)

	fd, err := os.Open("/usr/share/dict/web2")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	scanner := bufio.NewScanner(fd)

	fmt.Println("----")
	fmt.Println("Fingerprint = 1 byte and Bucketsize = 2 \nStart with 2 buckets\n")

	fmt.Println("Starting Insertion\n")

	values := make([][]byte, 0)
	i := 0
	for scanner.Scan() {
		s := []byte(scanner.Text())
		cf.Insert(s)
		values = append(values, s)
		i += 1
		if i == 4 {
			break
		}
	}

	count := cf.GetCount()
	if count != 4 {
		t.Errorf("Expected count = 4, instead count = %d", count)
	}

	fmt.Println("\n\nStarting Deletion\n")
	fmt.Print("Buckets: ")
	printBucket(cf.buckets)
	fmt.Print("Indicat: ")
	printIndicators(cf.indicators)
	fmt.Println("")

	for _, v := range values {
		cf.Delete(v)
	}

	fmt.Println("End Result\n")
	fmt.Print("Buckets: ")
	printBucket(cf.buckets)
	fmt.Print("Indicat: ")
	printIndicators(cf.indicators)
	fmt.Println("")

	count = cf.GetCount()
	if count != 0 {
		t.Errorf("Expected count = 0, instead count == %d", count)
	}
}
