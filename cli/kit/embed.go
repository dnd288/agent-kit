package kit

import "embed"

//go:embed all:skills all:templates all:claude all:openspec all:hooks all:ci all:docs
var Content embed.FS
