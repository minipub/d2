package main

const (
	// 1: argv[0] only
	only1st = 1 << iota
	// 2: argv[1] only
	only2nd
	// 4: argv[0] - argv[1]
	sub1st2nd
	// 8: argv[1] - argv[0]
	sub2nd1st
	// 16: verbose
	verbose
)

type intLevel uint8

func (l intLevel) isOnly1st() bool {
	return l&only1st == only1st
}

func (l intLevel) isOnly2nd() bool {
	return l&only2nd == only2nd
}

func (l intLevel) isSub1st2nd() bool {
	return l&sub1st2nd == sub1st2nd
}

func (l intLevel) isSub2nd1st() bool {
	return l&sub2nd1st == sub2nd1st
}

func (l intLevel) isVerbose() bool {
	return l&verbose == verbose
}
