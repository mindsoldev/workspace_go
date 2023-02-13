package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

// Fetchall does the same fetch of a URL’s contents as the pre vious example,
// but it fetch es many URLs, all con cur rently, so that the pro cess will take no lon g er than the
// longest fetch rat her than the sum of all the fetch times. This version of fetchall discards the
// resp ons es but rep orts the size and elaps ed time for each one

// Test it: go run . fetchAll16 https://google.com https://vra.hu https://hetpecset.hu
func fetchAll16(urls ...string) {
	// fmt.Println("I'm FetchAll")
	start := time.Now()
	chanel := make(chan string)
	for _, url := range urls {
		go fetch16(url, chanel)
	}
	for range urls {
		fmt.Println(<-chanel)
	}
	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}

func fetch16(url string, chanel chan<- string) {
	start := time.Now()
	response, err := http.Get(url)
	if err != nil {
		chanel <- fmt.Sprint(err)
		return
	}
	defer response.Body.Close()

	nbytes, err := io.Copy(io.Discard, response.Body)
	if err != nil {
		chanel <- fmt.Sprintf("while reading %s: %v", url, err)
		return
	}

	secs := time.Since(start).Seconds()
	chanel <- fmt.Sprintf("%.2fs %7d  %s", secs, nbytes, url)
}

// ------------------------------------------------------------------------------------
// Exercis e 1.10: Find a web sit e that pro duces a large amount of dat a. Invest igate caching by
// running fetchall twice in succession to see whether the rep orted time changes much. Do
// you get the same content each time? Modif y fetchall to print its out put to a file so it can be
// examined.

var stdout io.Writer = os.Stdout

const fileBase = "fetch_test_result_"
const fileExt = ".html"

// Test it: go run . fetchAll10 https://google.com https://vra.hu https://hetpecset.hu
func fetchAll10(urls ...string) {
	start := time.Now()
	chanel := make(chan string)
	for i, url := range urls {
		go fetch10(i, url, chanel)
	}
	for range urls {
		fmt.Fprintln(stdout, <-chanel)
	}
	fmt.Fprintf(stdout, "%.2fs elapsed\n", time.Since(start).Seconds())
}

func fetch10(index int, url string, chanel chan<- string) {
	start := time.Now()
	response, err := http.Get(url)
	if err != nil {
		chanel <- fmt.Sprint(err)
		return
	}
	defer response.Body.Close()

	fileName := fileBase + strconv.Itoa(index) + fileExt

	file, err := os.Create(fileName)
	if err != nil {
		chanel <- fmt.Sprintf("Failed to create file. error: %v", err)
	}
	defer file.Close()

	nbytes, err := io.Copy(io.Discard, response.Body)
	if err != nil {
		chanel <- fmt.Sprintf("while reading %s: %v", url, err)
		return
	}

	secs := time.Since(start).Seconds()
	chanel <- fmt.Sprintf("%.2fs %7d  %s", secs, nbytes, url)
}

// ---------------------------------------------------------------------------------------
// Exercis e 1.11: Try fetchall with longer argument lists, such as samples from the top million web sites
// available at alexa.com. How does the program behave if a web site just doesn't respond?

// Test it: go run . fetchAll11 http://hetpecset.hu http://vra.hu http://google.com http://alexa.com
// May be Alexa runs into timeout, others not
func fetchAll11(urls ...string) {
	start := time.Now()
	chanel := make(chan string)
	for _, url := range urls {
		go fetch11(url, chanel)
	}

	for range urls {
		fmt.Println(<-chanel)
	}
	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}

func fetch11(url string, chanel chan<- string) {
	start := time.Now()
	ctx := context.Background()
	// ctx, cancel := context.WithCancel(ctx)
	ctx, cancel := context.WithTimeout(ctx, 700*time.Millisecond)
	defer cancel()

	request, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		chanel <- fmt.Sprint(err)
	}

	fmt.Println("requesting", url)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		chanel <- fmt.Sprint(err)
		return
	}
	defer response.Body.Close()

	nbytes, err := io.Copy(io.Discard, response.Body)
	if err != nil {
		chanel <- fmt.Sprintf("while reading %s: %v", url, err)
		return
	}

	secs := time.Since(start).Seconds()
	chanel <- fmt.Sprintf("%.2fs %7d  %s", secs, nbytes, url)
}
