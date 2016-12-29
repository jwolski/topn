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
//		go build -o topn
//		topn -file=./data -n=15
//
// Example stdin usage:
//		go build -o topn
//		topn -n=15 < ./data
//
// Example contents of data file:
//		24557
//		13225
//		27592
//		23095
//		17253

import (
	"bufio"
	"container/heap"
	"flag"
	"fmt"
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

	// Setup scanner - if file flag has not been provided
	// read from stdin.
	var numberScanner *bufio.Scanner

	if *fileFlag == defaultFileName {
		numberScanner = bufio.NewScanner(os.Stdin)
	} else {
		dataFile, err := os.Open(*fileFlag)
		if err != nil {
			errLogger.Fatalf("Failed to open file - %s", err)
		}

		defer dataFile.Close()
		numberScanner = bufio.NewScanner(dataFile)
	}

	// Build min-heap from scanning list of numbers
	numberHeap, err := buildHeap(numberScanner, nFlag)
	if err != nil {
		errLogger.Fatalf("Failed to scan numbers - %s", err)
	}

	// Take the top N integers from the heap and print
	// them in descending order.
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
	if topHeap.Len() == 0 {
		return []int{}
	}

	selection := make([]int, 0)

	for i := uint(0); i < *nFlag && topHeap.Len() > 0; i++ {
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

// Creates new TopHeap instance.
func NewTopHeap(n uint) *TopHeap {
	return &TopHeap{}
}

// Returns number of elements in heap.
func (h IntHeap) Len() int {
	return len(h)
}

// Compares heap elements.
func (h IntHeap) Less(i, j int) bool {
	return h[i] < h[j]
}

// Swaps elements within the heap.
func (h IntHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

// Adds an element to the heap.
func (h *IntHeap) Push(x interface{}) {
	*h = append(*h, x.(int))
}

// Removes and returns the minimum element from the heap.
func (h *IntHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// Replaces the minimum element in the heap with the provided value
func (h *TopHeap) ReplaceMin(value interface{}) {
	heap.Pop(h)
	heap.Push(h, value)
}
