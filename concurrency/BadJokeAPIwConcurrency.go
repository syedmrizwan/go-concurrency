package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

type Response struct {
	ID     string `json:"id"`
	Joke   string `json:"joke"`
	Status int    `json:"status"`
}

var wg = sync.WaitGroup{}

func main() {
	start := time.Now()
	wg.Add(50)

	for i := 0; i < 50; i++ {
		go getBadJoke()
	}
	wg.Wait()
	elapsed := time.Since(start)

	fmt.Printf("Processes took %s", elapsed)
}

func getBadJoke() {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://icanhazdadjoke.com/", nil)
	if err != nil {
		fmt.Print(err.Error())
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)

	if err != nil {
		fmt.Print(err.Error())
	}

	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			fmt.Print(err.Error())
		}
	}(resp.Body)

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Print(err.Error())
	}

	var responseObject Response
	err = json.Unmarshal(bodyBytes, &responseObject)
	if err != nil {
		fmt.Print(err.Error())
	}

	fmt.Println("\n", responseObject.Joke)
	wg.Done()
}
