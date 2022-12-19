package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

var (
	typ string

	root cobra.Command
	args []argument
)

func main() {
	var revs []bool
	var paths, seps []string

	root = cobra.Command{
		Use:   "d2",
		Short: "D2",
		Long:  "Compute delta between columns of 2 files of same class",
		Run: func(_ *cobra.Command, _ []string) {
			switch typ {
			case "i", "f":
			default:
				fmt.Fprintf(os.Stderr, "type: %s\n", typ)
				fmt.Fprintln(os.Stderr, "dst type specify: i(int) or f(float).")
				fmt.Fprintln(os.Stderr)
				root.Help()
				os.Exit(-1)
			}

			if len(paths) != 2 || len(seps) != 2 || len(revs) != 2 {
				fmt.Fprintf(os.Stderr, "paths: %+v\n", paths)
				fmt.Fprintf(os.Stderr, "seps: %+v\n", seps)
				fmt.Fprintf(os.Stderr, "revs: %+v\n", revs)
				fmt.Fprintln(os.Stderr, "there should be 3 pairs of params(path, sep, rev), each pair has two values.")
				fmt.Fprintln(os.Stderr)
				root.Help()
				os.Exit(-1)
			}

			for i := 0; i < 2; i++ {
				if fi, err := os.Stat(paths[i]); err == nil {
					if fi.IsDir() {
						fmt.Fprintf(os.Stderr, "the path should be a file instead of a directory, path: %s.\n", paths[i])
						fmt.Fprintln(os.Stderr)
						root.Help()
						os.Exit(-1)
					}
				} else if os.IsNotExist(err) {
					fmt.Fprintf(os.Stderr, "the path file not exists, path: %s.\n", paths[i])
					fmt.Fprintln(os.Stderr)
					root.Help()
					os.Exit(-1)
				} else {
					fmt.Fprintf(os.Stderr, "stats error, path: %s, err: %+v\n", paths[i], err)
					fmt.Fprintln(os.Stderr)
					root.Help()
					os.Exit(-1)
				}

				args = append(args, argument{
					paths[i],
					seps[i],
					revs[i],
				})
			}

			fmt.Fprintf(os.Stdout, "args: %+v\n", args)

			compute()
		},
	}

	root.Flags().StringVarP(&typ, "type", "t", "i", "dst type specify: i(int) or f(float)")

	root.Flags().StringSliceVarP(&paths, "path", "p", []string{}, "dst file path")
	root.Flags().StringSliceVarP(&seps, "sep", "s", []string{" ", " "}, "dst separator")
	root.Flags().BoolSliceVarP(&revs, "rev", "r", []bool{false, false}, "dst column sequence: false(forward) or true(reverse)")

	root.Execute()

}

type argument struct {
	path string
	sep  string
	rev  bool
}

type number interface {
	int64 | float64
}

// type numbers[N number] []N

type cache[N number] map[string]N

type assignFunc[N number] func([2][]byte, cache[N], parseFunc[N])

type parseFunc[N number] func(string) N

type dst[N number] struct {
	path string
	sep  []byte
	c    cache[N]
	afn  assignFunc[N]
	pfn  parseFunc[N]
}

type dsts[N number] []dst[N]

func (ds *dsts[N]) set(v dst[N]) {
	*ds = append(*ds, v)
}

func (ds dsts[N]) compute() {
	for sk, sv := range ds[0].c {
		dv, ok := ds[1].c[sk]
		if !ok {
			fmt.Printf("only in %s: %+v, k: %s\n", ds[0].path, sv, sk)
			continue
		}

		if sv > dv {
			fmt.Printf("src - dst: %+v, k: %s\n", sv-dv, sk)
		} else if sv < dv {
			fmt.Printf("dst - src: %+v, k: %s\n", dv-sv, sk)
		}
	}

	for k, v := range ds[1].c {
		_, ok := ds[0].c[k]
		if !ok {
			fmt.Printf("only in %s, k: %s, v: %+v\n", ds[1].path, k, v)
			continue
		}
	}
}

func compute() {
	if typ == "f" {
		var ds dsts[float64]

		for _, a := range args {
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
			fmt.Fprintf(os.Stdout, "path: %s, count: %d\n", a.path, len(d.c))
			if err != nil {
				fmt.Fprintf(os.Stderr, "read error, path: %s, err: %+v\n", a.path, err)
				os.Exit(-1)
			}

			ds.set(d)
		}

		ds.compute()
	} else {
		var ds dsts[int64]

		for _, a := range args {
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
			fmt.Fprintf(os.Stdout, "path: %s, count: %d\n", a.path, len(d.c))
			if err != nil {
				fmt.Fprintf(os.Stderr, "read error, path: %s, err: %+v\n", a.path, err)
				os.Exit(-1)
			}

			ds.set(d)
		}

		ds.compute()
	}
}

func parseI64[N number](s string) (n N) {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "row value can't convert to i64, s: %+v.\n", s)
		os.Exit(-1)
	}
	return N(i)
}

func parseF64[N number](s string) (n N) {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "row value can't convert to f64, s: %+v.\n", s)
		os.Exit(-1)
	}
	return N(f)
}

func retain[N number](bs [2][]byte, c cache[N], pfn parseFunc[N]) {
	c[string(bs[0])] = pfn(string(bs[1]))
}

func exchange[N number](bs [2][]byte, c cache[N], pfn parseFunc[N]) {
	c[string(bs[1])] = pfn(string(bs[0]))
}

func extract(b []byte, sep []byte) [2][]byte {
	bs := bytes.Split(b, sep)
	if len(bs) != 2 {
		fmt.Fprintf(os.Stderr, "extract row hasn't a pair of data, data: %+v, num: %d, sep: %+v.\n", b, len(bs), sep)
		os.Exit(-1)
	}
	return [2][]byte{bs[0], bs[1]}
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
		// fmt.Printf("b2: %s\n", b2)
		d.afn(b2, d.c, d.pfn)
	}

	return
}
