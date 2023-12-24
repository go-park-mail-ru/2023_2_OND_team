package levenstein

func Levenshtein(s1, s2 []rune) int {
	col := make([]int, len(s1)+1)

	for y := 1; y <= len(s1); y++ {
		col[y] = y
	}
	for x := 1; x <= len(s2); x++ {
		col[0] = x
		lastkey := x - 1
		for y := 1; y <= len(s1); y++ {
			oldkey := col[y]
			var incr int
			if s1[y-1] != s2[x-1] {
				incr = 1
			}

			col[y] = min(col[y]+1, col[y-1]+1, lastkey+incr)
			lastkey = oldkey
		}
	}
	return col[len(s1)]
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
