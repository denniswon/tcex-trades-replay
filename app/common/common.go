package common

import (
	"errors"
	"strconv"
)

// ParseNumber - Given an integer as string, attempts to parse it
func ParseNumber(number string) (uint64, error) {

	_num, err := strconv.ParseUint(number, 10, 64)
	if err != nil {

		return 0, errors.New("Failed to parse integer")

	}

	return _num, nil
}

// RangeChecker - Checks whether given number range is at max
// `limit` far away
func RangeChecker(from string, to string, limit uint64) (uint64, uint64, error) {

	_from, err := ParseNumber(from)
	if err != nil {
		return 0, 0, errors.New("Failed to parse integer")
	}

	_to, err := ParseNumber(to)
	if err != nil {
		return 0, 0, errors.New("Failed to parse integer")
	}

	if !(_to-_from < limit) {
		return 0, 0, errors.New("Range too long")
	}

	return _from, _to, nil
}
