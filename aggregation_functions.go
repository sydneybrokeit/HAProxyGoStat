package HAProxyGoStat

func sumInt(stats []int) int {
	total := 0
	for _, stat := range stats {
		total += stat
	}

	return total
}

func averageInt(stats []int) float32 {
	sum := float32(sumInt(stats))
	average := sum / float32(len(stats))
	return average
}

func maxInt(stats []int) int {
	var max int
	for _, stat := range stats {
		if stat > max {
			max = stat
		}
	}
	return max
}
