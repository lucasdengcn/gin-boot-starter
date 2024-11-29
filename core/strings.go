package core

import (
	"strconv"
)

func StringFromUint(val uint) string {
	return strconv.Itoa(int(val))
}
