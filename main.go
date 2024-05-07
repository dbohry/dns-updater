package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

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
		fmt.Printf("%s Updating DNS for [dkb.crabdance.com]...\n", time.Now().Format(time.RFC3339))
		resp, err := http.Get("http://freedns.afraid.org/dynamic/update.php?b24zWGdmOGpMbXZhcmFCOVZHZXBlckM4OjIxNzg0NDM2")
		if err != nil {
			fmt.Printf("%s Error updating DNS: %s\n", time.Now().Format(time.RFC3339), err.Error())
		} else {
			body, err := ioutil.ReadAll(resp.Body)
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
	ip, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(ip), nil
}
