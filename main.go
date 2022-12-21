package main

import (
	"github.com/spf13/cobra"
)

func main() {
	var typ string
	var l string
	var r []bool
	var p, s []string

	root := cobra.Command{
		Use:   "d2",
		Short: "D2",
		Long:  "Compute delta between columns of 2 files of same class",
		Run: func(_ *cobra.Command, _ []string) {
			tv, lv, pv, sv, rv := digitType([]byte(typ)[0]), level(l), paths(p), seps(s), revs(r)
			prs := []optioner{
				tv, &lv, pv, sv, rv,
			}

			opts.init()

			for _, pr := range prs {
				pr.parse()
				pr.with(&opts)
			}

			opts.print()

			compute()
		},
	}

	root.Flags().StringVarP(&typ, "type", "t", "i", "digit type specify: i(int) or f(float)")
	root.Flags().StringVarP(&l, "level", "l", "1|2|4|8", `print level (unit: bit)
1: argv[0] only 
2: argv[1] only 
4: argv[0] - argv[1] 
8: argv[1] - argv[0] 
`)

	root.Flags().StringArrayVarP(&p, "path", "p", []string{}, "file path")
	root.Flags().StringArrayVarP(&s, "sep", "s", []string{" ", " "}, "separator")
	root.Flags().BoolSliceVarP(&r, "rev", "r", []bool{false, false}, "column sequence: false(forward) or true(reverse)")

	root.Execute()
}
