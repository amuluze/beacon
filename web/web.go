// Package web
// Date       : 2025/1/16 00:52
// Author     : Amu
// Description:
package web

import (
    "embed"
)

//go:embed dist/* dist/*/*
var FS embed.FS
