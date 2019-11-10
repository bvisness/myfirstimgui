package imath

func MinInt(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func MaxInt(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

func AbsInt(a int) int {
	if a < 0 {
		return -a
	} else {
		return a
	}
}
