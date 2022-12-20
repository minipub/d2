package main

import (
	"fmt"
	"os"
	"strconv"
)

type parseFunc[N number] func(string) N

func parseI64[N number](s string) (n N) {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: row value can't convert to i64, s: %+v.\n", s)
		os.Exit(-1)
	}
	return N(i)
}

func parseF64[N number](s string) (n N) {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: row value can't convert to f64, s: %+v.\n", s)
		os.Exit(-1)
	}
	return N(f)
}
