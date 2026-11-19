package services

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"time"
)

func FindLatestEpisodes() {
	baseUrl := "https://primewire.mov/filter?type=tv&sort=Latest%20Episode"

	ctn, err := GetContent(baseUrl)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(ctn)
}

func GetContent(url string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // Disable certificate verification
		},
	}

	req.Header.Set("Referer", "https://primewire.mov/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:150.0) Gecko/20100101 Firefox/150.0")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	client := http.Client{
		Transport: tr,
		Timeout:   60 * time.Second, // Set timeout for 10 seconds
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	bctn, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(bctn), nil
}
