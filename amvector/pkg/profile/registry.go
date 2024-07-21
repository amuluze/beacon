// Package profile
// Date: 2024/7/12 14:40
// Author: Amu
// Description:
package profile

type Registry struct {
	Services []*ServiceEntry
}

var RegistryServices = Registry{
	Services: []*ServiceEntry{
		{
			Name:    "network",
			Desc:    "网络配置",
			Builder: NewNetworkBuilder,
		},
		{
			Name:    "amprobe",
			Desc:    "Amprobe",
			Builder: NewAmprobeBuilder,
		},
	},
}
