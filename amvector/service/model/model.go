// Package model
// Date: 2024/06/10 18:44:12
// Author: Amu
// Description:
package model

import "github.com/google/wire"

var Set = wire.NewSet(NewModels)

type Models struct{}

func NewModels() *Models {
	return &Models{}
}

func (a *Models) GetAllModels() []interface{} {
	return []interface{}{
		new(Container),
		new(Docker),
		new(Image),
		new(Host),
		new(CPU),
		new(Memory),
		new(Disk),
		new(Net),
	}
}
