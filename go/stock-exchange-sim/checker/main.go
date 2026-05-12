package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"platform.zone01.gr/git/vstefano/stock-exchange-sim/domain"
	"platform.zone01.gr/git/vstefano/stock-exchange-sim/parser"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("Usage: %s <config_file> <log_file>\n", os.Args[0])
		return
	}

	configFile := os.Args[1]
	logFile := os.Args[2]

	// Parse the configuration file
	p := parser.NewParser()
	scheduler, err := p.Run(configFile)
	if err != nil {
		fmt.Printf("Error parsing configuration file: %v\n", err)
		return
	}

	// Read and parse the log file
	logs, err := readLogFile(logFile)
	if err != nil {
		fmt.Printf("Error reading log file: %v\n", err)
		return
	}

	// Validate the log
	if err := validateLog(scheduler, logs); err != nil {
		fmt.Println("Error detected")
		fmt.Println(err)
		return
	}

	fmt.Println("Trace completed, no error detected.")
}

// readLogFile reads and parses the log file
func readLogFile(filepath string) ([]domain.Log, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read log file: %w", err)
	}

	var logs []domain.Log

	// Try to parse as JSON first
	var jsonLogs []domain.Log
	if err := json.Unmarshal(data, &jsonLogs); err == nil {
		// Filter out empty process names
		for _, log := range jsonLogs {
			if log.ProcName != "" {
				logs = append(logs, log)
			}
		}
		return logs, nil
	}

	// If not JSON, parse as text format (cycle:process_name)
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Skip "No more process doable" messages
		if strings.HasPrefix(line, "No more process") {
			continue
		}

		// Parse cycle:process_name format
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		cycle, err := strconv.Atoi(strings.TrimSpace(parts[0]))
		if err != nil {
			return nil, fmt.Errorf("invalid cycle number in line: %s", line)
		}

		procName := strings.TrimSpace(parts[1])
		if procName == "" {
			continue
		}

		logs = append(logs, domain.Log{
			Cycle:    cycle,
			ProcName: procName,
		})
	}

	return logs, nil
}

// validateLog simulates the execution of processes and validates stock availability
func validateLog(scheduler *domain.Scheduler, logs []domain.Log) error {
	// Create a copy of the initial stock
	stock := make(domain.Resources)
	for k, v := range scheduler.Stock {
		stock[k] = v
	}

	// Create a map of processes by name for quick lookup
	processMap := make(map[string]domain.Process)
	for _, proc := range scheduler.Processes {
		processMap[proc.Name] = proc
	}

	// Track processes in progress (completion_cycle -> list of processes completing at that cycle)
	processingQueue := make(map[int][]domain.Process)

	for _, log := range logs {
		fmt.Printf("Evaluating: %d:%s\n", log.Cycle, log.ProcName)

		// Complete all processes that finish at or before current cycle
		for cycle := 0; cycle <= log.Cycle; cycle++ {
			if procs, exists := processingQueue[cycle]; exists {
				for _, proc := range procs {
					// Add outputs to stock
					for resource, quantity := range proc.Out {
						stock[resource] += quantity
					}
				}
				delete(processingQueue, cycle)
			}
		}

		// Find the process
		proc, exists := processMap[log.ProcName]
		if !exists {
			return fmt.Errorf("at %d:%s process not found in configuration", log.Cycle, log.ProcName)
		}

		// Check if we have enough stock to start this process
		for resource, needed := range proc.In {
			available := stock[resource]
			if available < needed {
				fmt.Println("Exiting...")
				return fmt.Errorf("at %d:%s stock insufficient", log.Cycle, log.ProcName)
			}
		}

		// Consume inputs from stock
		for resource, quantity := range proc.In {
			stock[resource] -= quantity
		}

		// Schedule the process to complete at (current cycle + process cycles)
		completionCycle := log.Cycle + proc.Cycles
		processingQueue[completionCycle] = append(processingQueue[completionCycle], proc)
	}

	// Complete all remaining processes
	for _, procs := range processingQueue {
		for _, proc := range procs {
			for resource, quantity := range proc.Out {
				stock[resource] += quantity
			}
		}
	}

	return nil
}
