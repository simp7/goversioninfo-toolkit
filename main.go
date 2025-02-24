package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/josephspurrier/goversioninfo"
)

type VersionNotation string

const (
	NotationSimple VersionNotation = "simple"
	NotationNormal VersionNotation = "normal"
	NotationDetail VersionNotation = "detail"
)

type VersionLevel string

const (
	LevelMajor VersionLevel = "major"
	LevelMinor VersionLevel = "minor"
	LevelPatch VersionLevel = "patch"
	LevelBuild VersionLevel = "build"
)

var (
	ErrInvalidStringVersion = errors.New("invalid error format")
)

type VersionTarget string

const (
	TargetBoth    VersionTarget = "both"
	TargetFile    VersionTarget = "file"
	TargetProduct VersionTarget = "product"
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

func getFileVersion(info goversioninfo.VersionInfo) (target goversioninfo.FileVersion, err error) {
	if !isVersionEmpty(info.FixedFileInfo.FileVersion) {
		target = info.FixedFileInfo.FileVersion
	} else {
		target, err = parseVersion(info.StringFileInfo.FileVersion)
	}
	return
}

func getProductVersion(info goversioninfo.VersionInfo) (target goversioninfo.FileVersion, err error) {
	if !isVersionEmpty(info.FixedFileInfo.ProductVersion) {
		target = info.FixedFileInfo.ProductVersion
	} else {
		target, err = parseVersion(info.StringFileInfo.ProductVersion)
	}
	return
}

func versionUp(prev goversioninfo.FileVersion, level VersionLevel) goversioninfo.FileVersion {
	result := prev
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

func main() {
	notationValue := flag.String("notation", string(NotationNormal), "notation for version - simple/normal/detail")
	flag.StringVar(notationValue, "n", *notationValue, "alias for -notation")

	levelValue := flag.String("level", string(LevelPatch), "level for versioning - major/minor/patch/build")
	flag.StringVar(levelValue, "l", *levelValue, "alias for -level")

	targetValue := flag.String("target", string(TargetBoth), "target for versioning - both/file/product")
	flag.StringVar(targetValue, "t", *targetValue, "alias for -target")

	outputName := flag.String("output", "", "output file name, blank for input itself")
	flag.StringVar(outputName, "o", *outputName, "alias for -output")

	flag.Parse()

	notation := VersionNotation(*notationValue)
	level := VersionLevel(*levelValue)
	target := VersionTarget(*targetValue)

	inputFileName := "versioninfo.json"
	args := flag.Args()
	if len(args) >= 1 {
		inputFileName = args[0]
	}
	outputFileName := inputFileName
	if *outputName != "" {
		outputFileName = *outputName
	}

	info, err := parseVersionInfoFromFile(inputFileName)
	if err != nil {
		log.Fatal(err)
	}

	fileVersion, err := getFileVersion(info)
	if err != nil {
		log.Fatal(err)
	}

	productVersion, err := getProductVersion(info)
	if err != nil {
		log.Fatal(err)
	}

	fileVersion = versionUp(fileVersion, level)
	productVersion = versionUp(productVersion, level)

	switch target {
	case TargetFile:
		setFileVersion(&info, fileVersion, notation)
	case TargetProduct:
		setProductVersion(&info, productVersion, notation)
	case TargetBoth:
		setFileVersion(&info, fileVersion, notation)
		setProductVersion(&info, productVersion, notation)
	}

	if err = overwriteVersionInfoToFile(outputFileName, info); err != nil {
		log.Fatal(err)
	}

	return
}
