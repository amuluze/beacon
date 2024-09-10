// Package repository
// Date: 2024/06/11 19:38:14
// Author: Amu
// Description:
package repository

import (
	"github.com/google/wire"
)

var Set = wire.NewSet(
	ContainerServiceSet,
)
