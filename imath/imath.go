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

func LerpIntF(a, b int, t float32) int {
	return int(float32(a)*(1-t) + float32(b)*t)
}

func LerpInt(a, b, tMin, tMax, t int) int {
	return LerpIntF(a, b, float32(t-tMin)/float32(tMax-tMin))
}
