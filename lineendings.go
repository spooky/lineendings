package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

type Endings struct {
	crlf, lf uint
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func getFileList(searchDir string) []string {
	fileList := []string{}
	filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() {
			fileList = append(fileList, path)
		}
		return nil
	})

	return fileList
}

func getFileEndings(filename string, c chan Endings) {
	fileContent, err := ioutil.ReadFile(filename)
	check(err)
	c <- countEndings(string(fileContent))
}

func countEndings(s string) Endings {
	crlf := regexp.MustCompile("\r\n")
	lf := regexp.MustCompile("\n")

	x := len(crlf.FindAllStringIndex(s, -1))
	y := len(lf.FindAllStringIndex(s, -1))

	return Endings{uint(x), uint(y - x)}
}

func splitIntoChunks(arr []string, chunkSize int) [][]string {
	if chunkSize <= 0 {
		panic("chunk size too small")
	}

	if len(arr) <= chunkSize {
		return [][]string{arr}
	}

	numChunks := int(len(arr)/chunkSize) + 1
	chunks := make([][]string, numChunks)

	for i := 0; i < numChunks; i++ {
		l := i * chunkSize
		u := l + chunkSize

		if u > len(arr) {
			u = len(arr)
		}

		chunks[i] = arr[l:u]
	}

	return chunks
}

func main() {
	searchDir := os.Args[1]
	c := make(chan Endings)
	chunkSize := 1000

	fileList := getFileList(searchDir)

	count := Endings{0, 0}
	for _, chunk := range splitIntoChunks(fileList, chunkSize) {
		for _, file := range chunk {
			go getFileEndings(file, c)
		}

		for _ = range chunk {
			result := <-c

			count.crlf += result.crlf
			count.lf += result.lf
		}
	}

	fmt.Printf("crlf: %d, lf: %d\n", count.crlf, count.lf)
}
