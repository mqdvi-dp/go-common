package version

func compare(appVersion []int, minAppVersion []int) int {
	if len(appVersion) == 0 || len(minAppVersion) == 0 {
		return 0
	}

	a := appVersion[0]
	b := minAppVersion[0]

	if a > b {
		return 1
	} else if a < b {
		return -1
	}

	return compare(appVersion[1:], minAppVersion[1:])
}
