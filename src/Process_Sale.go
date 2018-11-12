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
		fmt.Println(line)
		lineSplit := strings.Split(line, ",")
		salePrice, _ := strconv.ParseFloat(lineSplit[4], 64)
		m.Lock()
		total += salePrice
		fmt.Println(total)
		m.Unlock()
		time.Sleep(2 * time.Second)
	}
	wg.Done()
}

func workerPool(noOfWorkers int, done chan bool) {
	var wg sync.WaitGroup

	var m sync.Mutex
	for i := 0; i < noOfWorkers; i++ {
		wg.Add(1)
		go processLine(&wg, &m)
	}
	wg.Wait()
	done <- true
}
func main() {

	fmt.Println("Processing sales file")
	path := "sales/file_1.csv"

	go readFile(path)

	noOfWorkers := 1
	done := make(chan bool)
	workerPool(noOfWorkers, done)
	<-done
	fmt.Println("Total Sales:", total)
}
