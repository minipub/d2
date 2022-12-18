package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
)

var (
	srcFile = "../src.out"
	dstFile = "../dst.out"
)

type assignFunc func([2][]byte, map[string]int64)

func main() {
	srcM := make(map[string]int64)
	{
		err := read(srcFile, srcM, []byte("\t"), retain)
		fmt.Printf("%+v\n", len(srcM))
		fmt.Printf("err: %+v\n", err)
	}

	dstM := make(map[string]int64)
	{
		err := read(dstFile, dstM, []byte(" "), exchange)
		fmt.Printf("%+v\n", len(dstM))
		fmt.Printf("err: %+v\n", err)
	}

	for sk, sv := range srcM {
		dv, ok := dstM[sk]
		if !ok {
			fmt.Printf("only src: %d, k: %s\n", sv, sk)
			continue
		}

		if sv > dv {
			fmt.Printf("src - dst: %d, k: %s\n", sv-dv, sk)
		} else if sv < dv {
			fmt.Printf("dst - src: %d, k: %s\n", dv-sv, sk)
		}
	}

	for dk, dv := range dstM {
		_, ok := srcM[dk]
		if !ok {
			fmt.Printf("only dst: %d, k: %s\n", dv, dk)
			continue
		}
	}
}

func retain(bs [2][]byte, m map[string]int64) {
	i, err := strconv.ParseInt(string(bs[1]), 10, 64)
	if err != nil {
		panic(string(bs[1]))
	}
	m[string(bs[0])] = i
}

func exchange(bs [2][]byte, m map[string]int64) {
	i, err := strconv.ParseInt(string(bs[0]), 10, 64)
	if err != nil {
		panic(string(bs[0]))
	}
	m[string(bs[1])] = i
}

func extract(b []byte, delimiter []byte) [2][]byte {
	bs := bytes.Split(b, delimiter)
	if len(bs) != 2 {
		panic(len(bs))
	}
	return [2][]byte{bs[0], bs[1]}
}

func read(path string, m map[string]int64, delimiter []byte, fn assignFunc) (err error) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	r := bufio.NewReader(f)

	for {
		b, _, err := r.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		b2 := extract(b, delimiter)
		// fmt.Printf("b2: %s\n", b2)
		fn(b2, m)
	}

	return
}
