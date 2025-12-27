package cmd

import (
	"fmt"
	"net"
	"net/url"
	"strings"

	"github.com/spf13/cobra"
)

var (
	urlInput string
)

var resolveCmd = &cobra.Command{
	Use:   "resolve",
	Short: "Resolve IP address from URL or hostname",
	Long: `Resolves the IP address(es) for a given URL or hostname.
	
Examples:
  # Resolve IP for a URL
  metronet resolve -u https://google.com
  
  # Resolve IP for a hostname
  metronet resolve -u example.com
  
  # Show all IP addresses (including IPv6)
  metronet resolve -u google.com --all
  
  # Show only IPv6 addresses
  metronet resolve -u google.com --ipv6`,
	RunE: runResolve,
}

func init() {
	rootCmd.AddCommand(resolveCmd)

	// Define flags
	resolveCmd.Flags().StringVarP(&urlInput, "url", "u", "", "URL or hostname to resolve (required)")
	// resolveCmd.Flags().BoolVarP(&resolveIPv6, "ipv6", "6", false, "Show only IPv6 addresses")

	// Mark required flags
	resolveCmd.MarkFlagRequired("url")
}

func runResolve(cmd *cobra.Command, args []string) error {
	// Extract hostname from URL
	hostname, err := extractHostname(urlInput)
	if err != nil {
		return fmt.Errorf("error parsing URL: %v", err)
	}

	// Print resolution header
	fmt.Println("\n════════════════════════════════════════════════════════════")
	fmt.Println("    METRONET DNS RESOLVER")
	fmt.Println("════════════════════════════════════════════════════════════")
	fmt.Printf("Target:      %s\n", hostname)
	fmt.Println("════════════════════════════════════════════════════════════")

	// Resolve IP addresses
	ips, err := net.LookupIP(hostname)
	if err != nil {
		return fmt.Errorf("failed to resolve %s: %v", hostname, err)
	}

	if len(ips) == 0 {
		fmt.Println("⚠️  No IP addresses found for", hostname)
		return nil
	}

	// Filter and display IP addresses
	displayIPAddresses(hostname, ips)

	return nil
}

// extractHostname extracts the hostname from a URL or returns the input if it's already a hostname
func extractHostname(input string) (string, error) {
	// Clean up input
	input = strings.TrimSpace(input)

	// If input doesn't have a scheme, try parsing as hostname
	if !strings.Contains(input, "://") {
		// Check if it's a valid hostname
		if strings.Contains(input, ".") || input == "localhost" {
			return input, nil
		}
		// Try adding http:// scheme
		input = "http://" + input
	}

	// Parse as URL
	parsedURL, err := url.Parse(input)
	if err != nil {
		return "", err
	}

	hostname := parsedURL.Hostname()
	if hostname == "" {
		return "", fmt.Errorf("unable to extract hostname from: %s", input)
	}

	return hostname, nil
}

// displayIPAddresses filters and displays IP addresses based on flags
func displayIPAddresses(hostname string, ips []net.IP) {
	var ipv4Addrs []string
	var ipv6Addrs []string

	// Categorize IPs
	for _, ip := range ips {
		if ip.To4() != nil {
			ipv4Addrs = append(ipv4Addrs, ip.String())
		} else {
			ipv6Addrs = append(ipv6Addrs, ip.String())
		}
	}

	// Display based on flags
	if len(ipv4Addrs) > 0 {
		fmt.Println("IPv4 Addresses:")
		for _, ip := range ipv4Addrs {
			fmt.Printf("  ✓ %s\n", ip)
		}
	}
	if len(ipv6Addrs) > 0 {
		fmt.Println("\nIPv6 Addresses:")
		for _, ip := range ipv6Addrs {
			fmt.Printf("  ✓ %s\n", ip)
		}
	}

	// Print summary
	fmt.Println("\n────────────────────────────────────────────────────────────")
	fmt.Printf("Total Addresses Found: %d (%d IPv4, %d IPv6)\n",
		len(ipv4Addrs)+len(ipv6Addrs), len(ipv4Addrs), len(ipv6Addrs))
	fmt.Println("────────────────────────────────────────────────────────────")
}
