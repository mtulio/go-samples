package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type tInputDict map[string]int

const (
	filenameInput  = "./hack/input.txt"
	filenameDict   = "./hack/inputDict.txt"
	filenameOutput = "./hack/outputData.txt"
	bufLenProc     = 8192
	bufLenWriter   = 8192
)

var (
	inputDict   tInputDict
	chProcessor chan string
	chOutFile   chan string
	semControl  chan int
	semCounter  int
)

// checkErr Error controler
func checkErr(e error) {
	if e != nil {
		panic(e)
	}
}

// Initialization
func init() {
	// init dictionary
	inputDict = make(tInputDict)
	f, err := os.Open(filenameDict)
	checkErr(err)
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		inputDict[scanner.Text()] = 1
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

	// misc controler
	chProcessor = make(chan string, bufLenProc)
	chOutFile = make(chan string, bufLenWriter)
	semControl = make(chan int)
	semCounter = 1
}

func outputHandler() {
	f, err := os.Create(filenameOutput)
	if err != nil {
		fmt.Println("Error writing output file")
		return
	}

	defer f.Close()

	for {
		select {
		case line := <-chOutFile:
			f.Sync()
			w := bufio.NewWriter(f)
			nb, err := w.WriteString(fmt.Sprintf("%s\n", line))
			if err != nil {
				fmt.Println(err)
			}
			fmt.Printf("wrote %d bytes\n", nb)
			w.Flush()
			semControl <- -1
		}
	}
}

func lineProcessor(line string) {
	lineNew := ""
	re := regexp.MustCompile("[^A-Za-z]+")
	matchall := re.Split(line, -1)
	submatchall := re.FindAllString(line, -1)
	for i, v := range matchall {
		word := strings.ToLower(v)
		if _, b := inputDict[word]; b {
			if i < len(submatchall) {
				lineNew += fmt.Sprintf("%s", submatchall[i])
			}

		} else if i < len(submatchall) {
			lineNew += fmt.Sprintf("%s%s", v, submatchall[i])

		} else {
			lineNew += fmt.Sprintf("%s", v)
		}
	}
	chOutFile <- lineNew
}

func lineHandler() {

	for {
		select {
		case line := <-chProcessor:
			go lineProcessor(line)
		}
	}
}

func readInput() {

	// read input
	f, err := os.Open(filenameInput)
	checkErr(err)
	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		chProcessor <- scanner.Text()
		semControl <- 1
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
	semControl <- -1
}

func semControler() {

	for {
		select {
		case cn := <-semControl:
			semCounter += cn
		}
	}

}

func main() {

	go lineHandler()
	go outputHandler()
	go semControler()
	readInput()

	// semaphore controler, wait all writers
	for {
		if semCounter == 0 {
			break
		}
	}
}
