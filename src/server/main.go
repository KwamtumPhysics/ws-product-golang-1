package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
	"strconv"
)

type counters struct {
	sync.Mutex
	view  int
	click int
}

var (
	c = counters{}
	contentStack = make(map[string]string)
	content = []string{"sports", "entertainment", "business", "education"}
)

func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to EQ Works 😎")
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	data := content[rand.Intn(len(content))]
	currentTime := time.Now().String()

	c.Lock()
	c.view++
	c.Unlock()

	err := processRequest(r)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(400)
		return
	}

	// simulate random click call
	if rand.Intn(100) < 50 {
		processClick(data)
	}
	key := data + ":" + currentTime
	value := "views: " + strconv.Itoa(c.view) + ", clicks: "+ strconv.Itoa(c.click)
	contentStack[key] = value
}

func processRequest(r *http.Request) error {
	time.Sleep(time.Duration(rand.Int31n(50)) * time.Millisecond)
	return nil
}

func processClick(data string) error {
	c.Lock()
	c.click++
	c.Unlock()

	return nil
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	if !isAllowed() {
		w.WriteHeader(429)
		return
	} else {
		for key, value := range contentStack {
			fmt.Fprintln(w, key," ", value )
		}
	}
}

func isAllowed() bool {
	return true
}

func uploadCounters() error {
	return nil
}

func main() {
	http.HandleFunc("/", welcomeHandler)
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/stats/", statsHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
