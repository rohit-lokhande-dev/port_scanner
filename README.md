# MetroNet Port Scanner

A high-performance, concurrent network port scanner built in Go with service detection and banner grabbing capabilities.

## Features

✅ **Concurrent Scanning** - Uses Go routines with configurable concurrency  
✅ **Banner Grabbing** - Identifies services by capturing and analyzing banners  
✅ **Service Detection** - Recognizes common services (SSH, HTTP, MySQL, etc.)  
✅ **Multiple Scan Modes** - Scan common ports, specific ports, port ranges, or all 65535 ports  
✅ **Subnet Support** - Scan multiple hosts using CIDR notation (e.g., 192.168.1.0/24)  
✅ **Port Status Detection** - Distinguishes between open, closed, and filtered ports  
✅ **Randomization** - Randomize port scanning order to avoid pattern detection  
✅ **Rate Limiting** - Add delays between requests to prevent rate limiting  
✅ **Professional Output** - Clean, formatted output with statistics  

## Project Structure

```
metron_code_jam/
├── cmd/
│   ├── root.go          # Root command configuration
│   └── scan.go          # Scan command implementation
├── internal/
│   ├── scanner/
│   │   ├── types.go     # Data structures and types
│   │   ├── scanner.go   # Main scanner orchestrator
│   │   ├── port.go      # Port scanning logic
│   │   └── banner.go    # Banner grabbing & service detection
│   └── network/
│       └── host.go      # Network utilities (CIDR parsing, etc.)
├── main.go              # Application entry point
├── go.mod               # Go module file
└── README.md            # This file
```

## Installation

```bash
# Clone the repository
git clone <repo-url>
cd metron_code_jam

# Build the binary
go build -o metronet

# Run the scanner
./metronet scan -H scanme.nmap.org
```

## Usage

### Basic Scan (Common Ports)
```bash
./metronet scan -H scanme.nmap.org
```

### Scan Specific Ports
```bash
./metronet scan -H 192.168.1.1 -p 22,80,443
```

### Scan Port Range
```bash
./metronet scan -H example.com -p 1-1000
```

### Scan Multiple Hosts (Subnet)
```bash
./metronet scan -H 192.168.1.0/24 -p 22,80
```

### Full Port Scan (All 65535 Ports)
```bash
./metronet scan -H scanme.nmap.org --full -c 500
```

### Advanced Options
```bash
# Randomize port order with high concurrency
./metronet scan -H target.com -p 1-65535 -r -c 1000

# Add delay between requests (rate limiting)
./metronet scan -H target.com -p 1-1000 -d 50

# Show closed and filtered ports
./metronet scan -H target.com -p 1-1000 --show-closed

# Custom timeout
./metronet scan -H target.com -p 1-1000 -t 5
```

## Command-Line Options

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--host` | `-H` | *required* | Target host or subnet (CIDR) |
| `--ports` | `-p` | common ports | Ports to scan (e.g., 22,80,443 or 1-1000) |
| `--timeout` | `-t` | 2 | Connection timeout in seconds |
| `--concurrency` | `-c` | 100 | Maximum concurrent connections |
| `--randomize` | `-r` | false | Randomize port scanning order |
| `--delay` | `-d` | 0 | Delay between requests in milliseconds |
| `--show-closed` | | false | Show closed and filtered ports |
| `--full` | | false | Scan all 65535 ports |

## Output Format

```
════════════════════════════════════════════════════════════
    METRONET PORT SCANNER
════════════════════════════════════════════════════════════
Targets:     1 host(s)
Ports:       3 port(s)
Timeout:     2s
Concurrency: 100
════════════════════════════════════════════════════════════

╔═══════════════════════════════════════════════════════════════╗
║  Scanning Target: scanme.nmap.org                             ║
╚═══════════════════════════════════════════════════════════════╝

PORT   STATUS   SERVICE         BANNER
────   ──────   ───────         ──────
22     OPEN ✓   SSH (SSH)       SSH-2.0-OpenSSH_6.6.1p1
80     OPEN ✓   HTTP (Apache)   HTTP/1.1 200 OK...

────────────────────────────────────────────────────────────
SCAN STATISTICS
────────────────────────────────────────────────────────────
Total Ports Scanned:  3
Open Ports:           2 ✓
Closed Ports:         1
Filtered Ports:       0
Scan Duration:        1.001s
────────────────────────────────────────────────────────────
```

## Technical Implementation

### Concurrency Model
- Uses Go routines for parallel port scanning
- Semaphore pattern to control maximum concurrent connections
- WaitGroup for synchronization

### Service Detection
The scanner identifies services through:
1. **Port-based signatures** - Common services on well-known ports
2. **Banner analysis** - Pattern matching on service banners
3. **Protocol-specific requests** - HTTP requests for web services

### Supported Services
- SSH, HTTP/HTTPS, FTP, SMTP, DNS
- MySQL, PostgreSQL, MongoDB, Redis
- RDP, VNC, SMB, Telnet
- And many more...

## Performance Considerations

- Default concurrency: 100 connections (adjust based on network capacity)
- Default timeout: 2 seconds (increase for slow networks)
- Randomization helps avoid IDS/IPS detection
- Rate limiting with delays prevents network congestion

## Security & Ethics

⚠️ **IMPORTANT**: Only scan networks and systems you own or have explicit permission to test. Unauthorized port scanning may be illegal in your jurisdiction.

This tool is intended for:
- Network administration
- Security auditing (authorized)
- Educational purposes
- Infrastructure assessment (authorized)

## Requirements

- Go 1.25.0 or higher
- Network connectivity
- Appropriate permissions

## License

This project was created for the Metron Code Jam challenge.

## Author

Built with ❤️ using Go and Cobra
