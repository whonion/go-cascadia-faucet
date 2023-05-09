package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
	"github.com/whonion/go-parse-proxy-geonode"
)

func main() {
	GetFreeProxyList()
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
	httpClient := &http.Client{} // Create a single instance of httpClient

	for i, address := range addresses {
		wg.Add(1)
		go func(address string) {
			defer wg.Done()

			proxyStr := proxy[i%len(proxy)]

			// Parse the proxy URL
			var proxyURL *url.URL
			var err error
			if strings.HasPrefix(proxyStr, "http://") || strings.HasPrefix(proxyStr, "https://") {
				proxyURL, err = url.Parse(proxyStr)
				if err != nil {
					fmt.Println("Proxy URL parsing error:", err)
					return
				}
			} else {
				proxyURL, err = url.Parse("http://" + proxyStr)
				if err != nil {
					fmt.Println("Proxy URL parsing error:", err)
					return
				}
			}

			// Extract the host and port from the URL
			var proxyHost, proxyPort string
			if proxyURL.Scheme == "https" || proxyURL.Scheme == "http" {
				proxyHost, proxyPort, _ = net.SplitHostPort(proxyURL.Host)
			} else {
				proxyHost, proxyPort, _ = net.SplitHostPort(proxyStr)
			}

			// Check if the proxy is working
			dialer := net.Dialer{
				Timeout:   5 * time.Second,
				KeepAlive: 30 * time.Second,
			}
			conn, err := dialer.Dial("tcp", net.JoinHostPort(proxyHost, proxyPort))
			if err != nil {
				fmt.Printf("Proxy %s is not working: %v\n", proxyStr, err)
				return
			}
			defer conn.Close()

			reqBody := fmt.Sprintf("address=%s", url.QueryEscape(address))
			req, err := http.NewRequest("POST", "https://api.cascadia.foundation/api/get-faucet", strings.NewReader(reqBody))
			if err != nil {
				fmt.Println("Error in request creation:", err)
				return
			}

			req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
			req.Header.Set("User-Agent", userAgents[3]) //userAgents[rand.Intn(len(userAgents))]

			resp, err := httpClient.Do(req) // Use the single instance of httpClient
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

			fmt.Printf("Response for address %s via proxy %s: %s\n", address, proxyStr, string(body))
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
func go-parse-proxy-geonode.
