package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

func main() {
	addresses, err := readLines("addresses.txt")
	if err != nil {
		fmt.Println("Read address error:", err)
		return
	}

	proxy, err := readLines("proxy.txt")
	if err != nil {
		fmt.Println("Proxy reading error:", err)
		return
	}

	rand.Seed(time.Now().UnixNano()) // Run random number generator

	userAgents, err := readLines("useragents.txt")
	if err != nil {
		fmt.Println("UserAgents read error:", err)
		return
	}

	var wg sync.WaitGroup
	for _, address := range addresses {
		wg.Add(1)
		go func(address string) {
			defer wg.Done()
			proxy := proxy[rand.Intn(len(proxy))]
			var httpClient *http.Client
			if strings.Contains(proxy, "http://") {
				proxyURL, err := url.Parse(proxy)
				if err != nil {
					fmt.Println("Proxy URL parsing error:", err)
					return
				}
				httpClient = &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)}}
			} else if strings.Contains(proxy, "https://") {
				proxyURL, err := url.Parse(proxy)
				if err != nil {
					fmt.Println("Proxy URL parsing error:", err)
					return
				}
				httpClient = &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)}}
			} else {
				httpClient = &http.Client{}
			}

			reqBody := fmt.Sprintf("address=%s", url.QueryEscape(address))
			req, err := http.NewRequest("POST", "https://api.cascadia.foundation/api/get-faucet", strings.NewReader(reqBody))
			if err != nil {
				fmt.Println("Error in request creation:", err)
				return
			}

			req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
			req.Header.Set("User-Agent", userAgents[rand.Intn(len(userAgents))])

			resp, err := httpClient.Do(req)
			if err != nil {
				fmt.Println("Error while sending the request:", err)
				return
			}

			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Reply Body Read Error:", err)
				return
			}

			fmt.Printf("Response for address %s via proxy %s: %s\n", address, proxy, string(body))
		}(address)
	}
	wg.Wait()
}

func readLines(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
