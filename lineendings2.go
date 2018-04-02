package main

import (
    "flag"
    "fmt"
    "io/ioutil"
    "os"
    "path/filepath"
    "regexp"
    "sync"
)

var (
    crlf = regexp.MustCompile("\r\n")
    lf   = regexp.MustCompile("\n")
)

type Endings struct {
    crlf, lf uint
}

func (e *Endings) countEndings(fileName string) {
    c, err := ioutil.ReadFile(fileName)
    if err != nil {
        fmt.Printf("fail to read file %s: %v\n", fileName, err)
    }
    x := len(crlf.FindAllIndex(c, -1))
    y := len(lf.FindAllIndex(c, -1))

    e.crlf += uint(x)
    e.lf += uint(y - x)
}

func run(files <-chan string, results chan<- Endings, wg *sync.WaitGroup) {
    e := &Endings{}
    for f := range files {
        e.countEndings(f)
    }
    results <- *e
    wg.Done()
}

func main() {

    searchDir := flag.String("dir", "dir", "directory to search")
    nWorker := flag.Int("n", 2, "number of worker")
    flag.Parse()

    results := make(chan Endings, *nWorker)

    var wg sync.WaitGroup
    wg.Add(*nWorker)
    files := make(chan string, 1000)
    for i := 0; i < *nWorker; i++ {
        go run(files, results, &wg)
    }
    filepath.Walk(*searchDir, func(path string, f os.FileInfo, err error) error {
        if !f.IsDir() {
            files <- path
        }
        return nil
    })
    close(files)
    wg.Wait()
    close(results)

    count := &Endings{}
    for e := range results {
        count.crlf += e.crlf
        count.lf += e.lf
    }
    fmt.Printf("crlf: %d, lf: %d\n", count.crlf, count.lf)
}
