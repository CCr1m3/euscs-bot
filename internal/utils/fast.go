package utils

var _gSeed int = 644660825

func FastRand() int { // this was haashi's idea
	_gSeed = (214013)*_gSeed + 2531011
	return (_gSeed >> 16) & 0x7FFF
}

func FastRandN(n int) int {
	return FastRand() % n
}
