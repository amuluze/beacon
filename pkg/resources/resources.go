// Package resources
// Date: 2024/7/3 09:53
// Author: Amu
// Description:
package resources

import "path/filepath"

var (
	RootPath = "resources"

	SystemEnv           = "/etc/amvector/amvector.env"
	SystemVersionFolder = "/etc/amvector/versions"

	AmvectorFolder              = "amvector"
	AmvectorStorageFolder       = filepath.Join(AmvectorFolder, "storage")
	AmvectorStorageConfigDBPath = filepath.Join(AmvectorStorageFolder, "vector")
	AmvectorSockFolder          = filepath.Join(AmvectorFolder, "socks")
	AmvectorSockFile            = filepath.Join(AmvectorSockFolder, "vector.sock")

	AmprobeFolder       = "amprobe"
	AmprobeConfigFolder = filepath.Join(AmprobeFolder, "configs")
	AmprobeNginxFolder  = filepath.Join(AmprobeFolder, "nginx")
)
