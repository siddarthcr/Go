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

//Job to assign jobs
type Job struct {
	id       int
	randomno int
}

//Result to assign results
type Result struct {
	text string
}

var lines1 = make(chan string, 10)

//var results = make(chan Result, 10)

var total1 float64

func digits(number int) int {
	sum := 0
	no := number
	for no != 0 {
		digit := no % 10
		sum += digit
		no /= 10
	}
	time.Sleep(2 * time.Second)
	return sum
}
func worker(wg *sync.WaitGroup, m *sync.Mutex) {
	for line := range lines {
		lineSplit := strings.Split(line, ",")
		salePrice, err := strconv.ParseFloat(lineSplit[4], 64)
		if err == nil {
			m.Lock()
			total1 += salePrice
			m.Unlock()
		}
	}
	wg.Done()
}
func createWorkerPool(noOfWorkers int, done chan bool) {
	var wg sync.WaitGroup
	var m sync.Mutex
	for i := 0; i < noOfWorkers; i++ {
		wg.Add(1)
		go worker(&wg, &m)
	}
	wg.Wait()
	//close(results)
	fmt.Println("Result:")
	//}
	done <- true
}
func allocate(path string) {

	file, _ := os.Open(path)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines <- scanner.Text()

	}
	close(lines)
	/*for i := 0; i < noOfJobs; i++ {
		randomno := rand.Intn(999)
		job := Job{i, randomno}
		jobs <- job
	}
	close(jobs)
	*/
}

/*func result(done chan bool) {
	//for result := range results {
	fmt.Println("Result:")
	//}
	done <- true
}*/
func mainf() {
	startTime := time.Now()
	//noOfJobs := 100000
	path := "sales/file_1.csv"
	go allocate(path)
	done := make(chan bool)
	//go result(done)
	noOfWorkers := 50000

	go createWorkerPool(noOfWorkers, done)
	<-done
	endTime := time.Now()
	diff := endTime.Sub(startTime)
	fmt.Println("total time taken ", diff.Seconds(), "seconds & Total Value is", total1)
}
