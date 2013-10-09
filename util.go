package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

func ConcatFiles(files []string) string {

	var buffer bytes.Buffer
	for _, f := range files {
		// Append the basePath
		file := path.Join(*basePath, f) // TODO - remove the basePath
		content, err := ioutil.ReadFile(file)
		if err != nil {
			fmt.Printf("Failed to read file %s", file)
		}
		buffer.Write(content)
	}

	return buffer.String()
}

func DeleteFiles(files []string) {

	for _, f := range files {
		// Append the basePath
		file := path.Join(*basePath, f) // TODO - remove the basePath
		err := os.Remove(file)
		if err != nil {
			fmt.Printf("Failed to delete file %s", file)
		}
	}

}
