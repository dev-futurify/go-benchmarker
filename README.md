## GO Benchmarker

A benchmarking tool designed to measure and report on the performance of a system's CPU, memory, and disk (specifically SSD) under specified conditions. It utilizes Go's concurrency model, command-line interface (CLI) capabilities, and system command execution to accomplish its tasks.

At the start, the program defines command-line flags for the CPU load, memory size, and disk count using the flag package. These flags allow users to specify the parameters for the benchmark tests. The flags are parsed, and their values are validated to ensure they are positive integers. If the validation fails, the program prints usage information and exits.

The printSystemInfo function gathers and prints basic system information, including the operating system, architecture, CPU, memory, and disk information. It uses the exec.Command function to execute system commands (uname, grep, and df) and captures their output, which is then printed to the console.

The runBenchmarks function orchestrates the benchmarking process. It creates a channel named done to signal the completion of benchmarks. The CPU, memory, and disk benchmarks are initiated concurrently using Go routines. Each benchmark function receives the done channel and the relevant parameter (CPU load, memory size, or disk count) as arguments. The program waits for 10 seconds before signaling the completion of benchmarks by closing the done channel.

The CPU benchmark (runCPUBenchmark) measures CPU performance by calculating the Fibonacci sequence up to a specified number (representing CPU load). The memory benchmark (runMemoryBenchmark) measures memory bandwidth by copying data between slices, and the disk benchmark (runDiskBenchmark) measures disk (SSD) performance using the dd command to write a specified amount of data to a temporary file.

The disk benchmark uses regular expressions to parse the output of the dd command to extract the transfer rate, which is then converted to MB/s based on the unit of measurement provided in the output. If the parsing or conversion fails, the function prints an error message and returns a zero rate.

Throughout the code, error handling is minimal for simplicity. For example, the output of system commands is captured without checking for errors, and the disk benchmark assumes successful execution of the dd command.

This code demonstrates the use of Go's concurrency features (goroutines and channels), system command execution, regular expressions, and error handling to build a simple yet effective system benchmarking too