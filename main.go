package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
)

func main() {
	maxGoroutin := 5
	maxCountURL := 1000
	chURL := make(chan string, maxGoroutin)
	chRes := make(chan map[string]int, maxCountURL)
	toFind := "Go"
	wg := sync.WaitGroup{}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		chURL <- scanner.Text()
		wg.Add(1)

		go Get(chURL, chRes, &wg, toFind)
	}

	wg.Wait()

	ReadResult(chRes)
}

func Get(chUrl chan string, chRes chan map[string]int, wg *sync.WaitGroup, toFind string) {
	url := <- chUrl

	defer wg.Done()

	resp, err := http.Get(url)
	if err != nil {
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	cnt := strings.Count(string(body), toFind)

	chRes <- map[string]int{url: cnt}
}

func ReadResult(chRes chan map[string]int) {
	total := 0
	for {
		select {
		case res := <-chRes:
			for k, v := range res {
				fmt.Printf("Count for %s: %d\n", k, v)
				total += v
			}
		default:
			fmt.Printf("Total: %d\n", total)
			return
		}
	}
}