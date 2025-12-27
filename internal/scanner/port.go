package scanner

import (
	"fmt"
	"net"
	"time"
)

// ScanPort scans a single port and returns the result
func ScanPort(host string, port int, timeout time.Duration) ScanResult {
	result := ScanResult{
		Host:   host,
		Port:   port,
		Status: StatusClosed,
	}

	address := net.JoinHostPort(host, fmt.Sprintf("%d", port))
	fmt.Printf("Start scanning port %d\n", port)
	conn, err := net.DialTimeout("tcp", address, timeout)

	if err != nil {
		// Determine if port is filtered or closed
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			result.Status = StatusFiltered
		} else {
			result.Status = StatusClosed
		}
		fmt.Printf("End scanning port %d\n with error: %v", port, err)
		return result
	}
	defer conn.Close()

	// Port is open
	result.Status = StatusOpen

	// Attempt banner grabbing
	banner, err := GrabBanner(host, port, timeout)
	// fmt.Printf("banner %T: ",banner)
	if err == nil && banner != "" {
		result.Banner = cleanBanner(banner)
		result.Service = IdentifyService(port, banner)
		result.Body = GetBody(banner)
	} else {
		result.Service = IdentifyService(port, "")
	}
	fmt.Printf("End scanning port %d\n", port)
	return result
}

// cleanBanner removes non-printable characters and trims the banner
func cleanBanner(banner string) string {
	// Remove null bytes and other control characters
	cleaned := make([]rune, 0, len(banner))
	for _, r := range banner {
		if r >= 32 && r <= 126 || r == '\n' || r == '\r' || r == '\t' {
			cleaned = append(cleaned, r)
		}
	}

	result := string(cleaned)

	// Truncate if too long
	if len(result) > 200 {
		result = result[:200] + "..."
	}

	return result
}

func GetBody(cleanedBanner string) string {
	// Find <body> tag (case insensitive)
	bodyStart := -1
	bodyEnd := -1

	// Convert to lowercase for searching
	lowerBanner := ""
	for _, r := range cleanedBanner {
		if r >= 'A' && r <= 'Z' {
			lowerBanner += string(r + 32)
		} else {
			lowerBanner += string(r)
		}
	}

	// Find opening <body> tag
	for i := 0; i < len(lowerBanner)-5; i++ {
		if lowerBanner[i:i+5] == "<body" {
			// Find the end of the opening tag
			for j := i; j < len(lowerBanner); j++ {
				if lowerBanner[j] == '>' {
					bodyStart = j + 1
					break
				}
			}
			break
		}
	}

	// Find closing </body> tag
	for i := 0; i < len(lowerBanner)-7; i++ {
		if lowerBanner[i:i+7] == "</body>" {
			bodyEnd = i
			break
		}
	}

	// If both tags found, extract body content
	if bodyStart != -1 && bodyEnd != -1 && bodyStart < bodyEnd {
		return cleanedBanner[bodyStart:bodyEnd]
	}

	// Return empty string if no body tags found
	return ""
}

// GetAllPorts returns a range of ports (1-65535)
func GetAllPorts() []int {
	ports := make([]int, 65535)
	for i := 0; i < 65535; i++ {
		ports[i] = i + 1
	}
	return ports
}

// GetPortRange returns a range of ports
func GetPortRange(start, end int) []int {
	if start < 1 {
		start = 1
	}
	if end > 65535 {
		end = 65535
	}
	if start > end {
		start, end = end, start
	}

	ports := make([]int, end-start+1)
	for i := range ports {
		ports[i] = start + i
	}
	return ports
}
