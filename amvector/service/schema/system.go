// Package schema
// Date: 2024/6/24 09:37
// Author: Amu
// Description:
package schema

type RebootArgs struct{}

type RebootReply struct{}

type ShutdownArgs struct{}

type ShutdownReply struct{}

type GetDNSArgs struct{}

type GetDNSReply struct {
	DNS []string
}

type SetDNSArgs struct {
	DNS []string
}

type SetDNSReply struct{}

type GetSystemTimeArgs struct{}

type GetSystemTimeReply struct {
	SystemTime string
}

type SetSystemTimeArgs struct {
	SystemTime string
}

type SetSystemTimeReply struct{}

type GetSystemTimeZoneListArgs struct {
}

type GetSystemTimeZoneListReply struct {
	SystemTimeZoneList []string
}

type GetSystemTimeZoneArgs struct{}

type GetSystemTimeZoneReply struct {
	SystemTimeZone string
}

type SetSystemTimeZoneArgs struct {
	SystemTimeZone string
}

type SetSystemTimeZoneReply struct{}

type GetDockerRegistryMirrorsArgs struct{}

type GetDockerRegistryMirrorsReply struct {
	Mirrors []string
}

type SetDockerRegistryMirrorsArgs struct {
	Mirrors []string
}

type SetDockerRegistryMirrorsReply struct{}
