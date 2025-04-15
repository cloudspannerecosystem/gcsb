# GCSB - Google Cloud Spanner Benchmark

A tool for benchmarking Google Cloud Spanner performance.

## Features

- Configurable workload patterns
- Support for various Spanner configurations
- Real-time metrics and reporting
- Commit delay optimization support
- Easy instance creation and testing setup
- Automated dual-region testing

## Prerequisites

- Go 1.19 or later
- Google Cloud SDK
- Authenticated gcloud session
- Project with Spanner API enabled
- Permission to create Spanner instances

## Installation

There are two ways to install GCSB:

### Option 1: Install from source (Recommended)
```bash
# Clone the repository
git clone https://github.com/cloudspannerecosystem/gcsb.git
cd gcsb

# Build the binary
go build -o gcsb

# Verify the installation
./gcsb --help
```

### Option 2: Install using go install (Advanced)
```bash
go install github.com/cloudspannerecosystem/gcsb@latest
```
Note: If using `go install`, you'll need to update the paths in the test scripts to use the full path to your GCSB binary (usually in `$GOPATH/bin/gcsb`).

## Quick Start

1. Ensure you have the Google Cloud SDK installed and are authenticated:
```bash
# Install Google Cloud SDK if not already installed
# Visit https://cloud.google.com/sdk/docs/install for installation instructions

# Authenticate with Google Cloud
gcloud auth login

# Set your project ID
gcloud config set project YOUR_PROJECT_ID
```

2. Build the GCSB binary (if not already done):
```bash
go build -o gcsb
```

3. Create test instances using the provided script:
```bash
chmod +x create_test_instances.sh
./create_test_instances.sh
```

The script will:
- Show available Spanner configurations
- Let you select which configurations to create
- Create instances and databases with test schemas
- Provide example commands for running benchmarks

Available configurations include:
- Single-region (us-central1)
- Multi-region (nam3)
- Dual-region (australia-southeast1+australia-southeast2)
- Both STANDARD and ENTERPRISE versions

### Dual-Region Testing

For testing dual-region configurations with Enterprise Plus edition:

```bash
# First, make sure you've built the binary (if not already done)
go build -o gcsb

# Then run the test script
chmod +x run_dual_region_test.sh
./run_dual_region_test.sh
```

This script will:
- Create an Enterprise Plus instance in dual-region configuration
- Create a database with test schema
- Load initial test data
- Run a benchmark with commit delay optimization
- Save results to a timestamped file

The test configuration includes:
- 4 concurrent threads
- 100,000 operations (~5GB of data)
- 100ms commit delay
- 100% write operations
- Large string fields to help reach target data size

## Troubleshooting

Common issues and solutions:

1. "cannot execute binary file: Exec format error"
   - Make sure you've built the binary for your system using `go build -o gcsb`
   - If using `go install`, update the script to use the full path to your GCSB binary (usually in `$GOPATH/bin/gcsb`)

2. "Permission denied" when running scripts
   - Make sure the scripts are executable: `chmod +x *.sh`

3. "Project not found" or "Permission denied" for Spanner operations
   - Verify your gcloud authentication: `gcloud auth login`
   - Check your project ID: `gcloud config get-value project`
   - Ensure you have the necessary Spanner permissions

4. "Instance already exists" error
   - The script will automatically use existing instances if they match the naming pattern
   - To create a new instance, either delete the existing one or modify the instance name in the script

## Configuration Options

### Workload Configuration

- `threads`: Number of concurrent worker threads
- `operations`: Total number of operations to perform
- `commitDelay`: Maximum commit delay for write operations (e.g., "100ms")
- `tables`: List of tables to benchmark
  - `name`: Table name
  - `writeRatio`: Ratio of write operations (0.0 to 1.0)
  - `readRatio`: Ratio of read operations (0.0 to 1.0)
  - `rowCount`: Number of rows to operate on

### Data Generation

- `column`: Column name
- `type`: Data type (uuid, string, int)
- For strings:
  - `length`: Maximum string length
- For integers:
  - `min`: Minimum value
  - `max`: Maximum value

## Performance Optimization

The tool supports commit delay optimization for write operations. This feature:
- Buffers mutations in transactions
- Delays commits to improve throughput
- Reduces the number of round trips to Spanner

To enable commit delay, set the `commitDelay` parameter in your configuration:
```yaml
workload:
  commitDelay: 100ms  # Maximum delay for commits
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the Apache License 2.0 - see the LICENSE file for details.
