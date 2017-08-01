package main

import (
	"fmt"
	"math"
)

type bucket [4]uint16

func main() {
	for i := uint16(0); i <= 49; i++ {
		dict := map[bucket]bool{}
		for a := uint16(0); a < i; a++ {
			for b := a; b < i; b++ {
				for c := b; c < i; c++ {
					for d := c; d < i; d++ {
						x := bucket([4]uint16{a, b, c, d})
						dict[x] = true
					}
				}
			}
		}
		fmt.Println(i, math.Pow(2, float64(i)), len(dict))
	}
}
