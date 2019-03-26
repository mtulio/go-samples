package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"
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

// checkErr control the error exiting when it's triggered.
func checkErr(e error, msg string, panicEnabled bool) {
	if e != nil {
		log.Println("Error: ", msg)
		if panicEnabled {
			panic(e)
		}
	}
}

// Initialize channels and load dictionary
func init() {
	defer timeTrack(timeRunning("init()"))

	inputDict = make(tInputDict)
	f, err := os.Open(filenameDict)
	checkErr(err, "reading dict", true)
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

// outputHandler control the output file, writing line processed.
func outputHandler() {
	defer timeTrack(timeRunning("outputHandler()"))

	f, err := os.Create(filenameOutput)
	checkErr(err, "writing output file", true)
	defer f.Close()

	for {
		select {
		case line := <-chOutFile:
			f.Sync()
			w := bufio.NewWriter(f)
			_, err := w.WriteString(fmt.Sprintf("%s\n", line))
			checkErr(err, "writing line", false)
			// log.Printf("wrote %d bytes\n", nb)
			w.Flush()
			semControl <- -1
		}
	}
}

// lineProcessor process line sending newLine to the output bus.
func lineProcessor(line string) {
	lineNew := ""
	re := regexp.MustCompile("[^A-Za-z0-9]+")
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

// lineHandler handles processor chanel and dispatch to processor function.
func lineHandler() {
	defer timeTrack(timeRunning("lineHandler()"))

	for {
		select {
		case line := <-chProcessor:
			go lineProcessor(line)
		}
	}
}

// readInput reads the input file and send it to processor
func readInput() {
	defer timeTrack(timeRunning("readInput()"))

	f, err := os.Open(filenameInput)
	checkErr(err, "reading input", true)
	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		chProcessor <- scanner.Text()
		semControl <- 1
	}

	if err := scanner.Err(); err != nil {
		// fmt.Fprintln(os.Stderr, "reading standard input:", err)
		checkErr(err, "reading standard input: ", false)
	}
	semControl <- -1
}

// semControler controls the semaphore.
func semControler() {
	defer timeTrack(timeRunning("semControler()"))
	for {
		select {
		case cn := <-semControl:
			semCounter += cn
		}
	}
}

func timeRunning(s string) (string, time.Time) {
	log.Println("Start:", s)
	return s, time.Now()
}

func timeTrack(s string, startTime time.Time) {
	endTime := time.Now()
	log.Println("End:", s, "took", endTime.Sub(startTime))
}

func main() {
	defer timeTrack(timeRunning("Main()"))

	// start all processor handlers
	go semControler()
	go lineHandler()
	go outputHandler()

	// read data set
	readInput()

	// wait to processor threads finishes.
	for {
		if semCounter == 0 {
			break
		}
	}
}
