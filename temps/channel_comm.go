package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	links := []string{
		"http://qq.com",
		"https://amazon.com",
		"http://example.com",
		"https://us.gov",
		"http://onlinekhabar.com",
		"https://www.linkedin.com/feed/",
	}

	ch := make(chan string)
	for _, link := range links {
		go checkLink(link, ch)
		fmt.Println("------------Loop finished------------")
	}

	//fmt.Print receiving value from channel
	//this is blocking operation (waiting on channel is blocking operation)
	for l := range ch {
		go func(link string) {
			time.Sleep(5 * time.Second)
			fmt.Println("--- -- --- --- -From main body goroutine:", link)
			checkLink(link, ch)
		}(l)
	}
	//never try to access same (main)variable from child goroutine. Only share variable through value(pass by value) with function arguments
}

func checkLink(link string, ch chan string) {
	resp, err := http.Get(link)
	fmt.Println(resp.Status)
	if err != nil {
		fmt.Println("Down: ", link, err)
		msg := link
		ch <- msg
		return
	}
	defer resp.Body.Close()
	fmt.Println("Up: ", link)
	msg := link
	ch <- msg
}
