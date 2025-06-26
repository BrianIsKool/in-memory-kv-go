package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var mu sync.Mutex

type Entry struct {
	Key   int32  `json:"key"`
	Value string `json:"value"`
	TTL   int16  `json:"TTL"`
}

var Cache []Entry

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		fmt.Println("Method not allowed")
		return
	}

	data := new(Entry)
	err := json.NewDecoder(r.Body).Decode(data)

	if err != nil {
		http.Error(w, fmt.Sprintf("Internal server error: %s", err), http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	mu.Lock()
	Cache = append(Cache, *data)
	mu.Unlock()
	fmt.Println(Cache)
}

func get(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	keyStr := r.URL.Query().Get("key")
	keyInt, _ := strconv.Atoi(keyStr)

	var id int32 = int32(keyInt)

	mu.Lock()
	for i, element := range Cache {
		if element.Key == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(Cache[i])
			mu.Unlock()
			return
		}
	}
	mu.Unlock()
	http.Error(w, "Not found", http.StatusNotFound)
	// return
}

func deleteRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	key, err := strconv.Atoi(r.URL.Query().Get("key"))
	if err != nil {
		return
	}
	remove(int32(key))

}

func remove(key int32) {
	mu.Lock()
	for i, element := range Cache {
		if element.Key == key {
			Cache = append(Cache[:i], Cache[i+1:]...)
			mu.Unlock()
			return
		}
	}
	mu.Unlock()
}

func cacheCleaner() {
	for {
		time.Sleep(1 * time.Second)
		mu.Lock()
		var updated []Entry
		for _, element := range Cache {
			element.TTL -= 1
			if element.TTL > 0 {
				updated = append(updated, element)
			}
		}
		Cache = updated
		mu.Unlock()
	}

}

func main() {
	go cacheCleaner()
	http.HandleFunc("/set", handler)
	http.HandleFunc("/get", get)
	http.HandleFunc("/remove", deleteRequest)
	http.ListenAndServe(":8080", nil)
}
