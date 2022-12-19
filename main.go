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

func init() {
	var revs []bool
	var paths, seps []string

	root = cobra.Command{
		Use:   "",
		Short: "D2",
		Long:  "Compute delta between columns of 2 files of same class",
		Run: func(_ *cobra.Command, _ []string) {
			switch typ {
			case "i", "f":
			default:
				fmt.Fprintf(os.Stderr, "type: %s\n", typ)
				fmt.Fprintln(os.Stderr, "dst type specify: i(int) or f(float)")
				os.Exit(-1)
			}

			if len(paths) != 2 || len(seps) != 2 || len(revs) != 2 {
				fmt.Fprintf(os.Stderr, "paths: %+v\n", paths)
				fmt.Fprintf(os.Stderr, "seps: %+v\n", seps)
				fmt.Fprintf(os.Stderr, "revs: %+v\n", revs)
				fmt.Fprintln(os.Stderr, "there should be 3 pairs of params(path, sep, rev), each pair has two values.")
				os.Exit(-1)
			}

			for i := 0; i < 2; i++ {
				if fi, err := os.Stat(paths[i]); err == nil {
					if fi.IsDir() {
						fmt.Fprintf(os.Stderr, "the path should be a file instead of a directory, path: %s.\n", paths[i])
						os.Exit(-1)
					}
				} else if os.IsNotExist(err) {
					fmt.Fprintf(os.Stderr, "the path file not exists, path: %s.\n", paths[i])
					os.Exit(-1)
				} else {
					fmt.Fprintf(os.Stderr, "stats error, path: %s, err: %+v\n", paths[i], err)
					os.Exit(-1)
				}

				args = append(args, argument{
					paths[i],
					seps[i],
					revs[i],
				})
			}

			main()
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

type assignFunc[N number] func([2][]byte, cache[N])

type dst[N number] struct {
	path string
	sep  []byte
	c    cache[N]
	fn   assignFunc[N]
}

type dsts[N number] []dst[N]

func (ds dsts[N]) set(v dst[N]) {
	ds = append(ds, v)
}

func (ds dsts[N]) compute() {
	for sk, sv := range ds[0].c {
		dv, ok := ds[1].c[sk]
		if !ok {
			fmt.Printf("only in %s: %d, k: %s\n", ds[0].path, sv, sk)
			continue
		}

		if sv > dv {
			fmt.Printf("src - dst: %d, k: %s\n", sv-dv, sk)
		} else if sv < dv {
			fmt.Printf("dst - src: %d, k: %s\n", dv-sv, sk)
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

func main() {
	if typ == "f" {
		var ds dsts[float64]

		for _, a := range args {
			d := dst[float64]{
				a.path,
				[]byte(a.sep),
				make(map[string]float64),
				exchange[float64],
			}
			err := d.read()
			fmt.Printf("%+v\n", len(d.c))
			fmt.Printf("err: %+v\n", err)

			ds = append(ds, d)
		}

		ds.compute()
	} else {
		var ds dsts[int64]

		for _, a := range args {
			d := dst[int64]{
				a.path,
				[]byte(a.sep),
				make(map[string]int64),
				retain[int64],
			}
			err := d.read()
			fmt.Printf("%+v\n", len(d.c))
			fmt.Printf("err: %+v\n", err)

			ds = append(ds, d)
		}

		ds.compute()
	}
}

func retain[N number](bs [2][]byte, m cache[N]) {
	i, err := strconv.ParseInt(string(bs[1]), 10, 64)
	if err != nil {
		panic(string(bs[1]))
	}
	m[string(bs[0])] = N(i)
}

func exchange[N number](bs [2][]byte, m cache[N]) {
	i, err := strconv.ParseInt(string(bs[0]), 10, 64)
	if err != nil {
		panic(string(bs[0]))
	}
	m[string(bs[1])] = N(i)
}

func extract(b []byte, sep []byte) [2][]byte {
	bs := bytes.Split(b, sep)
	if len(bs) != 2 {
		panic(len(bs))
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
		d.fn(b2, d.c)
	}

	return
}
