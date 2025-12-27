# Metron Port Scanner - Quick Reference

## Quick Start
```bash
# Build
go build -o metron

# Basic scan
./metron scan -H scanme.nmap.org

# Scan specific ports
./metron scan -H target.com -p 22,80,443,8080
```

## Common Use Cases

### 1. Quick Web Server Check
```bash
./metron scan -H example.com -p 80,443,8080,8443
```

### 2. Database Server Audit
```bash
./metron scan -H db.example.com -p 3306,5432,27017,6379
```

### 3. Internal Network Discovery
```bash
./metron scan -H 192.168.1.0/24 -p 22,80,443,3389
```

### 4. Comprehensive Port Scan
```bash
./metron scan -H target.com -p 1-65535 -c 1000
```

### 5. Stealth Scan (Randomized + Delayed)
```bash
./metron scan -H target.com -p 1-1000 -r -d 100 -c 10
```

## Port Categories

### Web Services
- HTTP: 80, 8000, 8080
- HTTPS: 443, 8443

### Remote Access
- SSH: 22
- Telnet: 23
- RDP: 3389
- VNC: 5900

### Databases
- MySQL: 3306
- PostgreSQL: 5432
- MongoDB: 27017
- Redis: 6379

### Email
- SMTP: 25, 587
- POP3: 110, 995
- IMAP: 143, 993

## Performance Tuning

### Fast Local Network
```bash
./metron scan -H 192.168.1.1 -p 1-65535 -c 2000 -t 1
```

### Slow/Remote Network
```bash
./metron scan -H remote-host.com -p 1-1000 -c 50 -t 5
```

### Avoid Rate Limiting
```bash
./metron scan -H target.com -p 1-10000 -c 10 -d 50 -r
```

## Output Modes

### Show Only Open Ports (Default)
```bash
./metron scan -H target.com -p 1-1000
```

### Show All Ports (Open + Closed + Filtered)
```bash
./metron scan -H target.com -p 1-1000 --show-closed
```

## Architecture Overview

```
User Command
     │
     ▼
CLI Layer (cmd/scan.go)
     │
     ├─► Parse Arguments
     ├─► Validate Input
     └─► Configure Scanner
          │
          ▼
Scanner Layer (internal/scanner/)
     │
     ├─► Create Worker Pool
     ├─► Distribute Ports
     └─► Collect Results
          │
          ├─► Port Scanner (port.go)
          │    └─► TCP Connect
          │
          └─► Banner Grabber (banner.go)
               ├─► Capture Banner
               └─► Identify Service
```

## Tips & Tricks

1. **Speed vs. Accuracy**: Higher concurrency = faster but may miss slower services
2. **Timeout Tuning**: Increase timeout for slow networks or filtered firewalls
3. **Randomization**: Use `-r` flag to avoid pattern-based blocking
4. **Rate Limiting**: Use `-d` flag when scanning cloud providers or CDNs
5. **Local Testing**: Scan localhost first to verify setup

## Troubleshooting

### No Ports Found
- Check firewall rules
- Verify host is reachable (`ping <host>`)
- Try increasing timeout with `-t 10`

### Scan Too Slow
- Increase concurrency with `-c 500`
- Reduce timeout with `-t 1`
- Scan fewer ports

### Connection Refused
- Normal for closed ports
- Use `--show-closed` to see all results

### Too Many Open Files Error
- Reduce concurrency with `-c 50`
- Your OS has file descriptor limits

## Examples from Challenge Requirements

### Scan with Service Detection
```bash
./metron scan -H scanme.nmap.org -p 1-1000
```

### Subnet Scan
```bash
./metron scan -H 192.168.1.0/24 -p 22,80,443
```

### Full Scan with Statistics
```bash
./metron scan -H target.com --full -c 500
```

All scans include:
- ✅ Port status (open/closed/filtered)
- ✅ Service identification via banner
- ✅ Total scan time
- ✅ Concurrent execution with goroutines
