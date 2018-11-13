package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var total float64

var lines = make(chan string, 10)

func readFile(path string) {
	file, _ := os.Open(path)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines <- scanner.Text()

	}
	close(lines)
	file.Close()

}

func processLine(wg *sync.WaitGroup, m *sync.Mutex) {
	for line := range lines {
		//fmt.Println(line)
		lineSplit := strings.Split(line, ",")
		salePrice, err := strconv.ParseFloat(lineSplit[4], 64)
		if err != nil {
			fmt.Println("error")
		}
		m.Lock()
		total += salePrice
		fmt.Println("New Total:", total)
		m.Unlock()
		time.Sleep(1 * time.Second)
	}
	wg.Done()
}

func workerPool(noOfWorkers int) {

	var wg sync.WaitGroup
	var m sync.Mutex
	for i := 0; i < noOfWorkers; i++ {
		wg.Add(1)
		go processLine(&wg, &m)
	}

	wg.Wait()
}
func main() {
	startTime := time.Now()

	fmt.Println("Processing sales file")
	path := "sales/file_1.csv"

	go readFile(path)

	noOfWorkers := 5
	fmt.Println("No of workers", noOfWorkers)

	workerPool(noOfWorkers)

	endTime := time.Now()
	diff := endTime.Sub(startTime)

	fmt.Println("Total time taken", diff.Seconds(), "seconds")
	fmt.Println("Total Value is", total)
}
