// Package repository
// Date: 2024/3/27 16:39
// Author: Amu
// Description:
package repository

import "github.com/google/wire"

var Set = wire.NewSet(
	AuthRepoSet,
)
