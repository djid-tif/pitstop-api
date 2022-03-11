package utils

import (
	"strconv"
	"time"
)

func TimeToIso(t time.Time) string {
	return t.Format(time.RFC3339)
}

func StringToUint(str string) (uint, error) {
	res, err := strconv.ParseUint(str, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(res), nil
}

func ContainsUint(i uint, l []uint) bool {
	for _, j := range l {
		if i == j {
			return true
		}
	}

	return false
}

// RemoveUintFromList iterates through the given list and tries to remove the given element i.
// It returns whether it removed the element from the list or not
func RemoveUintFromList(i uint, l []uint) ([]uint, bool) {
	for index, j := range l {
		if i == j {
			return append(l[:index], l[index+1:]...), true
		}
	}

	return l, false
}
