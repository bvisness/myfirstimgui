package imath

func Min(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func Max(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

func Abs(a int) int {
	if a < 0 {
		return -a
	} else {
		return a
	}
}

func LerpF(a, b int, t float32) int {
	return int(float32(a)*(1-t) + float32(b)*t)
}

func Lerp(a, b, tMin, tMax, t int) int {
	return LerpF(a, b, float32(t-tMin)/float32(tMax-tMin))
}
