package storage

func ExtractFileFormat(filename string) string {
	lastDotIndex := -1
	for i := len(filename) - 1; i >= 0; i-- {
		if filename[i] == '.' {
			lastDotIndex = i
			break
		}
	}

	if lastDotIndex == -1 || lastDotIndex == len(filename)-1 {
		return ""
	}

	return filename[lastDotIndex+1:]
}
