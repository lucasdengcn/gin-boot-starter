package core

import (
	"strconv"
)

// StringFromUint
func StringFromUint(val uint) string {
	return strconv.Itoa(int(val))
}

// UintFromString
func UintFromString(val string) (uint, error) {
	intVal, err := strconv.Atoi(val)
	if err != nil {
		return 0, err
	}
	return uint(intVal), nil
}
