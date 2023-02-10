package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

// test it: go run . Fetch15 http://google.com https://hetpecset.hu https://vra.hu
func Fetch15(urls ...string) {
	fmt.Printf("urls: %T; %v\n", urls, urls)
	fmt.Printf("urls[0]: %T; %v\n", urls[0], urls[0])
	for _, url := range urls {
		fmt.Println("---------------------------------------")
		fmt.Printf("url: %T; %v\n", url, url)
		fmt.Println("---------------------------------------")
		response, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		defer response.Body.Close()

		bodyByteArray, err := io.ReadAll(response.Body)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch: reading %s: %v\n", url, err)
			log.Fatal(err)
		}
		fmt.Printf("%s\n", bodyByteArray)
		fmt.Println("---------------------------------------")
	}
}

// test it: go run . Fetch17 http://google.com https://hetpecset.hu https://vra.hu
func Fetch17(urls ...string) {
	fmt.Printf("urls: %T; %v\n", urls, urls)
	fmt.Printf("urls[0]: %T; %v\n", urls[0], urls[0])
	for _, url := range urls {
		fmt.Println("---------------------------------------")
		fmt.Printf("url: %T; %v\n", url, url)
		fmt.Println("---------------------------------------")
		response, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		defer response.Body.Close()

		copiedChar, err := io.Copy(os.Stdout, response.Body)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch: copying %s: %v\n", url, err)
			log.Fatal(err)
		}
		fmt.Println("\n---------------------------------------")
		log.Printf("Number of copíed chars: %v", copiedChar)
	}
}

// test it: go run . Fetch18 http://google.com https://hetpecset.hu https://vra.hu
func Fetch18(urls ...string) {
	fmt.Printf("urls: %T; %v\n", urls, urls)
	fmt.Printf("urls[0]: %T; %v\n", urls[0], urls[0])
	for _, url := range urls {
		url = strings.TrimLeft(url, " ")
		if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
			url = "http://" + url
		}
		fmt.Println("---------------------------------------")
		fmt.Printf("url: %T; %v\n", url, url)
		fmt.Println("---------------------------------------")
		response, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		defer response.Body.Close()

		copiedChar, err := io.Copy(os.Stdout, response.Body)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch: copying %s: %v\n", url, err)
			log.Fatal(err)
		}
		fmt.Println("\n---------------------------------------")
		log.Printf("Number of copíed chars: %v", copiedChar)
	}
}
