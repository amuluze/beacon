// Package service
// Date: 2022/11/9 10:18
// Author: Amu
// Description:
package service

import "github.com/google/wire"

var Set = wire.NewSet(
	AuditServiceSet,
)
