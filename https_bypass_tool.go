package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	tlsutls "github.com/refraction-networking/utls"
	"golang.org/x/net/proxy"
)

// Configurable settings
var (
	targetsFile   = "targets.txt"
	proxiesFile   = "proxies.txt"
	numRequests   = 50
	delayMin      = 500  // milliseconds
	delayMax      = 3000 // milliseconds
	tlsProfiles   = []string{"chrome", "firefox", "safari"}
)

type Proxy struct {
	Auth string
	Host string
}

func loadLines(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}
	return lines, scanner.Err()
}

func parseProxy(proxyLine string) Proxy {
	parts := strings.Split(proxyLine, "@")
	return Proxy{
		Auth: parts[0],
		Host: parts[1],
	}
}

func getRandomDelay() time.Duration {
	r := rand.Intn(delayMax-delayMin) + delayMin
	return time.Duration(r) * time.Millisecond
}

func randomTLSProfile() string {
	return tlsProfiles[rand.Intn(len(tlsProfiles))]
}

func makeTLSConfig(profile string, hostname string) (*tlsutls.UConn, error) {
	var spec tlsutls.ClientHelloID
	switch profile {
	case "chrome":
		spec = tlsutls.HelloChrome_Auto
	case "firefox":
		spec = tlsutls.HelloFirefox_Auto
	case "safari":
		spec = tlsutls.HelloIOS_Auto
	default:
		spec = tlsutls.HelloRandomized
	}

	dialConn, err := net.Dial("tcp", hostname+":443")
	if err != nil {
		return nil, err
	}

	uconn := tlsutls.UClient(dialConn, &tls.Config{ServerName: hostname}, spec)
	if err := uconn.Handshake(); err != nil {
		return nil, err
	}

	return uconn, nil
}

func makeRequest(url string, proxyData Proxy, profile string) error {
	proxyAuth := strings.Split(proxyData.Auth, ":")
	dialer, err := proxy.SOCKS5("tcp", proxyData.Host, &proxy.Auth{
		User:     proxyAuth[0],
		Password: proxyAuth[1],
	}, proxy.Direct)
	if err != nil {
		return fmt.Errorf("proxy error: %v", err)
	}

	dialContext := func(network, addr string) (net.Conn, error) {
		return dialer.Dial(network, addr)
	}

	transport := &http.Transport{
		Dial: dialContext,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // optional
		},
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   20 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	// Example browser headers
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/114 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Connection", "keep-alive")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode == 200 && !strings.Contains(string(body), "captcha") {
		fmt.Printf("[✓] %s | TLS: %s | Proxy: %s\n", resp.Status, profile, proxyData.Host)
	} else {
		fmt.Printf("[!] Blocked or Challenge | %d | Proxy: %s\n", resp.StatusCode, proxyData.Host)
	}

	return nil
}

func main() {
	rand.Seed(time.Now().UnixNano())
	targets, err := loadLines(targetsFile)
	if err != nil {
		panic(err)
	}

	proxies, err := loadLines(proxiesFile)
	if err != nil {
		panic(err)
	}

	for i := 0; i < numRequests; i++ {
		target := targets[rand.Intn(len(targets))]
		proxy := parseProxy(proxies[rand.Intn(len(proxies))])
		profile := randomTLSProfile()
		err := makeRequest(target, proxy, profile)
		if err != nil {
			fmt.Println("[-] Error:", err)
		}
		time.Sleep(getRandomDelay())
	}
}
