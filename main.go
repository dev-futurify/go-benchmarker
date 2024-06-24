package main

import (
	"flag"
	"fmt"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func main() {
	// Define command-line flags
	cpuLoad := flag.Int("cpu", 0, "CPU load value (positive integer)")
	memSize := flag.Int("memory", 0, "Memory size in MB (positive integer)")
	distCount := flag.Int("disk", 0, "Disk Count Round (positive integer)")

	// Parse command-line flags
	flag.Parse()

	// Validate flags
	if *cpuLoad <= 0 || *memSize <= 0 || *distCount <= 0 {
		fmt.Println("Usage: ./main --cpu <CPU> --memory <memory> --disk <disk>")
		fmt.Println("All values must be positive integers.")
		return
	}

	// Print basic system information
	printSystemInfo()

	// Create a channel to signal completion of benchmarks
	done := make(chan struct{})

	// Run benchmarks with a 10-second timeout
	go runBenchmarks(done, *cpuLoad, *memSize, *distCount)

	// Wait for 10 seconds and then signal completion
	select {
	case <-time.After(10 * time.Second):
		fmt.Println("Stopping benchmarks...")
		close(done)
	}
}

func printSystemInfo() {
	// Print basic system information
	fmt.Printf("OS: %s\n", runtime.GOOS)
	fmt.Printf("Architecture: %s\n", runtime.GOARCH)

	// Get CPU information
	cpuInfo, _ := exec.Command("uname", "-p").Output()
	fmt.Printf("CPU Info: %s\n", cpuInfo)

	// Get memory information
	memInfo, _ := exec.Command("grep", "MemTotal", "/proc/meminfo").Output()
	fmt.Printf("Memory Info: %s", memInfo)

	// Get disk information
	diskInfo, _ := exec.Command("df", "-h").Output()
	fmt.Printf("Disk Info:\n%s", diskInfo)
}

func runBenchmarks(done chan struct{}, cpuLoad int, memSize int, distCount int) {
	// Benchmark CPU
	fmt.Println("\nBenchmarking CPU...")
	cpuBenchmark := runCPUBenchmark(done, cpuLoad)
	fmt.Printf("CPU Benchmark Score: %d\n", cpuBenchmark)

	// Benchmark Memory
	fmt.Println("\nBenchmarking Memory...")
	memBenchmark := runMemoryBenchmark(done, memSize)
	fmt.Printf("Memory Benchmark Score: %.2f MB/s\n", memBenchmark)

	// Benchmark Disk (SSD) Performance
	fmt.Println("\nBenchmarking Disk (SSD) Performance...")
	diskBenchmark := runDiskBenchmark(done, distCount)
	fmt.Printf("Disk Benchmark Score: %.2f MB/s\n", diskBenchmark)
}

func runCPUBenchmark(done chan struct{}, cpuLoad int) int {
	// Measure CPU performance by calculating Fibonacci sequence
	fib := make(chan int)
	go func() {
		fib <- fibonacci(cpuLoad)
	}()
	select {
	case <-done:
		return 0
	case result := <-fib:
		return result
	}
}

func runMemoryBenchmark(done chan struct{}, memSize int) float64 {
	// Measure memory bandwidth by copying data between slices
	data := make([]byte, memSize*1024*1024) // Convert MB to bytes
	start := time.Now()
	for i := 0; i < len(data); i++ {
		select {
		case <-done:
			return 0
		default:
			data[i] = 1
		}
	}
	elapsed := time.Since(start)
	return float64(len(data)) / (float64(elapsed) / float64(time.Second)) / (1024 * 1024) // Convert to MB/s
}

func runDiskBenchmark(done chan struct{}, distCount int) float64 {
	// Measure disk (SSD) performance using the 'dd' command
	cmd := exec.Command("dd", "if=/dev/zero", "of=/tmp/testfile", "bs=1M", fmt.Sprintf("count=%d", distCount), "conv=sync")

	cmd.Stderr = nil
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error running disk benchmark:", err)
		return 0
	}

	// Use regular expression to extract any numeric value followed by a unit
	rateRegexp := regexp.MustCompile(`(\d+(\.\d+)?)\s*(\D+)/s`)
	matches := rateRegexp.FindStringSubmatch(string(output))
	if len(matches) < 4 {
		fmt.Println("Error parsing transfer rate: rate not found in output")
		return 0
	}

	// Parse the transfer rate to MB/s
	rate, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		fmt.Println("Error parsing transfer rate:", err)
		return 0
	}

	unit := strings.ToUpper(matches[3])
	switch unit {
	case "B":
		rate /= 1e6
	case "KB":
		rate /= 1e3
	case "MB":
	case "GB":
		rate *= 1e3
	default:
		fmt.Printf("Warning: Unsupported unit '%s', assuming MB/s\n", unit)
		unit = "MB"
	}

	return rate
}

func fibonacci(n int) int {
	if n <= 0 {
		return 0
	} else if n == 1 {
		return 1
	}
	return fibonacci(n-1) + fibonacci(n-2)
}
