package main

import (
	"io/ioutil"
	"path/filepath"
	"log"
)

func writeFile(data []byte, fileName, outPath string) error {
	fileName = fileName + ".json"
	if err := ioutil.WriteFile(filepath.Join(outPath, fileName), data, 0666); err != nil {
		log.Fatalf("Error during file creation: %v", err)
		return err
	}
	log.Printf("File %v has been created in %v", fileName, outPath)
	return nil
}
