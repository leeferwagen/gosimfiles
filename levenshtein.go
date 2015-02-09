package main

type levenshtein struct {
	refstr string
	reflen int
	column []int
}

func CreateLevenshtein(refstr string) *levenshtein {
	reflen := len(refstr)
	column := make([]int, reflen+1)
	for i := 1; i <= reflen; i++ {
		column[i] = i
	}
	return &levenshtein{
		refstr: refstr,
		reflen: reflen,
		column: column,
	}
}

// Shamelessly stolen from https://github.com/arbovm/levenshtein
func (l *levenshtein) Distance(s string) (int, float64) {
	if l.refstr == s {
		return 0, 100.0
	}

	var cost, lastdiag, olddiag int
	slen := len(s)

	column := make([]int, len(l.column))
	copy(column, l.column)

	for x := 1; x <= slen; x++ {
		column[0] = x
		lastdiag = x - 1
		for y := 1; y <= l.reflen; y++ {
			olddiag = column[y]
			cost = 0
			if l.refstr[y-1] != s[x-1] {
				cost = 1
			}
			column[y] = min(column[y]+1, column[y-1]+1, lastdiag+cost)
			lastdiag = olddiag
		}
	}
	distance := column[l.reflen]
	var similarity float64
	if distance > 0 {
		similarity = 100.0 - (100.0 / (float64(maxInt(l.reflen, slen)) / float64(distance)))
	} else {
		similarity = 0.0
	}
	return distance, similarity
}

func min(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
	} else {
		if b < c {
			return b
		}
	}
	return c
}
