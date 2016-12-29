package main

// This program can be used to output the N largest integers from a file or
// from stdin. It assumes that there is one integer per line. If the line
// cannot be converted to an integer, it will be skipped. The integers are
// printed on a single line in descending order.
//
// Both N and the file can be specified as command-line parameters. If the
// file option is not specified, the program defaults to reading from stdin.
// If the n option is not specified, the program defaults to a small integer.
//
// If any errors occur during execution, the program will exit with exit code 1.
//
// Example file usage:
//		top_n -file=./data -n=15
//
// Example stdin usage:
//		top_n -n=15 < ./data

import (
	"bufio"
	"container/heap"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

const (
	defaultFileName = ""
	defaultN        = 5
)

var errLogger = log.New(os.Stderr, "ERROR: ", log.Ltime)

// DRIVER FUNCTION

// Runs top N program.
func main() {
	fileFlag, nFlag := setupFlags()

	numberScanner, err := setupScanner(fileFlag)
	if err != nil {
		errLogger.Fatalf("Failed to setup scanner - %s", err)
	}

	numberHeap, err := buildHeap(numberScanner, nFlag)
	if err != nil {
		errLogger.Fatalf("Failed to scan numbers - %s", err)
	}

	numbers := takeTopN(numberHeap, nFlag)
	printNumbers(numbers)
}

// PROGRAM SETUP FUNCTIONS

// Sets up flags to be used as command-line options
func setupFlags() (*string, *uint) {
	var fileFlag = flag.String("file", defaultFileName, "file to read")
	var nFlag = flag.Uint("n", defaultN, "amount of numbers to select")
	flag.Parse()

	return fileFlag, nFlag
}

// Sets up logger. Returns a *Logger.
func setupLogger(out io.Writer, prefix string) *log.Logger {
	return log.New(out, prefix, log.Ltime)
}

// Sets up number scanner. If file flag is specified, numbers will
// be scanned from file. Otherwise, numbers will be scanned from stdin.
func setupScanner(fileFlag *string) (*bufio.Scanner, error) {
	if *fileFlag == defaultFileName {
		return bufio.NewScanner(os.Stdin), nil
	}

	dataFile, err := os.Open(*fileFlag)
	if err != nil {
		return nil, err
	}

	return bufio.NewScanner(dataFile), nil
}

// PROGRAM ALGORITHM FUNCTIONS

// Scans numbers with number scanner and builds min-heap. Returns a fully
// constructed min-heap if no error occurred during scan. Otherwise, returns
// partially constructed min-heap and error.
func buildHeap(numberScanner *bufio.Scanner, nFlag *uint) (*TopHeap, error) {
	topHeap := NewTopHeap(*nFlag)
	heap.Init(topHeap)

	for numberScanner.Scan() {
		// Skip lines that can't be converted to ints
		value, err := strconv.Atoi(numberScanner.Text())
		if err != nil {
			continue
		}

		// Fill up the heap until n-elements have been added.
		if topHeap.Len() < int(*nFlag) {
			heap.Push(topHeap, value)
			continue
		}

		// If the value is less than the minimum, we don't need to
		// add it to the heap. We only want the N-highest.
		if value < topHeap.IntHeap[0] {
			continue
		}

		// If we've got a value that's higher than the minimum, make
		// room for the new value by replacing the minimum.
		topHeap.ReplaceMin(value)
	}

	return topHeap, numberScanner.Err()
}

// Selects largest N numbers by popping them off the heap.
func takeTopN(topHeap *TopHeap, nFlag *uint) []int {
	selection := make([]int, 0)

	for i := uint(0); i < *nFlag; i++ {
		selection = append(selection, heap.Pop(topHeap).(int))
	}

	return selection
}

// PROGRAM OUTPUT FUNCTIONS

// Prints numbers in on line, highest first. numbers are expected to be in
// ascending order.
func printNumbers(numbers []int) {
	// Start from the end of numbers in order to print highest first.
	for i := len(numbers) - 1; i >= 0; i-- {
		number := numbers[i]

		// Don't print trailing whitespace if last number in range.
		if i == 0 {
			fmt.Printf("%d", number)
			continue
		}

		fmt.Printf("%d ", number)
	}

	fmt.Println()
}

// SUPPORTING DATA STRUCTURE
// This is a 'no frills' min-heap implementation. Most of this code was taken
// from the min-heap example on http://golang.org. It suited my needs exactly.
// I did add a couple of convenience functions like the constructor function
// (NewTopHeap) and `ReplaceMin`.

type IntHeap []int

type TopHeap struct {
	IntHeap
}

func NewTopHeap(n uint) *TopHeap {
	return &TopHeap{}
}

func (h IntHeap) Len() int {
	return len(h)
}

func (h IntHeap) Less(i, j int) bool {
	return h[i] < h[j]
}

func (h IntHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *IntHeap) Push(x interface{}) {
	*h = append(*h, x.(int))
}

func (h *IntHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func (h *TopHeap) ReplaceMin(value interface{}) {
	heap.Pop(h)
	heap.Push(h, value)
}
