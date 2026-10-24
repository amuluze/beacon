// Package schema
// Date: 2022/11/9 10:18
// Author: Amu
// Description:
package schema

type GetDockerRegistryMirrorsArgs struct{}

type GetDockerRegistryMirrorsReply struct {
	Mirrors []string `json:"mirrors"`
}

type SetDockerRegistryMirrorsArgs struct {
	Mirrors []string `json:"mirrors"`
}

type SetDockerRegistryMirrorsReply struct{}
