package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
	"time"

	"metron_code_jam/internal/network"
	"metron_code_jam/internal/scanner"
	"metron_code_jam/internal/constants"

	"github.com/spf13/cobra"
)

var (
	host        string
	ports       string
	timeout     int
	concurrency int
	randomize   bool
	delay       int
	showClosed  bool
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Network port scanner",
	Long: `A concurrent network port scanner with service detection and banner grabbing.
	
Examples:
  # Scan common ports on a host
  metronet scan -h scanme.nmap.org
  
  # Scan specific ports
  metronet scan -h 192.168.1.1 -p 22,80,443
  
  # Scan a port range
  metronet scan -h example.com -p 1-1000
  
  # Scan a subnet
  metronet scan -h 192.168.1.0/24 -p 22,80
  
  # Full port scan with high concurrency
  metronet scan -h scanme.nmap.org --full -c 500`,
	RunE: runScan,
}

func init() {
	rootCmd.AddCommand(scanCmd)
	// Define flags
	scanCmd.Flags().StringVarP(&host, "host", "H", "", "Target host or subnet (required)")
	scanCmd.Flags().StringVarP(&ports, "ports", "p", "", "Ports to scan (e.g., 22,80,443 or 1-1000)")
	scanCmd.Flags().IntVarP(&timeout, "timeout", "t", constants.Timeout, "Connection timeout in seconds")
	scanCmd.Flags().IntVarP(&concurrency, "concurrency", "c", constants.Concurrency, "Maximum concurrent connections")
	scanCmd.Flags().BoolVarP(&randomize, "randomize", "r", false, "Randomize port scanning order")
	scanCmd.Flags().IntVarP(&delay, "delay", "d", constants.Delay, "Delay between requests in milliseconds")
	scanCmd.Flags().BoolVar(&showClosed, "show-closed", false, "Show closed and filtered ports")

	// Mark required flags
	scanCmd.MarkFlagRequired("host")
}

func runScan(cmd *cobra.Command, args []string) error {
	// Validate and parse host
	hosts, err := network.ParseHosts(host)
	if err != nil {
		return fmt.Errorf("error parsing host: %v", err)
	}

	// Parse ports
	var portList []int
	if ports != "" {
		portList, err = network.ParsePortRange(ports)
		if err != nil {
			return fmt.Errorf("error parsing ports: %v", err)
		}
	} else {
		portList = scanner.GetAllPorts()
		fmt.Println("⚠️  Full scan mode: scanning all 65535 ports (this may take a while)")
	}

	// Print scan configuration
	printScanHeader(hosts, portList)

	// Scan each host
	for _, targetHost := range hosts {
		if err := scanHost(targetHost, portList); err != nil {
			fmt.Fprintf(os.Stderr, "Error scanning %s: %v\n", targetHost, err)
			continue
		}
	}

	return nil
}

func scanHost(targetHost string, portList []int) error {
	fmt.Printf("\n╔═══════════════════════════════════════════════════════════════╗\n")
	fmt.Printf("║  Scanning Target: %-43s ║\n", targetHost)
	fmt.Printf("╚═══════════════════════════════════════════════════════════════╝\n\n")

	// Configure scanner
	config := scanner.ScanConfig{
		Host:           targetHost,
		Ports:          portList,
		Timeout:        time.Duration(timeout) * time.Second,
		MaxConcurrency: concurrency,
		RandomizeOrder: randomize,
		DelayBetween:   time.Duration(delay) * time.Millisecond,
	}

	// Create and run scanner
	s := scanner.NewScanner(config)
	results, stats, err := s.Scan() // Scan here
	if err != nil {
		return err
	}

	// Sort results by port number
	sort.Slice(results, func(i, j int) bool {
		return results[i].Port < results[j].Port
	})

	// Display results
	displayResults(results, stats)

	return nil
}

func printScanHeader(hosts []string, portList []int) {
	fmt.Println("\n════════════════════════════════════════════════════════════")
	fmt.Println("    METRONET PORT SCANNER")
	fmt.Println("════════════════════════════════════════════════════════════")
	fmt.Printf("Targets:     %d host(s)\n", len(hosts))
	fmt.Printf("Ports:       %d port(s)\n", len(portList))
	fmt.Printf("Timeout:     %ds\n", timeout)
	fmt.Printf("Concurrency: %d\n", concurrency)
	fmt.Printf("Randomize:   %v\n", randomize)
	if delay > 0 {
		fmt.Printf("Delay:       %dms\n", delay)
	}
	fmt.Println("════════════════════════════════════════════════════════════")
}

func displayResults(results []scanner.ScanResult, stats scanner.ScanStatistics) {
	// Create table writer
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)

	fmt.Fprintln(w, "PORT\tSTATUS\tSERVICE\tBANNER\tBODY")
	fmt.Fprintln(w, "────\t──────\t───────\t──────\t──────")

	openCount := 0
	for _, result := range results {
		// Skip closed/filtered ports unless requested
		if !showClosed && result.Status != scanner.StatusOpen {
			continue
		}

		if result.Status == scanner.StatusOpen {
			openCount++
		}

		// Format status with color indicators
		statusStr := formatStatus(result.Status)

		// Truncate banner for display
		bannerStr := strings.TrimSpace(result.Banner)
		bodyStr := strings.TrimSpace(result.Body)
		if len(bannerStr) > 60 {
			bannerStr = bannerStr[:60] + "..."
		}
		bannerStr = strings.ReplaceAll(bannerStr, "\n", " ")
		bannerStr = strings.ReplaceAll(bannerStr, "\r", "")

		fmt.Println("bannerStr: ",bannerStr)
		fmt.Println("body: ",bodyStr)

		fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\n",
			result.Port,
			statusStr,
			result.Service,
			bannerStr,
			bodyStr,
		)
	}

	w.Flush()

	// Print statistics
	fmt.Printf("\n────────────────────────────────────────────────────────────\n")
	fmt.Printf("SCAN STATISTICS\n")
	fmt.Printf("────────────────────────────────────────────────────────────\n")
	fmt.Printf("Total Ports Scanned:  %d\n", stats.TotalPorts)
	fmt.Printf("Open Ports:           %d ✓\n", stats.OpenPorts)
	fmt.Printf("Closed Ports:         %d\n", stats.ClosedPorts)
	fmt.Printf("Filtered Ports:       %d\n", stats.FilteredPorts)
	fmt.Printf("Scan Duration:        %v\n", stats.ScanDuration.Round(time.Millisecond))
	fmt.Printf("────────────────────────────────────────────────────────────\n\n")

	if openCount == 0 {
		fmt.Println("⚠️  No open ports found")
	} else {
		fmt.Printf("✓ Found %d open port(s)\n", openCount)
	}
}

func formatStatus(status scanner.PortStatus) string {
	switch status {
	case scanner.StatusOpen:
		return "OPEN ✓"
	case scanner.StatusClosed:
		return "CLOSED"
	case scanner.StatusFiltered:
		return "FILTERED"
	default:
		return string(status)
	}
}
