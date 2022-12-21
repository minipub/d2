package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
)

type dst[N number] struct {
	path string
	sep  []byte
	c    cache[N]
	afn  assignFunc[N]
	pfn  parseFunc[N]
}

func (d dst[_]) read() (err error) {
	f, err := os.Open(d.path)
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

		b2 := extract(b, d.sep)
		d.afn(b2, d.c, d.pfn)
	}

	return
}

type dsts[N number] []dst[N]

func (ds *dsts[N]) set(v dst[N]) {
	*ds = append(*ds, v)
}

func (ds dsts[N]) compute() {
	for fk, fv := range ds[0].c {
		bv, ok := ds[1].c[fk]
		if !ok {
			if opts.isOnly1st() {
				fmt.Fprintf(os.Stdout, "[argv[1st] only] k: %s, v: %+v\n", fk, fv)
			}
			continue
		}

		if fv > bv {
			if opts.isSub1st2nd() {
				fmt.Fprintf(os.Stdout, "[argv[1st](%+v) - argv[2nd](%+v)] k: %s, v: %+v\n", fv, bv, fk, fv-bv)
			}
		} else if fv < bv {
			if opts.isSub2nd1st() {
				fmt.Fprintf(os.Stdout, "[argv[2nd](%+v) - argv[1st](%+v)] k: %s, v: %+v\n", bv, fv, fk, bv-fv)
			}
		}
	}

	for k, v := range ds[1].c {
		_, ok := ds[0].c[k]
		if !ok {
			if opts.isOnly2nd() {
				fmt.Fprintf(os.Stdout, "[argv[2nd] only] k: %s, v: %+v\n", k, v)
			}
			continue
		}
	}
}

func compute() {
	if opts.typ == 'f' {
		var ds dsts[float64]

		for _, a := range opts.posOpts {
			afn := retain[float64]
			if a.rev {
				afn = exchange[float64]
			}

			d := dst[float64]{
				a.path,
				[]byte(a.sep),
				make(map[string]float64),
				afn,
				parseF64[float64],
			}

			err := d.read()
			if opts.isVerbose() {
				fmt.Fprintf(os.Stdout, "path: %s, line count: %d\n", a.path, len(d.c))
			}
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: read error, path: %s, err: %+v\n", a.path, err)
				os.Exit(-1)
			}

			ds.set(d)
		}

		ds.compute()
	} else {
		var ds dsts[int64]

		for _, a := range opts.posOpts {
			afn := retain[int64]
			if a.rev {
				afn = exchange[int64]
			}

			d := dst[int64]{
				a.path,
				[]byte(a.sep),
				make(map[string]int64),
				afn,
				parseI64[int64],
			}

			err := d.read()
			if opts.isVerbose() {
				fmt.Fprintf(os.Stdout, "path: %s, line count: %d\n", a.path, len(d.c))
			}
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: read error, path: %s, err: %+v\n", a.path, err)
				os.Exit(-1)
			}

			ds.set(d)
		}

		ds.compute()
	}
}

func extract(b []byte, sep []byte) [2][]byte {
	bs := bytes.Split(b, sep)
	if len(bs) != 2 {
		fmt.Fprintf(os.Stderr, "Error: extract row hasn't a pair of data, data: %+v, num: %d, sep: %+v.\n", b, len(bs), sep)
		os.Exit(-1)
	}
	return [2][]byte{bs[0], bs[1]}
}
