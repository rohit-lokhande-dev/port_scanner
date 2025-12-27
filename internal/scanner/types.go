package scanner

import "time"

// PortStatus represents the status of a scanned port
type PortStatus string

const (
	StatusOpen     PortStatus = "OPEN"
	StatusClosed   PortStatus = "CLOSED"
	StatusFiltered PortStatus = "FILTERED"
)

// ScanResult represents the result of scanning a single port
type ScanResult struct {
	Host    string
	Port    int
	Status  PortStatus
	Service string
	Banner  string
	Body    string
}

// ScanConfig holds configuration for the scanner
type ScanConfig struct {
	Host           string
	Ports          []int
	Timeout        time.Duration
	MaxConcurrency int
	RandomizeOrder bool
	DelayBetween   time.Duration
}

// ScanStatistics holds overall scan statistics
type ScanStatistics struct {
	TotalPorts    int
	OpenPorts     int
	ClosedPorts   int
	FilteredPorts int
	ScanDuration  time.Duration
}
