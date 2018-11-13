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

func processLine(wg *sync.WaitGroup, m *sync.Mutex, ledger map[string]float64) {
	for line := range lines {
		//fmt.Println(line)
		lineSplit := strings.Split(line, ",")
		salePrice, err := strconv.ParseFloat(lineSplit[4], 64)
		if err != nil {
			fmt.Println("error")
		}
		dateString := lineSplit[3]
		dateString = dateString[0:10]
		fmt.Println("Reading invoice id:", lineSplit[0])
		m.Lock()
		currentValue, _ := ledger[dateString]
		ledger[dateString] = currentValue + salePrice
		m.Unlock()
		time.Sleep(1 * time.Second)
	}
	wg.Done()
}

func workerPool(noOfWorkers int, ledger map[string]float64) {

	var wg sync.WaitGroup
	var m sync.Mutex

	for i := 0; i < noOfWorkers; i++ {
		wg.Add(1)
		go processLine(&wg, &m, ledger)
	}

	wg.Wait()
}
func main() {
	startTime := time.Now()
	fmt.Println("\nProcessing sales file")
	path := "sales/Chennai.csv"

	go readFile(path)

	ledger := make(map[string]float64)
	noOfWorkers := 6

	fmt.Println()
	fmt.Println("NUM OF WORKERS", noOfWorkers)

	fmt.Println()
	workerPool(noOfWorkers, ledger)

	endTime := time.Now()
	diff := endTime.Sub(startTime)

	fmt.Println("\nTotal time taken", diff.Seconds(), "seconds")
	fmt.Println("\nTotal Sales value by date:")
	for entryDate, entryValue := range ledger {
		fmt.Printf("%s Rs.%.2f\n", entryDate, entryValue)
	}
	fmt.Println()
}
