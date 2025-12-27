# MetroNet Port Scanner

A high-performance, concurrent network port scanner and DNS resolver built in Go with service detection, banner grabbing, and comprehensive network analysis capabilities.

## Features

✅ **Concurrent Scanning** - Uses Go routines with configurable concurrency  
✅ **Banner Grabbing** - Identifies services by capturing and analyzing banners  
✅ **HTTP Body Extraction** - Extracts and displays HTML body content from web services  
✅ **Service Detection** - Recognizes common services (SSH, HTTP, MySQL, etc.)  
✅ **DNS Resolution** - Resolve URLs and hostnames to IPv4 and IPv6 addresses  
✅ **Multiple Scan Modes** - Scan all ports, specific ports, or port ranges  
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
│   ├── scan.go          # Scan command implementation
│   └── resolve.go       # DNS resolution command
├── internal/
│   ├── constants/
│   │   └── constants.go # Default configuration constants
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

# Resolve a domain
./metronet resolve -u google.com
```

## Usage

### Scan Command

The `scan` command performs port scanning on target hosts with various options.

#### Basic Scan (Common Ports)
```bash
./metronet scan -H scanme.nmap.org
```

#### Scan Specific Ports
```bash
./metronet scan -H 192.168.1.1 -p 22,80,443
```

#### Scan Port Range
```bash
./metronet scan -H example.com -p 1-1000
```

#### Scan Multiple Hosts (Subnet)
```bash
./metronet scan -H 192.168.1.0/24 -p 22,80
```

#### Full Port Scan (All 65535 Ports)
```bash
./metronet scan -H scanme.nmap.org --full -c 500
```

#### Advanced Scan Options
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

### Resolve Command

The `resolve` command resolves URLs or hostnames to their IP addresses.

#### Resolve IP from URL
```bash
./metronet resolve -u https://google.com
```

#### Resolve IP from Hostname
```bash
./metronet resolve -u example.com
```

#### Resolve Local Hostname
```bash
./metronet resolve -u localhost
```

## Command-Line Options

### Scan Command Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--host` | `-H` | *required* | Target host or subnet (CIDR) |
| `--ports` | `-p` | all ports | Ports to scan (e.g., 22,80,443 or 1-1000) |
| `--timeout` | `-t` | 2 | Connection timeout in seconds |
| `--concurrency` | `-c` | 100 | Maximum concurrent connections |
| `--randomize` | `-r` | false | Randomize port scanning order |
| `--delay` | `-d` | 0 | Delay between requests in milliseconds |
| `--show-closed` | | false | Show closed and filtered ports |

### Resolve Command Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--url` | `-u` | *required* | URL or hostname to resolve |

## Output Format

### Scan Command Output

```
════════════════════════════════════════════════════════════
    METRONET PORT SCANNER
════════════════════════════════════════════════════════════
Targets:     1 host(s)
Ports:       3 port(s)
Timeout:     2s
Concurrency: 100
Randomize:   false
════════════════════════════════════════════════════════════

╔═══════════════════════════════════════════════════════════════╗
║  Scanning Target: scanme.nmap.org                             ║
╚═══════════════════════════════════════════════════════════════╝

PORT   STATUS   SERVICE         BANNER                BODY
────   ──────   ───────         ──────                ────
22     OPEN ✓   SSH (SSH)       SSH-2.0-OpenSSH_6.6.1 
80     OPEN ✓   HTTP (Apache)   HTTP/1.1 200 OK...    <html>...</html>

────────────────────────────────────────────────────────────
SCAN STATISTICS
────────────────────────────────────────────────────────────
Total Ports Scanned:  3
Open Ports:           2 ✓
Closed Ports:         1
Filtered Ports:       0
Scan Duration:        1.001s
────────────────────────────────────────────────────────────

✓ Found 2 open port(s)
```

### Resolve Command Output

```
════════════════════════════════════════════════════════════
    METRONET DNS RESOLVER
════════════════════════════════════════════════════════════
Target:      google.com
════════════════════════════════════════════════════════════
IPv4 Addresses:
  ✓ 142.250.192.46

IPv6 Addresses:
  ✓ 2607:f8b0:4004:c07::71
  ✓ 2607:f8b0:4004:c07::65
  ✓ 2607:f8b0:4004:c07::8b
  ✓ 2607:f8b0:4004:c07::64

────────────────────────────────────────────────────────────
Total Addresses Found: 5 (1 IPv4, 4 IPv6)
────────────────────────────────────────────────────────────
```

## Technical Implementation

### Concurrency Model
- Uses Go routines for parallel port scanning
- Semaphore pattern to control maximum concurrent connections
- WaitGroup for synchronization
- Efficient worker pool pattern for high-performance scanning

### DNS Resolution
The resolver provides:
1. **Hostname extraction** - Parses URLs to extract hostnames
2. **IPv4 and IPv6 support** - Resolves both IP address types
3. **Clean output** - Categorizes and displays addresses by type
4. **Flexible input** - Accepts URLs, hostnames, or IP addresses

### Service Detection
The scanner identifies services through:
1. **Port-based signatures** - Common services on well-known ports
2. **Banner analysis** - Pattern matching on service banners
3. **Protocol-specific requests** - HTTP requests for web services
4. **HTTP body extraction** - Captures HTML content from web services

### Configuration Management
- **Constants package** - Centralized default values for timeout, concurrency, and delay
- **Command-line flags** - Override defaults on a per-scan basis
- **Flexible configuration** - Easy to adjust for different network environments

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

## Author

Built with ❤️ using Go and Cobra
