package internal

import "os"

var powers = []int{1}

// GetPower returns the power of 95 with the given
// exponent.
func GetPower(exp int) int {
	if exp >= len(powers) {
		for i := len(powers); i <= exp; i++ {
			powers = append(powers, powers[i-1]*95)
		}
	}
	return powers[exp]
}

// ExtendState appends a character c to a state p with
// a depth n.
func ExtendState(n, p, s int) int {
	return (p%GetPower(n-1))*95 + s
}

// ResetDir prepares a directory. After calling, the given
// directory will exists and be empty.
func ResetDir(path string) {
	os.RemoveAll(path)
	os.MkdirAll(path, os.FileMode(0777))
}
