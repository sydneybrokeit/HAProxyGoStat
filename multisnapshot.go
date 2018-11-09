package HAProxyGoStat

import (
	"fmt"
	"github.com/oleiade/reflections"
)

var HAProxyFieldAggregationFunctions = map[string]string{
	"ProxyName":          "string",
	"ServerName":         "string",
	"CurrentQueue":       "sum",
	"MaxQueue":           "max",
	"CurrentSessions":    "sum",
	"MaxSessions":        "max",
	"SessionLimit":       "max",
	"TotalSessions":      "sum",
	"BytesIn":            "sum",
	"BytesOut":           "sum",
	"DeniedRequests":     "sum",
	"DeniedResponses":    "sum",
	"ErrorRequests":      "sum",
	"ErrorConnections":   "sum",
	"ErrorResponses":     "sum",
	"ServerRetry":        "sum",
	"ServerRedispatch":   "sum",
	"Status":             "string",
	"Weight":             "int",
	"ActiveServers":      "max",
	"BackupServers":      "max",
	"FailedChecks":       "sum",
	"DownTransitions":    "sum",
	"LastTransition":     "sum",
	"Downtime":           "sum",
	"QueueLimit":         "max",
	"PID":                "pass",
	"IID":                "pass",
	"ServerID":           "pass",
	"Throttle":           "average",
	"LoadBalanceTotal":   "sum",
	"Tracked":            "pass",
	"Type":               "string",
	"Rate":               "sum",
	"RateLimit":          "max",
	"RateMax":            "max",
	"CheckStatus":        "string",
	"CheckCode":          "pass",
	"CheckDuration":      "average",
	"HTTP1xx":            "sum",
	"HTTP2xx":            "sum",
	"HTTP3xx":            "sum",
	"HTTP4xx":            "sum",
	"HTTP5xx":            "sum",
	"HTTPOther":          "sum",
	"HANAFail":           "string",
	"HTTPRequestRate":    "sum",
	"HTTPRequestRateMax": "max",
	"HTTPRequestsTotal":  "sum",
	"ClientAbort":        "sum",
	"ServerAbort":        "sum",
	"CompressorIn":       "sum",
	"CompressorOut":      "sum",
	"CompressorBypass":   "sum",
	"CompressorResponse": "sum",
	"LastSession":        "max",
	"LastCheck":          "string",
	"LastAgent":          "string",
	"QueueTime":          "average",
	"ConnectTime":        "average",
	"ResponseTime":       "average",
	"TotalTime":          "average",
	"AgentStatus":        "string",
	"AgentCode":          "pass",
	"AgentDuration":      "average",
	"CheckDescription":   "string",
	"AgentDescription":   "string",
	"CheckRise":          "sum",
	"CheckFall":          "sum",
	"CheckHealth":        "sum",
	"AgentRise":          "sum",
	"AgentFall":          "sum",
	"AgentHealth":        "sum",
	"Address":            "string",
	"Cookie":             "string",
	"Mode":               "string",
	"Algorithm":          "string",
	"ConnectionRate":     "sum",
	"ConnectionRateMax":  "max",
	"ConnectionTotal":    "sum",
	"Intercepted":        "sum",
	"DeniedConnections":  "sum",
	"DeniedSessions":     "sum",
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
	snapshotLen := len(multisnapshot.Snapshots)
	for i, _ := range multisnapshot.Snapshots[0].Stats {
		aggregate.Stats = append(aggregate.Stats, *new(HAProxyStat))
		for _, attr := range fields {
			fieldType, err := reflections.GetFieldType(multisnapshot.Snapshots[0].Stats[0], attr)
			if err != nil {
				panic(err.Error())
			}
			if fieldType == "int" {
				ints := make([]int, 0)
				for x := 0; x < snapshotLen; x++ {
					inner, err := reflections.GetField(multisnapshot.Snapshots[x].Stats[i], attr)
					if err != nil {
						panic(err.Error())
					}
					innerInt := inner.(int)
					ints = append(ints, innerInt)
				}
				var newValue int
				switch HAProxyFieldAggregationFunctions[attr] {
				case "sum":
					newValue = sumInt(ints)
				case "average":
					newValue = int(averageInt(ints))
				case "max":
					newValue = maxInt(ints)
				}
				err := reflections.SetField(&aggregate.Stats[i], attr, newValue)
				if err != nil {
					panic(err.Error())
				}
			} else {
				newValue, _ := reflections.GetField(multisnapshot.Snapshots[0].Stats[i], attr)
				err := reflections.SetField(&aggregate.Stats[i], attr, newValue)
				if err != nil {
					fmt.Println(aggregate)
					panic(err.Error())
				}
			}
		}
	}
	return aggregate
}
