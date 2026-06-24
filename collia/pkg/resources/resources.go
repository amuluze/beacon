// Package resources
// Date: 2024/7/3 09:53
// Author: Amu
// Description:
package resources

import "path/filepath"

var (
	RootPath = "resources"

	SystemEnv           = "/etc/collia/collia.env"
	SystemVersionFolder = "/etc/collia/versions"

	ColliaFolder              = "collia"
	ColliaStorageFolder       = filepath.Join(ColliaFolder, "storage")
	ColliaStorageConfigDBPath = filepath.Join(ColliaStorageFolder, "collia")
	ColliaSockFolder          = filepath.Join(ColliaFolder, "socks")
	ColliaSockFile            = filepath.Join(ColliaSockFolder, "collia.sock")
)
