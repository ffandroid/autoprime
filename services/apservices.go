package services

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/ffandroid/autoprime/cfp/cfrepos"
	"github.com/ffandroid/autoprime/db"
	"github.com/ffandroid/autoprime/repos"
)

func GetMovieVideos(tmdb_id int64) []db.Servers {
	var servers []db.Servers
	var pResponse db.PResponse

	baseURL := fmt.Sprintf(
		"https://primesrc.me/api/v1/s?tmdb=%d&type=movie",
		tmdb_id)
	fmt.Println("baseUrl", baseURL)

	ctn, err := GetContent(baseURL)
	if err != nil {
		fmt.Println(err)
		return servers
	}

	fmt.Println("ctn", ctn)
	json.Unmarshal([]byte(ctn), &pResponse)
	servers = pResponse.Servers

	return servers

	/*var formatted bytes.Buffer
	if err := json.Indent(&formatted, []byte(ctn), "", "  "); err != nil {
		fmt.Printf("Invalid JSON response: %v\n", err)
		return pVideos
	}

	filePath := fmt.Sprintf("/Users/mpro/go/src/autoprime/v_%d_%d_%d.json", tmdb_id, s, e)
	if err := os.WriteFile(filePath, formatted.Bytes(), 0600); err != nil {
		fmt.Println(err)
		return pVideos
	}

	fmt.Println(formatted.String())*/
}

func GetEpisodeVideos(tmdb_id int64, s, e int) []db.Servers {
	var servers []db.Servers
	var pResponse db.PResponse

	baseURL := fmt.Sprintf(
		"https://primesrc.me/api/v1/s?tmdb=%d&season=%d&episode=%d&type=tv",
		tmdb_id, s, e,
	)
	fmt.Println("baseUrl", baseURL)

	ctn, err := GetContent(baseURL)
	if err != nil {
		fmt.Println(err)
		return servers
	}

	fmt.Println("ctn", ctn)
	json.Unmarshal([]byte(ctn), &pResponse)
	servers = pResponse.Servers

	return servers

	/*var formatted bytes.Buffer
	if err := json.Indent(&formatted, []byte(ctn), "", "  "); err != nil {
		fmt.Printf("Invalid JSON response: %v\n", err)
		return pVideos
	}

	filePath := fmt.Sprintf("/Users/mpro/go/src/autoprime/v_%d_%d_%d.json", tmdb_id, s, e)
	if err := os.WriteFile(filePath, formatted.Bytes(), 0600); err != nil {
		fmt.Println(err)
		return pVideos
	}

	fmt.Println(formatted.String())*/
}

func AddVideo(tmdb_id int64, s, e int) {
	var servers []db.Servers
	if s == 0 && e == 0 {

		servers = GetMovieVideos(tmdb_id)
	} else {
		servers = GetEpisodeVideos(tmdb_id, s, e)
	}
	for _, video := range servers {
		var domainID int
		switch video.Name {
		case "Luluvdoo":
			domainID = 10
		case "Vidmoly":
			domainID = 14
		case "Mixdrop":
			domainID = 6
		case "Vidara":
			domainID = 24
		case "Vinovo":
			domainID = 15
		default:
			continue
		}

		resolved := GetFinalUrl(video.Key)
		if videoID := videoIDFromLink(resolved.Link); videoID != "" {
			PostVideo(tmdb_id, s, e, domainID, videoID)
		}
		time.Sleep(5 * time.Second)
	}
}

func videoIDFromLink(link string) string {
	link = strings.TrimSpace(link)
	parsed, err := url.Parse(link)
	if err != nil {
		return ""
	}
	if parsed.Host == "" {
		if parsed.Scheme == "" && parsed.Path == link && !strings.Contains(link, "/") {
			return link
		}
		return ""
	}

	videoID := path.Base(parsed.EscapedPath())
	if videoID == "." || videoID == ".." || videoID == "/" {
		return ""
	}
	videoID, err = url.PathUnescape(videoID)
	if err != nil || strings.Contains(videoID, "/") {
		return ""
	}
	if strings.HasPrefix(videoID, "embed-") && strings.HasSuffix(videoID, ".html") {
		videoID = strings.TrimSuffix(strings.TrimPrefix(videoID, "embed-"), ".html")
	}
	return videoID
}

func PostVideo(tmdb_id int64, s, e, domain_id int, video_id string) {
	const origin = "http://v63662:2516"
	const submitURL = origin + "/submit"

	ctype := "tv"
	if s == 0 && e == 0 {
		ctype = "movie"
	}
	form := url.Values{
		"content_type": {ctype},
		"domain_id":    {strconv.Itoa(domain_id)},
		"vid":          {video_id},
		"tmdb_id":      {strconv.FormatInt(tmdb_id, 10)},
		"season":       {strconv.Itoa(s)},
		"episode":      {strconv.Itoa(e)},
	}
	req, err := http.NewRequest(http.MethodPost, submitURL, strings.NewReader(form.Encode()))
	if err != nil {
		fmt.Printf("PostVideo could not build request: %v\n", err)
		return
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:160.0) Gecko/20100101 Firefox/160.0")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Origin", origin)
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Referer", submitURL)
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Priority", "u=0, i")

	client := &http.Client{
		Timeout: 60 * time.Second,
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("PostVideo request failed: %v\n", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		fmt.Printf("PostVideo received HTTP %s for domain %d\n", resp.Status, domain_id)
	}
}

func GetFinalUrl(key string) db.PVideoInfo {
	pVideoInfo := repos.FindPVideoInfo(key)
	if pVideoInfo.Link != "" {
		fmt.Println("Get", key, "=>", pVideoInfo.Link, " in DB")
		return pVideoInfo
	}
	tailscaleIP := strings.TrimSpace(os.Getenv("TAILSCALE_IP"))
	if net.ParseIP(tailscaleIP) == nil {
		fmt.Println("GetFinalUrl requires TAILSCALE_IP to be a valid IP address")
		return pVideoInfo
	}
	cf := cfrepos.FindCFClearanceByHostname("v62203")
	cookie := cf.Clearance
	ua := cf.UserAgent

	murl := fmt.Sprintf("https://primesrc.me/api/v1/l?key=%s", key)

	req, err := http.NewRequest("GET", murl, nil)
	if err != nil {
		fmt.Println("err", err)
		return pVideoInfo
	}
	tr := &http.Transport{
		Proxy: http.ProxyURL(&url.URL{
			Scheme: "socks5",
			Host:   net.JoinHostPort(tailscaleIP, "1080"),
		}),
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // Disable certificate verification
		},
	}
	//fc, err := os.ReadFile("/Users/mpro/go/src/autoprime/cookie.txt")
	//cookie := string(fc)

	req.Header.Set("Referer", "https://primewire.mov/")
	req.Header.Set("User-Agent", ua)
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Cookie", cookie)

	client := http.Client{
		Transport: tr,
		Timeout:   60 * time.Second, // Set timeout for 10 seconds
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("GetFinalUrl request failed: %v\n", err)
		return pVideoInfo
	}
	defer resp.Body.Close()
	bctn, err := io.ReadAll(resp.Body)
	//fmt.Println(string(bctn))
	if err != nil {
		return pVideoInfo
	}
	json.Unmarshal(bctn, &pVideoInfo)
	if pVideoInfo.Link != "" {
		pVideoInfo.ID = key
	}

	if pVideoInfo.Link != "" {
		repos.SavePVideoInfo(pVideoInfo)
		fmt.Println("Save", key, "=>", pVideoInfo.Link, " in DB")
	}
	return pVideoInfo

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
