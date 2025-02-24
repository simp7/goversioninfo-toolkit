package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/josephspurrier/goversioninfo"
)

type VersionNotation int

const (
	NotationSimple VersionNotation = iota
	NotationNormal
	NotationDetail
)

var (
	ErrInvalidStringVersion = errors.New("invalid error format")
)

func parseVersionInfo(data []byte) (version goversioninfo.VersionInfo, err error) {
	err = json.Unmarshal(data, &version)
	return
}

func stringifyVersionInfo(info goversioninfo.VersionInfo) ([]byte, error) {
	return json.MarshalIndent(info, "", "\t")
}

func parseVersionInfoFromFile(fileName string) (goversioninfo.VersionInfo, error) {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return goversioninfo.VersionInfo{}, err
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	data, err := io.ReadAll(file)
	if err != nil {
		log.Println(err)
	}
	if err != nil {
		return goversioninfo.VersionInfo{}, err
	}

	return parseVersionInfo(data)
}

func overwriteVersionInfoToFile(fileName string, info goversioninfo.VersionInfo) error {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	err = file.Truncate(0)
	if err != nil {
		log.Println(err)
	}

	data, err := stringifyVersionInfo(info)
	if err != nil {
		return err
	}

	_, err = file.Write(data)
	return err
}

func isVersionEmpty(version goversioninfo.FileVersion) bool {
	return version.Major == 0 && version.Minor == 0 && version.Patch == 0 && version.Build == 0
}

func parseVersion(versionString string) (result goversioninfo.FileVersion, err error) {
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
		fallthrough
	default:
		err = ErrInvalidStringVersion
	}
	return
}

func stringifyVersion(version goversioninfo.FileVersion, notation VersionNotation) string {
	result := fmt.Sprintf("%d.%d", version.Major, version.Minor)
	switch notation {
	case NotationNormal:
		result += fmt.Sprintf(".%d", version.Patch)
	case NotationDetail:
		result += fmt.Sprintf(".%d.%d", version.Patch, version.Build)
	}
	return result
}

func setFileVersion(target *goversioninfo.VersionInfo, version goversioninfo.FileVersion, notation VersionNotation) {
	target.FixedFileInfo.FileVersion = version
	target.StringFileInfo.FileVersion = stringifyVersion(version, notation)
}

func setProductVersion(target *goversioninfo.VersionInfo, version goversioninfo.FileVersion, notation VersionNotation) {
	target.FixedFileInfo.ProductVersion = version
	target.StringFileInfo.ProductVersion = stringifyVersion(version, notation)
}

func main() {
	fileName := "versioninfo.json"
	if len(os.Args) > 1 {
		fileName = os.Args[1]
	}

	info, err := parseVersionInfoFromFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	var target goversioninfo.FileVersion
	if !isVersionEmpty(info.FixedFileInfo.FileVersion) {
		target = info.FixedFileInfo.FileVersion
	} else {
		var err error
		target, err = parseVersion(info.StringFileInfo.FileVersion)
		if err != nil {
			log.Fatal(err)
		}
	}

	setFileVersion(&info, target, NotationDetail)
	setProductVersion(&info, target, NotationDetail)
	if err = overwriteVersionInfoToFile("result.json", info); err != nil {
		log.Fatal(err)
	}

	return
}
