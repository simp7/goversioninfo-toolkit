# goversioninfo-toolkit

This is toolkit for [github.com/josephspurrier/goversioninfo](https://github.com/josephspurrier/goversioninfo)

## Prerequisites

As this is toolkit for goversioninfo, you need to follow goversioninfo's manual to use bellow tools.

## exevup - tool for versioning up

### Usage

To install, run the following command:

```
go install github.com/simp7/goversioninfo-toolkit/cmd/exevup
```

After satisfying the prerequisites with installing the command, Run following command for versioning

```
exevup {file name} {flags}
```

file name is input file for versioning. default file is versioninfo.json, which is provided by [goversioninfo repository](https://github.com/josephspurrier/goversioninfo/blob/master/testdata/resource/versioninfo.json).

### Command-Line Flags

```
  -level(-l)=[major/minor/patch/build]: level for versioning, default is patch
  -notation(-n)=[simple/normal/detail]: notation for version, default is normal
  -output(-o)={file name}: output file name, default is input file itself
  -target(-t)=[both/file/product]: target for versioning, default is both
```

You also can see description for flags by typing following command
```
exevup --help
```

## Issues

If you notice some problems, please let me know by publishing issues. I will cope with the problem as soon as possible.
