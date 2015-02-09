package main

import (
	"fmt"
	"gopkg.in/alecthomas/kingpin.v1"
	"io/ioutil"
	"os"
	"runtime"
	"sort"
	"sync"
	"time"
)

type CompareFile struct {
	file    string
	content string
}

var (
	app          = kingpin.New("simfile", "A command-line tool to search for similarities between files.")
	reffile      = app.Flag("reffile", "Input file.").Required().ExistingFile()
	minsim       = app.Flag("minsim", "Minimal similarity (%).").Default("90.0").Float()
	slicelen     = app.Flag("slicelen", "Slice length of files before comparing.").Default("-1").Int()
	showProgress = app.Flag("progress", "Show progress output").Bool()
	infiles      = app.Arg("--", "One or more files to compare with reffile").Required().Strings()

	wg sync.WaitGroup
)

func main() {
	app.Version("0.0.2")
	app.Parse(os.Args[1:])

	refContent, err := readSliceOfFile(*reffile, *slicelen)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not read input file: %v\n", err)
		return
	}

	var results []string
	lev := CreateLevenshtein(refContent)
	resultChannel := make(chan string, 0)
	workerChannel := make(chan *CompareFile, runtime.NumCPU()*2)
	contentReaderChannel := make(chan string, 1)

	runtime.GOMAXPROCS(runtime.NumCPU() * 3)
	wg.Add(len(*infiles))

	if *showProgress {
		go func() { // Stats printer
			t := time.Second / 3
			for {
				fmt.Printf("\r\033[2K %+v", wg)
				time.Sleep(t)
			}
		}()
	}

	go func() { // Results collector
		for result := range resultChannel {
			results = append(results, result)
		}
	}()

	go func() { // Comparison Worker
		for compareFile := range workerChannel {
			go func(cf *CompareFile) {
				defer wg.Done()
				_, similarity := lev.Distance(cf.content)
				if similarity >= *minsim {
					resultChannel <- fmt.Sprintf("%6.2f%% %s", similarity, cf.file)
				}
			}(compareFile)
		}
	}()

	go func() { // Content Reader
		for file := range contentReaderChannel {
			content, err := readSliceOfFile(file, *slicelen)
			if err != nil {
				wg.Done()
				fmt.Fprintf(os.Stderr, "Failed to read %s: %v\n", file, err)
			}
			workerChannel <- &CompareFile{
				file:    file,
				content: content,
			}
		}
	}()

	// Push all input files to the Content Reader
	for _, file := range *infiles {
		contentReaderChannel <- file
	}

	wg.Wait()
	fmt.Print("\n")

	// Finally, when all work is done, print all results
	sort.Strings(results)
	for _, result := range results {
		fmt.Printf("%s\n", result)
	}
}

func maxInt(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

func readSliceOfFile(filename string, slicelen int) (string, error) {
	contentBytes, err := ioutil.ReadFile(filename)
	if err == nil && slicelen > 0 && len(contentBytes) > slicelen {
		contentBytes = contentBytes[0:slicelen]
	}
	return string(contentBytes), err
}
