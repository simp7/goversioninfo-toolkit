package main

import (
	"flag"
	"io"
	"log"
	"os"

	"github.com/simp7/goversioninfo-toolkit/model"
)

func parseVersionInfoFromFile(fileName string) (model.Info, error) {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return model.Info{}, err
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	data, err := io.ReadAll(file)
	if err != nil {
		return model.Info{}, err
	}

	return model.ParseVersionInfo(data)
}

func overwriteVersionInfoToFile(fileName string, info model.Info) error {
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

	data, err := model.StringifyVersionInfo(model.Info(info))
	if err != nil {
		return err
	}

	_, err = file.Write(data)
	return err
}

func main() {
	notationValue := flag.String("notation", string(model.NotationNormal), "notation for version - simple/normal/detail")
	flag.StringVar(notationValue, "n", *notationValue, "alias for -notation")

	levelValue := flag.String("level", string(model.LevelPatch), "level for versioning - major/minor/patch/build")
	flag.StringVar(levelValue, "l", *levelValue, "alias for -level")

	targetValue := flag.String("target", string(model.TargetBoth), "target for versioning - both/file/product")
	flag.StringVar(targetValue, "t", *targetValue, "alias for -target")

	outputName := flag.String("output", "", "output file name, blank for input itself")
	flag.StringVar(outputName, "o", *outputName, "alias for -output")

	flag.Parse()

	notation := model.VersionNotation(*notationValue)
	level := model.VersionLevel(*levelValue)
	target := model.VersionTarget(*targetValue)

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

	fileVersion, err := info.GetFileVersion()
	if err != nil {
		log.Fatal(err)
	}

	productVersion, err := info.GetProductVersion()
	if err != nil {
		log.Fatal(err)
	}

	fileVersion = fileVersion.Updated(level)
	productVersion = productVersion.Updated(level)

	info = info.VersionUpdated(fileVersion, productVersion, target, notation)

	if err = overwriteVersionInfoToFile(outputFileName, info); err != nil {
		log.Fatal(err)
	}

}
