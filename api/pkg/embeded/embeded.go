package embeded

import (
	"embed"
)

//go:embed all:ui/dist/**
var UserInterfaceFS embed.FS
