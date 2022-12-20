package main

type assignFunc[N number] func([2][]byte, cache[N], parseFunc[N])

func retain[N number](bs [2][]byte, c cache[N], pfn parseFunc[N]) {
	c.set(string(bs[0]), pfn(string(bs[1])))
}

func exchange[N number](bs [2][]byte, c cache[N], pfn parseFunc[N]) {
	c.set(string(bs[1]), pfn(string(bs[0])))
}
