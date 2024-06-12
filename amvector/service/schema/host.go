// Package schema
// Date: 2024/06/11 00:40:00
// Author: Amu
// Description:
package schema

type HostArgs struct {
}

type CPUArgs struct {
}

type CPUUsageArgs struct {
	StartTime int64
	EndTime   int64
}

type MemoryArgs struct {
}

type MemoryUsageArgs struct {
	StartTime int64
	EndTime   int64
}

type DiskArgs struct {
}

type DiskUsageArgs struct {
	StartTime int64
	EndTime   int64
}

type NetworkUsageArgs struct {
	StartTime int64
	EndTime   int64
}
