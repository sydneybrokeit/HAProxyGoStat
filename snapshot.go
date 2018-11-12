package HAProxyGoStat

import "time"

// A single snapshot, a collection of individual stats.
type HAProxyStatSnapshot struct {
	Created time.Time
	Stats []HAProxyStat
}

// takes a function, f, with a boolean return, and filters to those values that return True
func (stats *HAProxyStatSnapshot) Filter(f func(HAProxyStat) bool) *HAProxyStatSnapshot {
	snapshot := GenerateNewSnapshot()
	for _, stat := range stats.Stats {
		if f(stat) {
			snapshot.Stats = append(snapshot.Stats, stat)
		}
	}
	return snapshot
}

// Instantiate a new Snapshot
func GenerateNewSnapshot() *HAProxyStatSnapshot {
	snapshot := new(HAProxyStatSnapshot)
	snapshot.Stats = make([]HAProxyStat, 0)
	return snapshot
}
