package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	ignore "github.com/sabhiram/go-gitignore"
)

var git = flag.Bool("git", false, "follow the gitignore")

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
		gitignore, err := ignore.CompileIgnoreFile(".gitignore")
		if err != nil {
			log.Panicln(err)
		}

		filepath.Walk(".", func(path string, info fs.FileInfo, err error) error {
			if path == ".git" || path == zipFile.Name() {
				return nil
			}
			if gitignore.MatchesPath(path) {
				log.Printf("skipping: %s\n", path)
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

			return nil
		})

		return
	}

	filesToAdd := flag.Args()[2:]
	for _, filename := range filesToAdd {
		fmt.Println(filename)

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
