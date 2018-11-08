# HAProxyGoStat

## Install
`go get -u github.com/hmschreck/HAProxyGoStat`

## Use
Create a line parser for your specific version of HAProxy by passing in the header line of the CSV output.

```parser := HAProxyGoStat.CreateHAPRoxyCSVParser(header_line)```

You can then create an HAProxyStat by calling the parser : `start := parser(line)`