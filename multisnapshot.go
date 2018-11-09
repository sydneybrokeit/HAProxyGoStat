package HAProxyGoStat

import "github.com/oleiade/reflections"

var HAProxyFieldAggregationFunctions = map[string]string {

}

// filters all snapshots in a multisnapshot
func (multisnapshot *HAProxyMultiSnapshot) Filter(f func(HAProxyStat) bool) *HAProxyMultiSnapshot {
	newSnapshot := GenerateNewMultiSnapshot()
	for _, snapshot := range multisnapshot.Snapshots {
		filtered := snapshot.Filter(f)
		newSnapshot.Snapshots = append(newSnapshot.Snapshots, *filtered)
	}
	return newSnapshot
}

// a collection of multiple snapshots; useful for multiple sockets
type HAProxyMultiSnapshot struct {
	Snapshots []HAProxyStatSnapshot
}

//instantiate a new multisnapshot
func GenerateNewMultiSnapshot() *HAProxyMultiSnapshot {
	snapshot := new(HAProxyMultiSnapshot)
	snapshot.Snapshots = make([]HAProxyStatSnapshot, 0)
	return snapshot
}

// Aggregate a multisnapshot into a single snapshot
func (multisnapshot *HAProxyMultiSnapshot) Aggregate() *HAProxyStatSnapshot {
	aggregate := GenerateNewSnapshot()
	fields, err := reflections.Fields(multisnapshot.Snapshots[0].Stats[0])
	if err != nil {
		panic(err.Error())
	}
	for _, attr := range fields {
		switch HAProxyFieldAggregationFunctions[attr] {
		case "sum":
			ints := make([]int, 0)
			for i, snapshot := range multisnapshot.Snapshots {
				innerStat, err := reflections.GetField(snapshot.Stats[i], attr)
				if err != nil {
					panic(err.Error())
				}
				innerStatInt := innerStat.(int)
				ints = append(ints, innerStatInt)
			}
			sumOfSnapshots := sumInt(ints)
			err := reflections.SetField(&aggregate.Stats, attr, sumOfSnapshots)
			if err != nil {
				panic(err.Error())
			}
		case "average":
			ints := make([]int, 0)
			for i, snapshot := range multisnapshot.Snapshots {
				innerStat, err := reflections.GetField(snapshot.Stats[i], attr)
				if err != nil {
					panic(err.Error())
				}
				innerStatInt := innerStat.(int)
				ints = append(ints, innerStatInt)
			}
			averageOfSnapshots := int(averageInt(ints))
			err := reflections.SetField(&aggregate.Stats, attr, averageOfSnapshots)
			if err != nil {
				panic(err.Error())
			}
		case "max":
			ints := make([]int, 0)
			for i, snapshot := range multisnapshot.Snapshots {
				innerStat, err := reflections.GetField(snapshot, attr)
				if err != nil {
					panic(err.Error())
				}
				innerStatInt := innerStat.(int)
				ints = append(ints, innerStatInt)
			}
			maxOfSnapshots := maxInt(ints)
			err := reflections.SetField(&aggregate.Stats, attr, maxOfSnapshots)
			if err != nil {
				panic(err.Error())
			}
		case "passthrough":
			passthrough, err := reflections.GetField(multisnapshot.Snapshots[0], att)
			err := reflections.SetField(&aggregate, )
		}
	}
	return aggregate
}
