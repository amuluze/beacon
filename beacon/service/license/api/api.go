// Package api
// Date: 2023/6/7 12:56
// Author: Amu
// Description:
package api

import "github.com/google/wire"

var Set = wire.NewSet(NewLicenseAPI)
