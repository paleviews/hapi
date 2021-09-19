package runtime

import (
	"strconv"
)

func ParseBool(s string) (bool, error) {
	return strconv.ParseBool(s)
}

func ParseInt32(s string) (int32, error) {
	i, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return 0, err
	}
	return int32(i), nil
}

func ParseInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

func ParseUint32(s string) (uint32, error) {
	i, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint32(i), nil
}

func ParseUint64(s string) (uint64, error) {
	return strconv.ParseUint(s, 10, 64)
}

func ParseFloat32(s string) (float32, error) {
	f, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return 0, err
	}
	return float32(f), nil
}

func ParseFloat64(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}
