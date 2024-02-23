# ZipAdd

A CLI tool to add files to zips or just to say specific files to add to the zip

## Installation

```bash
$ go install github.com/lukasmwerner/ZipAdd@v0.0.2
```

## Args usage

Usage:

```bash
$ zipadd Archive.zip file1 file2 file3
```

## Gitignore usage

when there is a gitignore file in the local directory it will use the gitignore
file to ignore those files and then zip everything else Usage:

```bash
$ zipadd -git Archive.zip
```
