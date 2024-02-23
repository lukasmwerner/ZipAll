package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"strings"

	"github.com/iriri/minimal/gitignore"
)

var git = flag.Bool("git", false, "follow the gitignore")
var verbose = flag.Bool("v", false, "verbose command output")

func main() {
	flag.Parse()
	archiveName := flag.Arg(0)

	archive_parts := strings.Split(archiveName, ".")
	if "zip" != archive_parts[len(archive_parts)-1] {
		fmt.Println("archive name does not contain 'zip'")
		return
	}
	zipFile, err := os.Create(archiveName)

	if err != nil {
		log.Fatalf("unable to create file %s: %s", archiveName, err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	if *git {
		ignorelist, err := gitignore.FromGit()
		if err != nil {
			log.Panicln(err)
		}

		every := func(path string) error {
			if path == ".gitignore" || path == ".git" || path == zipFile.Name() || path == "." {
				return nil
			}

			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			zipInfo, err := file.Stat()
			if err != nil {
				return err
			}

			header, err := zip.FileInfoHeader(zipInfo)
			if err != nil {
				return err
			}

			header.Name = path
			header.Method = zip.Deflate

			writer, err := zipWriter.CreateHeader(header)
			if err != nil {
				return err
			}
			_, err = io.Copy(writer, file)

			if *verbose {
				log.Printf("adding: %s\n", path)
			}

			return nil
		}

		ignorelist.Walk(".", func(path string, info fs.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			return every(path)
		})

		return
	}

	filesToAdd := flag.Args()[2:]
	for _, filename := range filesToAdd {
		if *verbose {
			log.Printf("adding: %s\n", filename)
		}

		file, err := os.Open(filename)
		if err != nil {
			log.Panicf("unable to read file %s: %s", filename, err)
		}
		defer file.Close()

		zipInfo, err := file.Stat()
		if err != nil {
			log.Panicf("unable to access file info %s: %s", archiveName, err)
		}

		header, err := zip.FileInfoHeader(zipInfo)
		if err != nil {
			log.Panicf("unable to read file header %s: %s", archiveName, err)
		}

		header.Method = zip.Deflate

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			log.Panicf("unable to make zip file header %s: %s", archiveName, err)
		}
		_, err = io.Copy(writer, file)

	}

}
