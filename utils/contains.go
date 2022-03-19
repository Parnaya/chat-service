package utils

func ContainsAll(container []string, elements []string) bool {
	is := true
	for _, elem := range elements {
		is = is && Contains(container, elem)
	}
	return is
}

func Contains(container []string, element string) bool {
	for _, entry := range container {
		if entry == element {
			return true
		}
	}
	return false
}
