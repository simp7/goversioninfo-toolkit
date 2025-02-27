package model

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/josephspurrier/goversioninfo"
)

var (
	ErrInvalidStringVersion = errors.New("invalid error format")
)

type Version goversioninfo.FileVersion

type VersionLevel string

const (
	LevelMajor VersionLevel = "major"
	LevelMinor VersionLevel = "minor"
	LevelPatch VersionLevel = "patch"
	LevelBuild VersionLevel = "build"
)

type VersionNotation string

const (
	NotationSimple VersionNotation = "simple"
	NotationNormal VersionNotation = "normal"
	NotationDetail VersionNotation = "detail"
)

type VersionTarget string

const (
	TargetBoth    VersionTarget = "both"
	TargetFile    VersionTarget = "file"
	TargetProduct VersionTarget = "product"
)

func parseVersion(versionString string) (result Version, err error) {
	separated := strings.Split(versionString, ".")
	var element int
	if versionString == "" {
		//비어있는 문자열의 경우 빈 값으로 바로 생성
		return
	}

	switch len(separated) {
	case 4:
		element, err = strconv.Atoi(separated[3])
		if err != nil {
			return
		}
		result.Build = element
		fallthrough
	case 3:
		element, err = strconv.Atoi(separated[2])
		if err != nil {
			return
		}
		result.Patch = element
		fallthrough
	case 2:
		element, err = strconv.Atoi(separated[1])
		if err != nil {
			return
		}
		result.Minor = element
		fallthrough
	case 1:
		element, err = strconv.Atoi(separated[0])
		if err != nil {
			return
		}
		result.Major = element
	default:
		err = ErrInvalidStringVersion
	}
	return
}

func (v Version) isEmpty() bool {
	return v.Major == 0 && v.Minor == 0 && v.Patch == 0 && v.Build == 0
}

func (v Version) Updated(level VersionLevel) Version {
	result := v
	switch level {
	case LevelMajor:
		result.Major += 1
		result.Minor = 0
		result.Patch = 0
		result.Build = 0
	case LevelMinor:
		result.Minor += 1
		result.Patch = 0
		result.Build = 0
	case LevelPatch:
		result.Patch += 1
		result.Build = 0
	case LevelBuild:
		result.Build += 1
	}
	return result
}

func (v Version) String(notation VersionNotation) string {
	result := fmt.Sprintf("%d.%d", v.Major, v.Minor)
	switch notation {
	case NotationNormal:
		result += fmt.Sprintf(".%d", v.Patch)
	case NotationDetail:
		result += fmt.Sprintf(".%d.%d", v.Patch, v.Build)
	}
	return result
}
