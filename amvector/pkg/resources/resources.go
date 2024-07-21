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
	SystemVersionPath   = filepath.Join(SystemVersionFolder, ".TXT")

	AmvectorFolder              = "amvector"
	AmvectorServiceProfilePath  = filepath.Join(AmvectorFolder, "service_profile.yml")
	AmvectorStorageFolder       = filepath.Join(AmvectorFolder, "storage")
	AmvectorStorageConfigDBPath = filepath.Join(AmvectorStorageFolder, "vector")
	AmvectorSockFolder          = filepath.Join(AmvectorFolder, "socks")
	AmvectorSockFile            = filepath.Join(AmvectorSockFolder, ".sock")

	AmprobeFolder       = "amprobe"
	AmprobeConfigFolder = filepath.Join(AmprobeFolder, "configs")
	AmprobeNginxFolder  = filepath.Join(AmprobeFolder, "nginx")
	AmprobeNginxConfig  = filepath.Join(AmprobeNginxFolder, "nginx.conf")
)
