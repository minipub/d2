package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strconv"
	"strings"

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
				fmt.Fprintf(os.Stderr, "Error: invalid type: %s\n", typ)
				os.Exit(-1)
			}

			if len(paths) != 2 {
				f := `invalid paths: []`
				if len(paths) > 0 {
					f = fmt.Sprintf(`invalid paths: ["%+v"]`, strings.Join(paths, `","`))
				}
				fmt.Fprintln(os.Stderr, "Error:", f)
				os.Exit(-1)
			}

			if len(seps) != 2 {
				f := `invalid seps: []`
				if len(seps) > 0 {
					f = fmt.Sprintf(`invalid seps: ["%+v"]`, strings.Join(seps, `","`))
				}
				fmt.Fprintln(os.Stderr, "Error:", f)
				os.Exit(-1)
			}

			if len(revs) != 2 {
				fmt.Fprintf(os.Stderr, "Error: invalid revs: [%+v]\n", formatBool(revs))
				os.Exit(-1)
			}

			for i := 0; i < 2; i++ {
				if fi, err := os.Stat(paths[i]); err == nil {
					if fi.IsDir() {
						fmt.Fprintf(os.Stderr, "Error: the path should be a file instead of a directory, path: %s.\n", paths[i])
						os.Exit(-1)
					}
				} else if os.IsNotExist(err) {
					fmt.Fprintf(os.Stderr, "Error: the path file not exists, path: %s.\n", paths[i])
					os.Exit(-1)
				} else {
					fmt.Fprintf(os.Stderr, "Error: stats error, path: %s, err: %+v\n", paths[i], err)
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

	root.Flags().StringArrayVarP(&paths, "path", "p", []string{}, "dst file path")
	root.Flags().StringArrayVarP(&seps, "sep", "s", []string{" ", " "}, "dst separator")
	root.Flags().BoolSliceVarP(&revs, "rev", "r", []bool{false, false}, "dst column sequence: false(forward) or true(reverse)")

	root.Execute()

}

func formatBool(bs []bool) string {
	if len(bs) == 0 {
		return ""
	}

	var s []string
	for _, b := range bs {
		s = append(s, strconv.FormatBool(b))
	}
	return strings.Join(s, ",")
}

func parseExpr() {
	expr, err := parser.ParseExpr("1+2|4|8")
	if err != nil {
		fmt.Printf("err: %+v\n", err)
		return
	}

	ast.Print(nil, expr)

	rs, err := interpret(expr)
	if err != nil {
		fmt.Printf("err: %+v\n", err)
		return
	}

	fmt.Printf("rs: %d\n", rs)
}

func interpret(expr ast.Node) (int, error) {
	switch a := expr.(type) {
	case *ast.BinaryExpr:
		x, err := interpret(a.X)
		if err != nil {
			return -4, err
		}

		y, err := interpret(a.Y)
		if err != nil {
			return -5, err
		}

		switch a.Op {
		case token.OR:
			return x + y, nil
		default:
			return -6, fmt.Errorf("unknown op")
		}

	case *ast.BasicLit:
		switch a.Kind {
		case token.INT:
			i, err := strconv.Atoi(a.Value)
			if err != nil {
				return -2, err
			}
			return i, nil
		default:
			return -1, fmt.Errorf("unknown lit")
		}
	}

	return -3, fmt.Errorf("unknown type")
}
