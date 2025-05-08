// Package repository
// Date:   2024/10/14 16:08
// Author: Amu
// Description:
package repository

import "github.com/google/wire"

var Set = wire.NewSet(
	MailRepositorySet,
)
