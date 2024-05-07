package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type Config struct {
	Target string `json:"target"`
}

func main() {
	// Disable TLS certificate verification
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

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
		if targetUrl == "" {
			targetUrl, err = getTargetUrl()
			if err != nil {
				fmt.Printf("%s Failed to retrieve target url from config.json", time.Now().Format(time.RFC3339))
				return
			}
		}

		fmt.Printf("%s Updating DNS for [dkb.crabdance.com]...\n", time.Now().Format(time.RFC3339))

		resp, err := http.Get(targetUrl)
		if err != nil {
			fmt.Printf("%s Error updating DNS: %s\n", time.Now().Format(time.RFC3339), err.Error())
		} else {
			body, err := io.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				fmt.Printf("%s Error reading response body: %s\n", time.Now().Format(time.RFC3339), err.Error())
			} else {
				fmt.Printf("%s Updated. Response: %s\n", time.Now().Format(time.RFC3339), body)
			}
		}

		// Sleep for 600 seconds (10 minutes)
		time.Sleep(600 * time.Second)
	}
}

func getCurrentIP() (string, error) {
	resp, err := http.Get("https://ifconfig.me")
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

func getTargetUrl() (string, error) {
	configFile, err := os.Open("config.json")
	if err != nil {
		fmt.Println("Error opening config file:", err)
		return "", nil
	}
	defer configFile.Close()
	var config Config
	if err := json.NewDecoder(configFile).Decode(&config); err != nil {
		fmt.Println("Error decoding config file:", err)
		return "", nil
	}
	return string(config.Target), nil
}
