// Package schema
// Date: 2024/3/6 13:19
// Author: Amu
// Description:
package schema

type Pagination struct {
	Total int64 `json:"total"`
	Size  int64 `json:"size"`
	Page  int64 `json:"page"`
}
