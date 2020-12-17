package main

import (
	"sync"
	"fmt"
	"sort"
	"strings"
	"strconv"
)

const TH = 6

func ExecutePipeline(jobs ...job)  {
	wg := &sync.WaitGroup{}
	in := make(chan interface {})
	
	for _, job := range jobs {
		wg.Add(1)
		out := make(chan interface {})
		go workerPipeline(wg, job, in, out)
		in = out
	}

	wg.Wait()
}

func workerPipeline(wg *sync.WaitGroup, jobFunc job, in, out chan interface{}) {
	defer wg.Done()
	defer close(out)
	jobFunc(in, out)
}

func SingleHash(in, out chan interface{})  {
	wg := &sync.WaitGroup{}
	
	for val := range in {
		data := fmt.Sprintf("%v", val)
		wg.Add(1)
		crcMd5 := DataSignerMd5(data)
		go workerSingleHash(wg, data, crcMd5, out)
	}

	wg.Wait()
}

func workerSingleHash(wg *sync.WaitGroup, data string, crcMd5 string, out chan interface{}) {
	defer wg.Done()
	
	crc32Chan := make(chan string)
	crcMd5Chan := make(chan string)

	go func ()  {
		crc32Chan <- DataSignerCrc32(data)
	}()

	go func ()  {
		crcMd5Chan <- DataSignerCrc32(crcMd5)
	}()

	var crc32Hash = <- crc32Chan
	var crcMd5Hash = <- crcMd5Chan

	out <- crc32Hash + "~" + crcMd5Hash
}

func MultiHash(in, out chan interface{})  {
	wg := &sync.WaitGroup{}
	
	for val := range in {
		wg.Add(1)
		go workerMultiHash(wg, val, out)
	}

	wg.Wait()
}

func workerMultiHash(wg *sync.WaitGroup, val interface{}, out chan interface{}) {
	defer wg.Done()
	wgMulti := &sync.WaitGroup{}

	hashArray := make([]string, TH)

	for i := 0; i < TH; i++ {
		wgMulti.Add(1)
		data := strconv.Itoa(i) + fmt.Sprintf("%v", val)
		go func (index int) {
			defer wgMulti.Done()

			hashArray[index] = DataSignerCrc32(data)
		}(i)
	}

	wgMulti.Wait()

	multiHash := strings.Join(hashArray, "")
	out <- multiHash
}


func CombineResults(in, out chan interface{}) {
	var hashArray []string

	for i := range in {
		hashArray = append(hashArray, i.(string))
	}

	sort.Strings(hashArray)
	combineResults := strings.Join(hashArray, "_")
	out <- combineResults
}

func main()  {
	
}