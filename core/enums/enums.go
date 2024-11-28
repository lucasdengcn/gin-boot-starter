package enums

import "strings"

type Gender int8

const (
	Unknown Gender = iota
	Male
	Female
)

func (g Gender) String() string {
	names := [...]string{"Unknown", "Male", "Female"}
	if g < Male || g > Female {
		return "Unknown"
	}
	return names[g]
}

func ParseGender(genderStr string) Gender {
	switch strings.ToLower(genderStr) {
	case "male":
		return Male
	case "female":
		return Female
	default:
		return Unknown
	}
}
