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

	values := make([][]byte, 0)
	i := 0
	for scanner.Scan() {
		s := []byte(scanner.Text())
		cf.InsertUnique(s)
		values = append(values, s)
		i += 1
		if i == 5 {
			break
		}
	}

	count := cf.GetCount()
	if count != 3 {
		t.Errorf("Expected count = 235041, instead count = %d", count)
	}

	fmt.Println(cf.buckets)

	for _, v := range values {
		cf.Delete(v)
	}

	count = cf.GetCount()
	if count != 0 {
		t.Errorf("Expected count = 0, instead count == %d", count)
	}
}
