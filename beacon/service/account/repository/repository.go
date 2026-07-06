// Package repository
// Date       : 2024/9/4 14:55
// Author     : Amu
// Description:
package repository

import "github.com/google/wire"

var Set = wire.NewSet(
	AccountRepositorySet,
)
