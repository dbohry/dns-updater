package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	// Disable TLS certificate verification
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	var currentIp = ""

	for {
		// Get current IP address
		ip, err := getCurrentIP()
		if err != nil {
			fmt.Printf("%s Error getting current IP address: %s\n", time.Now().Format(time.RFC3339), err.Error())
		} else {
			fmt.Printf("%s Current IP address is: %s\n", time.Now().Format(time.RFC3339), ip)
		}

		// Update DNS
		targetUrl := os.Getenv("TARGET_URL")
		notifyUrl := os.Getenv("NOTIFY_URL")
		frequency := os.Getenv("FREQUENCY")

		if notifyUrl == "" {
			notifyUrl = "https://default.url"
		}

		if frequency == "" {
			frequency = "600"
		}

		if targetUrl == "" {
			fmt.Printf("%s TARGET_URL env not found\n", time.Now().Format(time.RFC3339))
			return
		}

		fmt.Printf("%s Updating DNS...\n", time.Now().Format(time.RFC3339))

		if currentIp == ip {
			fmt.Printf("%s No IP changes...\n", time.Now().Format(time.RFC3339))
		} else {
			resp, err := http.Get(targetUrl)
			if err != nil {
				fmt.Printf("%s Error updating DNS: %s\n", time.Now().Format(time.RFC3339), err.Error())
			} else {
				body, err := io.ReadAll(resp.Body)
				resp.Body.Close()
				if err != nil {
					fmt.Printf("%s Error reading response body: %s\n", time.Now().Format(time.RFC3339), err.Error())
				} else {
					fmt.Printf("%s 200 OK: %s\n", time.Now().Format(time.RFC3339), body)
					currentIp = ip

					notify("Your IP was updated to: "+ip, notifyUrl)
				}
			}
		}

		fmt.Printf("%s Sleeping for: %s seconds\n", time.Now().Format(time.RFC3339), frequency)
		duration, err := strconv.Atoi(frequency)
		time.Sleep(time.Duration(duration) * time.Second)
	}
}

func getCurrentIP() (string, error) {
	ipProvider := os.Getenv("IP_PROVIDER")
	if ipProvider == "" {
		ipProvider = "https://api.ipify.org"
	}
	resp, err := http.Get(ipProvider)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	ip, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(ip), nil
}

func notify(message, url string) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	resp, err := http.Post(url, "application/x-www-form-urlencoded", strings.NewReader(message))
	if err != nil {
		fmt.Println("Failed to send notification:", err)
		return
	}
	defer resp.Body.Close()
}
