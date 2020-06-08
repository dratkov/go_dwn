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

type CountGo struct {
	URL     string
	CountGo int
}

func main() {
	maxGoroutin := 5
	maxCountURL := 1000
	chURL := make(chan string, maxGoroutin)
	chRes := make(chan CountGo, maxCountURL)
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

func Get(chUrl chan string, chRes chan CountGo, wg *sync.WaitGroup, toFind string) {
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

	chRes <- CountGo{URL: url, CountGo: cnt}
}

func ReadResult(chRes chan CountGo) {
	total := 0
	for {
		select {
		case res := <-chRes:
			fmt.Printf("Count for %s: %d\n", res.URL, res.CountGo)
			total += res.CountGo
		default:
			fmt.Printf("Total: %d\n", total)
			return
		}
	}
}