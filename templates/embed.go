// Package templateassets embeds managed FreeRADIUS template assets.
package templateassets

import "embed"

// Files contains the embedded managed template sets and overlays.
//
//go:embed sets overlays
var Files embed.FS
