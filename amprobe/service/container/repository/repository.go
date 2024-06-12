// Package repository
// Date: 2024/3/6 12:44
// Author: Amu
// Description:
package repository

import (
	"github.com/google/wire"
)

var Set = wire.NewSet(
	ContainerRepoSet,
)
