package utils

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

func Mutil() string {
	return "Utils V1"
}

func PostContent(jsonData []byte) (*http.Response, error) {
	url := "https://api.eduide.cc/v1/chat/completions"

	// Create a POST request with JSON payload
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Send the request using http.Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Request failed with status: %s\n", resp.Status)
		return nil, err
	}

	return resp, nil
}
func GetContent(url string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		//log.Fatal(err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36")
	//req.Header.Set("X-Requested-With", "XMLHttpRequest")
	//req.Header.Set("Referer", "https://streamtape.com/")
	//req.Header.Set("Host", "streamtape.com")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Accept", "*/*")
	var serverIP net.IP
	transport := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			conn, err := (&net.Dialer{
				Timeout: 30 * time.Second,
			}).DialContext(ctx, network, addr)

			if err != nil {
				return nil, err
			}

			if tcpConn, ok := conn.(*net.TCPConn); ok {
				serverIP = tcpConn.RemoteAddr().(*net.TCPAddr).IP
			}

			return conn, nil
		},
	}
	client := &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	if err != nil {
		return "", err
	}
	bctn, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	fmt.Println("serverIP", serverIP)
	return string(bctn), nil
}

func GetContentWH(url string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		//log.Fatal(err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:124.0) Gecko/20100101 Firefox/124.0")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	//req.Header.Set("Accept-Encoding", "gzip, deflate, br, zstd")
	req.Header.Set("Origin", "https://www.vidbinge.com")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "cross-site")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Priority", "u=4")
	req.Header.Set("TE", "trailers")
	var serverIP net.IP
	transport := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			conn, err := (&net.Dialer{
				Timeout: 30 * time.Second,
			}).DialContext(ctx, network, addr)

			if err != nil {
				return nil, err
			}

			if tcpConn, ok := conn.(*net.TCPConn); ok {
				serverIP = tcpConn.RemoteAddr().(*net.TCPAddr).IP
			}

			return conn, nil
		},
	}
	client := &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	if err != nil {
		return "", err
	}
	bctn, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	fmt.Println("resp.StatusCode", resp.StatusCode)
	fmt.Println("serverIP", serverIP)
	return string(bctn), nil
}

func GetBaseURL(r *http.Request) string {
	// Get the scheme (http or https) from the X-Forwarded-Proto header
	scheme := "http" // Default to http
	if proto := r.Header.Get("X-Forwarded-Proto"); proto != "" {
		scheme = proto
	}

	// Get the host (including port if provided)
	host := r.Host // e.g., "localhost:22100"

	// Combine scheme and host to form the base URL
	baseURL := fmt.Sprintf("%s://%s", scheme, host)

	return baseURL
}
