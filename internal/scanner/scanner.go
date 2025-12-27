package scanner

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Scanner is the main scanner orchestrator
type Scanner struct {
	config ScanConfig
	mu     sync.Mutex
	stats  ScanStatistics
}

// NewScanner creates a new scanner with the given configuration
func NewScanner(config ScanConfig) *Scanner {
	// Set defaults if not specified
	if config.Timeout == 0 {
		config.Timeout = 2 * time.Second
	}
	if config.MaxConcurrency == 0 {
		config.MaxConcurrency = 100
	}

	return &Scanner{
		config: config,
	}
}

// Scan performs the port scan and returns results
func (s *Scanner) Scan() ([]ScanResult, ScanStatistics, error) {
	startTime := time.Now()

	// Validate host
	if s.config.Host == "" {
		return nil, ScanStatistics{}, fmt.Errorf("host cannot be empty")
	}

	// Prepare ports
	ports := s.config.Ports
	if len(ports) == 0 {
		return nil, ScanStatistics{}, fmt.Errorf("ports cannot be empty")
	}

	// Randomize port order if requested
	if s.config.RandomizeOrder {
		s.customShufflePorts(ports)
		fmt.Println(s.config.Ports)
	}

	// Initialize statistics
	s.stats = ScanStatistics{
		TotalPorts: len(ports),
	}

	// Scan ports concurrently
	results := s.scanConcurrent(ports) // Concurrent Scan here

	// Calculate statistics
	s.stats.ScanDuration = time.Since(startTime)

	return results, s.stats, nil
}

// scanConcurrent scans ports using a worker pool pattern for proper concurrency control
func (s *Scanner) scanConcurrent(ports []int) []ScanResult {
	results := make([]ScanResult, 0, len(ports))
	resultsChan := make(chan ScanResult, len(ports))
	portsChan := make(chan int, len(ports))
	var wg sync.WaitGroup

	// Start worker pool with MaxConcurrency workers
	for i := 0; i < s.config.MaxConcurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			// Each worker processes ports from the channel
			for port := range portsChan {
				// Add delay if configured (between scans, not while idle)
				if s.config.DelayBetween > 0 {
					time.Sleep(s.config.DelayBetween)
				}

				// Scan the port
				result := ScanPort(s.config.Host, port, s.config.Timeout)

				// Update statistics
				s.updateStats(result)

				// Send result
				resultsChan <- result
			}
		}()
	}

	// Send all ports to the work channel
	go func() {
		for _, port := range ports {
			portsChan <- port
		}
		close(portsChan)
	}()

	// Wait for all workers to complete and close results channel
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// Collect results
	for result := range resultsChan {
		results = append(results, result)
	}

	return results
}

// updateStats updates scan statistics thread-safely
func (s *Scanner) updateStats(result ScanResult) {
	s.mu.Lock()
	defer s.mu.Unlock()

	switch result.Status {
	case StatusOpen:
		s.stats.OpenPorts++
	case StatusClosed:
		s.stats.ClosedPorts++
	case StatusFiltered:
		s.stats.FilteredPorts++
	}
}

// shufflePorts randomizes the order of ports
func (s *Scanner) shufflePorts(ports []int) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(ports), func(i, j int) {
		ports[i], ports[j] = ports[j], ports[i]
	})
}

// customShufflePorts randomizes the order of ports using Fisher-Yates algorithm
// This is a completely custom implementation without using any external packages
// It implements its own Linear Congruential Generator (LCG) for random numbers
func (s *Scanner) customShufflePorts(ports []int) {
	fmt.Println("Shuffling ports using custom algorithm")
	n := len(ports)
	if n <= 1 {
		return
	}

	// Initialize seed from current time
	seed := time.Now().UnixNano()

	// Fisher-Yates shuffle algorithm with custom random number generator
	for i := n - 1; i > 0; i-- {
		// Generate a random index between 0 and i (inclusive)
		// Using Linear Congruential Generator (LCG): next = (a * seed + c) mod m
		seed = (seed*1103515245 + 12345) & 0x7fffffff // LCG algorithm
		// fmt.Println("seed: ",seed)
		j := int(seed) % (i + 1)
		// fmt.Println("j",j)
		// fmt.Println("i",i)
		ports[i], ports[j] = ports[j], ports[i]
	}
}
