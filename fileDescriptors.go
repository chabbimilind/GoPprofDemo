package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

var totalSuccess int64
var totalFail int64

func Do(where string, numIter int, numBytes int, wg *sync.WaitGroup) {
	var success int64
	var failure int64
	for i := 0; i < numIter; i++ {
		f, err := os.Open(where)
		if err != nil {
			failure++
		} else {
			success++
		}
		b := make([]byte, numBytes)
		f.Read(b)
		f.Close()
	}
	atomic.AddInt64(&totalSuccess, success)
	atomic.AddInt64(&totalFail, failure)
	wg.Done()
}
func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	var wg sync.WaitGroup
	itr := flag.Int("i", 1000, "Loop trip per go routine.")
	b := flag.Int("n", 100, "File size to read.")
	g := flag.Int("g", 1000, "Level of concurrency while reading files.")
	p := flag.Bool("p", false, "Collect PMU profile.")
	flag.Parse()
	numIter := *itr
	numBytes := *b
	numGoRoutines := *g
	pmu := *p
	fmt.Printf("numIter=%v, numBytes=%v, numGoRoutines=%v, PMU=%v\n", numIter, numBytes, numGoRoutines, pmu)

	for i := 0; i < numGoRoutines; i++ {
		s := strconv.Itoa(i)
		f, err := os.Create("/tmp/xx/" + s)
		check(err)
		d2 := make([]byte, numBytes, numBytes)
		_, err = f.Write(d2)
		check(err)
		f.Sync()
		f.Close()
	}

	file, err := os.Create("file_io.prof")
	if err != nil {
		log.Fatal(err)
	}
	if pmu {
		if err = pprof.StartCPUProfileWithConfig(pprof.CPUCycles(file, 30000000)); err != nil {
			log.Fatal(err)
		}
	}

	start := time.Now()
	for i := 0; i < numGoRoutines; i++ {
		wg.Add(1)
		go Do("/tmp/xx/"+strconv.Itoa(i), numIter, numBytes, &wg)
	}
	wg.Wait()
	elapsed := time.Since(start)
	if pmu {
		pprof.StopCPUProfile()
	}
	file.Close()
	fmt.Printf("FD open:\n Successful=%v (%e/sec)\n Failed=%v (%e/sec)\n %f %% failure rate\n", totalSuccess, float64(totalSuccess)/elapsed.Seconds(), totalFail, float64(totalFail)/elapsed.Seconds(), float64(totalFail*100)/float64(totalFail+totalSuccess))
}
