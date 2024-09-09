// Package rpc
// Date: 2022/11/9 10:18
// Author: Amu
// Description:
package schema

type RebootArgs struct{}

type RebootReply struct{}

type ShutdownArgs struct{}

type ShutdownReply struct{}

type GetDNSArgs struct{}

type GetDNSReply struct {
	DNS []string `json:"dns"`
}

type SetDNSArgs struct {
	DNS []string `json:"dns"`
}

type SetDNSReply struct{}

type GetSystemTimeArgs struct{}

type GetSystemTimeReply struct {
	SystemTime string `json:"system_time"`
}

type SetSystemTimeArgs struct {
	SystemTime string `json:"system_time"`
}

type SetSystemTimeReply struct{}

type GetSystemTimeZoneListArgs struct{}

type GetSystemTimeZoneListReply struct {
	SystemTimeZoneList []string `json:"system_time_zone_list"`
}

type GetSystemTimeZoneArgs struct{}

type GetSystemTimeZoneReply struct {
	SystemTimeZone string `json:"system_time_zone"`
}

type SetSystemTimeZoneArgs struct {
	SystemTimeZone string `json:"system_time_zone"`
}

type SetSystemTimeZoneReply struct{}

type GetDockerRegistryMirrorsArgs struct{}

type GetDockerRegistryMirrorsReply struct {
	Mirrors []string `json:"mirrors"`
}

type SetDockerRegistryMirrorsArgs struct {
	Mirrors []string `json:"mirrors"`
}

type SetDockerRegistryMirrorsReply struct{}
