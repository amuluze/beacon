// Package schema
// Date: 2024/6/24 09:37
// Author: Amu
// Description:
package schema

type GetDNSArgs struct {
}

type GetDNSReply struct {
	DNS []string
}

type SetDNSArgs struct {
	DNS []string
}

type SetDNSReply struct {
}
