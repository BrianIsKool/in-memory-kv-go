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

// type Answer struct {
// 	Body   string `json:"Body"`
// 	Id     int32  `json:"Id"`
// 	Title  string `json:"Title"`
// 	UserId int32  `json:"UserId"`
// }

// var raw map[string]json.RawMessage

// func main() {
// 	data := []Post{{Body: "just test", Id: 312, Title: "jsjsjad", UserId: 71273}}
// 	answ := make(map[string]Answer)

// 	postdata, err := json.Marshal(data)
// 	if err != nil {
// 		fmt.Println(err)
// 		panic(err)
// 	}

// 	request, err := http.Post("https://jsonplaceholder.typicode.com/posts", "application/json", bytes.NewBuffer(postdata))
// 	if err != nil {
// 		fmt.Println(err)
// 		panic(err)
// 	}
// 	body, _ := io.ReadAll(request.Body)
// 	json.Unmarshal(body, &answ)

// 	fmt.Println(string(body))
// 	// fmt.Println(answ)

// }

// func main() {

// 	data := []Post{}

// 	resp, err := http.Get("https://jsonplaceholder.typicode.com/posts")
// 	if err != nil {
// 		fmt.Println(err)
// 		panic(err)
// 	}

// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		fmt.Println(err)
// 		panic(err)
// 	}
// 	// fmt.Println(string(body))
// 	err = json.Unmarshal(body, &data)
// 	fmt.Println(err)

// 	pretty, _ := json.MarshalIndent(data, "", "  ")
// 	fmt.Println(string(pretty))

// }

// package main

// import "fmt"

// func main() {
// 	test := make(map[string]string)
// 	test["first"] = "test"

// 	if item, status := test["first"]; status {
// 		fmt.Println(item)
// 		// fmt.Println(item)

// 	}
// }
