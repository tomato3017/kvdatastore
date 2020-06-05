package util

func RemoveSliceValue(valSlice []string, val string) ([]string, bool) {
	for i, v := range valSlice {
		if v == val {
			valSlice = append(valSlice[:i], valSlice[i+1:]...)
			return valSlice, true
		}
	}

	return valSlice, false
}

func EqualStringSlices(val1 []string, val2 []string) bool {
	if len(val1) != len(val2) {
		return false
	}

	tmpVal2 := make([]string, len(val2))
	copy(tmpVal2, val2)

	for _, v := range val1 {
		var ok bool
		tmpVal2, ok = RemoveSliceValue(tmpVal2, v)
		if !ok {
			return false
		}
	}

	if len(tmpVal2) != 0 {
		return false
	}

	return true
}
