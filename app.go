// main.go
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"sync"
	"time"
)

type RequestPayload struct {
	ToSort [][]int `json:"to_sort"`
}

type ResponsePayload struct {
	SortedArrays [][]int `json:"sorted_arrays"`
	TimeNs       int64   `json:"time_ns"`
}

func sortSequential(arrays [][]int) [][]int {
	sortedArrays := make([][]int, len(arrays))
	for i, subArray := range arrays {
		sortedSubArray := make([]int, len(subArray))
		copy(sortedSubArray, subArray)
		sort.Ints(sortedSubArray)
		sortedArrays[i] = sortedSubArray
	}
	return sortedArrays
}

func sortConcurrent(arrays [][]int) [][]int {
	var wg sync.WaitGroup
	var mu sync.Mutex
	sortedArrays := make([][]int, len(arrays))

	for i, subArray := range arrays {
		wg.Add(1)
		go func(i int, subArray []int) {
			defer wg.Done()

			sortedSubArray := make([]int, len(subArray))
			copy(sortedSubArray, subArray)
			sort.Ints(sortedSubArray)

			mu.Lock()
			sortedArrays[i] = sortedSubArray
			mu.Unlock()
		}(i, subArray)
	}

	wg.Wait()
	return sortedArrays
}

func processSingleHandler(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload
	err := json.NewDecoder(r.Body).Decode(&requestPayload)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	startTime := time.Now()
	sortedArrays := sortSequential(requestPayload.ToSort)
	endTime := time.Now()

	responsePayload := ResponsePayload{
		SortedArrays: sortedArrays,
		TimeNs:       endTime.Sub(startTime).Nanoseconds(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responsePayload)
}

func processConcurrentHandler(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload
	err := json.NewDecoder(r.Body).Decode(&requestPayload)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	startTime := time.Now()
	sortedArrays := sortConcurrent(requestPayload.ToSort)
	endTime := time.Now()

	responsePayload := ResponsePayload{
		SortedArrays: sortedArrays,
		TimeNs:       endTime.Sub(startTime).Nanoseconds(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responsePayload)
}

func main() {
	http.HandleFunc("/process-single", processSingleHandler)
	http.HandleFunc("/process-concurrent", processConcurrentHandler)

	fmt.Println("Server is running on :8000")
	http.ListenAndServe(":8000", nil)
}
