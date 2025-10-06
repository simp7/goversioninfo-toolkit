package model

import (
	"encoding/json"

	"github.com/josephspurrier/goversioninfo"
)

type Info goversioninfo.VersionInfo

func ParseVersionInfo(data []byte) (version Info, err error) {
	err = json.Unmarshal(data, &version)
	return
}

func StringifyVersionInfo(info Info) ([]byte, error) {
	return json.MarshalIndent(info, "", "\t")
}

func (i Info) GetFileVersion() (target Version, err error) {
	target = Version(i.FixedFileInfo.FileVersion)
	if target.isEmpty() {
		target, err = parseVersion(i.StringFileInfo.FileVersion)
	}
	return
}

func (i Info) GetProductVersion() (target Version, err error) {
	target = Version(i.FixedFileInfo.ProductVersion)
	if target.isEmpty() {
		target, err = parseVersion(i.StringFileInfo.ProductVersion)
	}
	return
}

func (i Info) FileVersionUpdated(version Version, notation VersionNotation) (result Info) {
	result = i

	result.FixedFileInfo.FileVersion = goversioninfo.FileVersion(version)
	result.StringFileInfo.FileVersion = version.String(notation)

	return result
}

func (i Info) ProductVersionUpdated(version Version, notation VersionNotation) (result Info) {
	result = i

	result.FixedFileInfo.ProductVersion = goversioninfo.FileVersion(version)
	result.StringFileInfo.ProductVersion = version.String(notation)

	return result
}

func (i Info) VersionUpdated(fileVersion Version, productVersion Version, target VersionTarget, notation VersionNotation) (result Info) {
	result = i
	switch target {
	case TargetFile:
		result = result.FileVersionUpdated(fileVersion, notation)
	case TargetProduct:
		result = result.ProductVersionUpdated(productVersion, notation)
	case TargetBoth:
		result = result.FileVersionUpdated(fileVersion, notation)
		result = result.ProductVersionUpdated(productVersion, notation)
	}
	return
}
