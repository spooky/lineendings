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

func main() {
	searchDir := "./data" //os.Args[1]
	c := make(chan Endings)

	fileList := getFileList(searchDir)
	for _, file := range fileList {
		go getFileEndings(file, c)
	}

	count := Endings{0, 0}
	for _ = range fileList {
		result := <-c

		count.crlf += result.crlf
		count.lf += result.lf
	}

	fmt.Printf("crlf: %d, lf: %d\n", count.crlf, count.lf)
}
