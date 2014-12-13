package main

type key struct {
    Command string
    Help    string
}

type keyMap map[int]key

func (k keyMap) Merge(m keyMap) {
    for code, key := range m {
        k[code] = key
    }
}
