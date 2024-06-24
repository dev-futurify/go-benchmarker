package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func main() {
	// Print basic system information
	printSystemInfo()

	// Create a channel to signal completion of benchmarks
	done := make(chan struct{})

	// Run benchmarks with a 10-second timeout
	go runBenchmarks(done)

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
	cpuInfo, _ := exec.Command("sysctl", "-n", "machdep.cpu.brand_string").Output()
	fmt.Printf("CPU Info: %s\n", cpuInfo)

	// Get memory information
	memInfo, _ := exec.Command("sysctl", "-n", "hw.memsize").Output()
	fmt.Printf("Memory Info: %s bytes\n", memInfo)

	// Get disk information
	diskInfo, _ := exec.Command("df", "-h").Output()
	fmt.Printf("Disk Info:\n%s", diskInfo)
}

func runBenchmarks(done chan struct{}) {
	// Benchmark CPU
	fmt.Println("\nBenchmarking CPU...")
	cpuBenchmark := runCPUBenchmark(done)
	fmt.Printf("CPU Benchmark Score: %d\n", cpuBenchmark)

	// Benchmark Memory
	fmt.Println("\nBenchmarking Memory...")
	memBenchmark := runMemoryBenchmark(done)
	fmt.Printf("Memory Benchmark Score: %.2f GB/s\n", memBenchmark)

	// Benchmark Disk (SSD) Performance
	fmt.Println("\nBenchmarking Disk (SSD) Performance...")
	diskBenchmark := runDiskBenchmark(done)
	fmt.Printf("Disk Benchmark Score: %.2f MB/s\n", diskBenchmark)
}

func runCPUBenchmark(done chan struct{}) int {
	// Measure CPU performance by calculating Fibonacci sequence
	fib := make(chan int)
	go func() {
		fib <- fibonacci(40)
	}()
	select {
	case <-done:
		return 0
	case result := <-fib:
		return result
	}
}

func runMemoryBenchmark(done chan struct{}) float64 {
	// Measure memory bandwidth by copying data between slices
	const dataSize = 1e9 // 1 GB
	data := make([]byte, dataSize)
	start := time.Now()
	for i := 0; i < dataSize; i++ {
		select {
		case <-done:
			return 0
		default:
			data[i] = 1
		}
	}
	elapsed := time.Since(start)
	return float64(dataSize) / (float64(elapsed) / float64(time.Second))
}

func runDiskBenchmark(done chan struct{}) float64 {
	// Measure disk (SSD) performance using the 'dd' command
	cmd := exec.Command("dd", "if=/dev/zero", "of=/tmp/testfile", "bs=1M", "count=1000", "conv=sync")
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
		// No conversion needed for MB
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
