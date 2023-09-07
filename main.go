package main

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	archiveName := os.Args[1]
	filesToAdd := os.Args[2:]

	zipFile, err := os.Create(archiveName)
	if err != nil {
		log.Fatalf("unable to create file %s: %s", archiveName, err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

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
