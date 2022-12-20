package main

type cache[N number] map[string]N

func (c cache[N]) set(k string, v N) {
	c[k] = v
}
