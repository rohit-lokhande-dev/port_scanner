package network

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

// ParseHosts parses a host string which can be a single host or CIDR subnet
func ParseHosts(hostStr string) ([]string, error) {
	// Check if it's a CIDR notation
	if strings.Contains(hostStr, "/") {
		return parseCIDR(hostStr)
	}

	// Single host
	return []string{hostStr}, nil
}

// parseCIDR parses a CIDR notation and returns all hosts in the subnet
func parseCIDR(cidr string) ([]string, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, fmt.Errorf("invalid CIDR notation: %v", err)
	}

	var hosts []string
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		hosts = append(hosts, ip.String())
	}

	// Remove network and broadcast addresses for typical subnets
	if len(hosts) > 2 {
		hosts = hosts[1 : len(hosts)-1]
	}

	return hosts, nil
}

// inc increments an IP address
func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

// ParsePortRange parses a port range string (e.g., "1-1024", "80", "22,80,443")
func ParsePortRange(portStr string) ([]int, error) {
	if portStr == "" {
		return nil, nil
	}

	var ports []int

	// Handle comma-separated ports
	parts := strings.Split(portStr, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)

		// Check if it's a range
		if strings.Contains(part, "-") {
			rangeParts := strings.Split(part, "-")
			if len(rangeParts) != 2 {
				return nil, fmt.Errorf("invalid port range: %s", part)
			}

			start, err := strconv.Atoi(strings.TrimSpace(rangeParts[0]))
			if err != nil {
				return nil, fmt.Errorf("invalid start port: %s", rangeParts[0])
			}

			end, err := strconv.Atoi(strings.TrimSpace(rangeParts[1]))
			if err != nil {
				return nil, fmt.Errorf("invalid end port: %s", rangeParts[1])
			}

			if start < 1 || start > 65535 || end < 1 || end > 65535 {
				return nil, fmt.Errorf("port numbers must be between 1 and 65535")
			}

			for i := start; i <= end; i++ {
				ports = append(ports, i)
			}
		} else {
			// Single port
			port, err := strconv.Atoi(part)
			if err != nil {
				return nil, fmt.Errorf("invalid port: %s", part)
			}

			if port < 1 || port > 65535 {
				return nil, fmt.Errorf("port number must be between 1 and 65535")
			}

			ports = append(ports, port)
		}
	}

	return ports, nil
}

