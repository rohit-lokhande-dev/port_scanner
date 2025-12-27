package scanner

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

// ServiceSignatures maps common port numbers to service names
var ServiceSignatures = map[int]string{
	21:    "FTP",
	22:    "SSH",
	23:    "Telnet",
	25:    "SMTP",
	53:    "DNS",
	80:    "HTTP",
	110:   "POP3",
	143:   "IMAP",
	443:   "HTTPS",
	445:   "SMB",
	3306:  "MySQL",
	3389:  "RDP",
	5432:  "PostgreSQL",
	5900:  "VNC",
	6379:  "Redis",
	8080:  "HTTP-Proxy",
	8443:  "HTTPS-Alt",
	9200:  "Elasticsearch",
	27017: "MongoDB",
}

// GrabBanner attempts to grab a banner from an open port
func GrabBanner(host string, port int, timeout time.Duration) (string, error) {
	address := net.JoinHostPort(host, fmt.Sprintf("%d", port))
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	// Set read deadline
	conn.SetReadDeadline(time.Now().Add(timeout))

	// For some services, we need to send a request first
	if needsRequest(port) {
		if err := sendInitialRequest(conn, port); err != nil {
			return "", err
		}
	}

	// Read banner
	reader := bufio.NewReader(conn)
	banner := make([]byte, 1024)
	n, err := reader.Read(banner)
	if err != nil && n == 0 {
		return "", err
	}
	// fmt.Println("banner: ",string(banner[:n]))
	return string(banner[:n]), nil
}

// IdentifyService attempts to identify the service based on port and banner
func IdentifyService(port int, banner string) string {
	// First check known port signatures
	if service, ok := ServiceSignatures[port]; ok {
		// If we have a banner, try to get more specific info
		if banner != "" {
			detected := detectServiceFromBanner(banner)
			if detected != "" {
				return fmt.Sprintf("%s (%s)", service, detected)
			}
		}
		return service
	}

	// Try to detect from banner if port is unknown
	if banner != "" {
		if detected := detectServiceFromBanner(banner); detected != "" {
			return detected
		}
	}

	return "Unknown"
}

// detectServiceFromBanner analyzes banner content to identify service
func detectServiceFromBanner(banner string) string {
	bannerLower := strings.ToLower(banner)

	signatures := map[string]string{
		"ssh-":       "SSH",
		"http/":      "HTTP",
		"ftp":        "FTP",
		"smtp":       "SMTP",
		"pop3":       "POP3",
		"imap":       "IMAP",
		"mysql":      "MySQL",
		"postgresql": "PostgreSQL",
		"redis":      "Redis",
		"mongodb":    "MongoDB",
		"nginx":      "nginx",
		"apache":     "Apache",
		"microsoft":  "Microsoft",
	}

	for sig, service := range signatures {
		if strings.Contains(bannerLower, sig) {
			return service
		}
	}

	return ""
}

// needsRequest returns true if the service requires an initial request
func needsRequest(port int) bool {
	// HTTP-like services need a request
	return port == 80 || port == 443 || port == 8080 || port == 8443
}

// sendInitialRequest sends an initial request to services that need it
func sendInitialRequest(conn net.Conn, port int) error {
	var request string

	switch port {
	case 80, 8080:
		request = "GET / HTTP/1.0\r\n\r\n"
	case 443, 8443:
		request = "GET / HTTP/1.0\r\n\r\n"
	default:
		return nil
	}

	_, err := conn.Write([]byte(request))
	return err
}
