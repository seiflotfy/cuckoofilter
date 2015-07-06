package cuckoofilter

import (
	"bufio"
	"fmt"
	"os"
	"testing"
)

func TestInsertion(t *testing.T) {

	cf := NewCuckooFilter(1000000)

	fd, err := os.Open("/usr/share/dict/web2")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	scanner := bufio.NewScanner(fd)

	i := 0
	values := make([][]byte, 0)
	for scanner.Scan() {
		s := []byte(scanner.Text())
		cf.InsertUnique(s)
		i++
		values = append(values, s)
	}
	fmt.Println(cf.count, i)
	for _, v := range values {
		cf.Delete(v)
	}
	fmt.Println(cf.count, i)
}
