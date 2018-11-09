package HAProxyGoStat

import (
	"github.com/oleiade/reflections"
	"strconv"
	"strings"
)

// Create a new parser.
// Pass in the header line (first line) of the socket output
// parser := CreateHAProxyCSVParser(headers)
// stat = parser(line)
func CreateHAProxyCSVParser(headers string) func(statsLine string) HAProxyStat {
	HeaderMap := strings.Split(strings.TrimSpace(strings.TrimPrefix(headers, "# ")), ",")
	for i, header := range HeaderMap {
		HeaderMap[i] = HAProxyColumnNames[header]
	}
	return func(statsLine string) HAProxyStat {
		statsLineSplit := strings.Split(strings.TrimSpace(statsLine), ",")
		stat := new(HAProxyStat)
		for i, header := range HeaderMap {
			if strings.HasPrefix(statsLineSplit[i], "# ") {
				continue
			}
			statInt, err := strconv.Atoi(statsLineSplit[i])
			if err != nil {
				err = reflections.SetField(stat, header, statsLineSplit[i])
			} else {
				err = reflections.SetField(stat, header, statInt)
			}
			if err != nil {
				continue
			}
		}
		return *stat
	}
}
