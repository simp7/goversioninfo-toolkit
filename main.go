package main

import (
	"fmt"

	"github.com/josephspurrier/goversioninfo"
)

func main() {
	var version goversioninfo.VersionInfo
	fmt.Println(version.FixedFileInfo.FileVersion.Major)
	return
}
