// Package schema
// Date:   2024/10/15 10:59
// Author: Amu
// Description:
package schema

type CPUAlarmQueryArgs struct {
	StartTime int64 `json:"start_time"`
	EndTime   int64 `json:"end_time"`
}

type CPUAlarmQueryReply struct {
	Data []Usage `json:"data"`
}

type MemoryAlarmQueryArgs struct {
	StartTime int64 `json:"start_time"`
	EndTime   int64 `json:"end_time"`
}

type MemoryAlarmQueryReply struct {
	Data []Usage `json:"data"`
}

type DiskAlarmQueryArgs struct {
	StartTime int64 `json:"start_time"`
	EndTime   int64 `json:"end_time"`
}

type DiskAlarmQueryReply struct {
	Data []Disk `json:"data"`
}

type ServiceAlarmQueryArgs struct {
}

type ServiceAlarmQueryReply struct {
	Data []Container `json:"data"`
}
