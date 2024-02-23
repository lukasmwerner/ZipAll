# ZipAll

A CLI tool to add files to zips or just to say specific files to add to the zip

## Installation

```bash
$ go install github.com/lukasmwerner/ZipAll@v0.0.3
```

## Args usage

Usage:

```bash
$ ZipAll Archive.zip file1 file2 file3
```

## Gitignore usage

when there is a gitignore file in the local directory it will use the gitignore
file to ignore those files and then zip everything else Usage:

```bash
$ ZipAll -git Archive.zip
```
