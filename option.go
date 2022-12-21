package main

import (
	"fmt"
	"go/ast"
	goparser "go/parser"
	"go/token"
	"os"
	"strconv"
	"strings"
)

var opts option

type posOpt struct {
	path string
	sep  string
	rev  bool
}

type posOpts []posOpt

type option struct {
	typ   byte
	level uint8
	posOpts
}

func (p *option) init() {
	p.posOpts = make(posOpts, 0, 2)
}

func (p *option) print() {
	fmt.Fprintf(os.Stdout, "opts: %+v\n", p)
}

type optioner interface {
	parser
	wither
}

type parser interface {
	parse()
}

type wither interface {
	with(*option)
}

type digitType byte

func (dt digitType) parse() {
	switch dt {
	case 'i', 'f':
	default:
		fmt.Fprintf(os.Stderr, "Error: invalid digit type: %c\n", dt)
		os.Exit(-1)
	}
}

func (dt digitType) with(opt *option) {
	opt.typ = byte(dt)
}

type paths []string

func (p paths) parse() {
	if len(p) != 2 {
		f := `invalid paths: []`
		if len(p) > 0 {
			f = fmt.Sprintf(`invalid paths: ["%+v"]`, strings.Join(p, `","`))
		}
		fmt.Fprintln(os.Stderr, "Error:", f)
		os.Exit(-1)
	}

	for i := 0; i < 2; i++ {
		if fi, err := os.Stat(p[i]); err == nil {
			if fi.IsDir() {
				fmt.Fprintf(os.Stderr, "Error: the path should be a file instead of a directory, path: %s.\n", p[i])
				os.Exit(-1)
			}
		} else if os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Error: the path file not exists, path: %s.\n", p[i])
			os.Exit(-1)
		} else {
			fmt.Fprintf(os.Stderr, "Error: stats error, path: %s, err: %+v\n", p[i], err)
			os.Exit(-1)
		}
	}
}

func (p paths) with(opt *option) {
	for i := 0; i < 2; i++ {
		opts.posOpts[i].path = p[i]
	}
}

type seps []string

func (s seps) parse() {
	if len(s) != 2 {
		f := `invalid seps: []`
		if len(s) > 0 {
			f = fmt.Sprintf(`invalid seps: ["%+v"]`, strings.Join(s, `","`))
		}
		fmt.Fprintln(os.Stderr, "Error:", f)
		os.Exit(-1)
	}
}

func (s seps) with(opt *option) {
	for i := 0; i < 2; i++ {
		opts.posOpts[i].sep = s[i]
	}
}

type revs []bool

func (r revs) parse() {
	if len(r) != 2 {
		fmt.Fprintf(os.Stderr, "Error: invalid revs: [%+v]\n", formatBool(r))
		os.Exit(-1)
	}
}

func (r revs) with(opt *option) {
	for i := 0; i < 2; i++ {
		opts.posOpts[i].rev = r[i]
	}
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

type level string

func (l *level) parse() {
	expr, err := goparser.ParseExpr(string(*l))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: invalid level, err: [%+v]\n", err)
		os.Exit(-1)
	}
	// ast.Print(nil, expr)

	i, err := interpret(expr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: interpret level, err: [%+v]\n", err)
		os.Exit(-1)
	}
	// fmt.Printf("i: %d\n", i)

	*l = level(strconv.Itoa(i))
}

func (l *level) with(opt *option) {
	i, _ := strconv.Atoi(string(*l))
	opts.level = uint8(i)
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
