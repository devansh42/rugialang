package util

func IsDigit(b byte) bool {
	return b >= 48 && b <= 57
}

func IsUnderscore(b byte) bool {
	return b == '_'
}

func IsAlpha(b byte) bool {
	return (b >= 65 && b <= 90) || (b >= 97 && b <= 122)
}
