package scheduler

import (
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"os"

	"platform.zone01.gr/git/vstefano/stock-exchange-sim/domain"
)

// ===================================================================
// Process Helpers
// ===================================================================

// Returns the first process found in s.processes that provides the resource.
// Returns 'ErrNoProcessFound' if no process provides the resource.
func findProcessByOutput(resource string, allProcesses []domain.Process) (*domain.Process, error) {
	for _, proc := range allProcesses {
		if _, ok := proc.Out[resource]; ok {
			return &proc, nil
		}
	}
	return nil, fmt.Errorf("%w: find process by output", ErrNoProcessFound)
}

// Returns the process with 'pName' from s.processes.
// Returns 'ErrNoProcessFound' if not found or no name is given.
func findProcessByName(pName string, processList []domain.Process) (proc domain.Process, err error) {
	if pName == "" {
		return domain.Process{}, fmt.Errorf("no process name given")
	}

	for _, p := range processList {
		if p.Name == pName {
			return p, nil
		}
	}
	return domain.Process{}, ErrNoProcessFound
}

// Checks s.optimize for optimization targets and adds them to 's.parents' slice.
// If s.optimize is empty and s.time is true then adds all processes to p.entryProc.
func (s *scheduler) setEntryProcs() error {
	// check for time
	if len(s.optimize) == 0 && s.time {
		s.parents = append(s.parents, s.allProcesses...)
		s.onlyOptimizingTime = true
		fmt.Println("Optimizing Time")
		return nil
	}

	// check for optimize resources
	for _, resource := range s.optimize {
		proc, err := findProcessByOutput(resource, s.allProcesses)
		if err != nil {
			continue
		}
		s.parents = append(s.parents, *proc)
	}

	if len(s.parents) == 0 {
		return ErrNoProcessFound
	}
	fmt.Printf("Optimizing %v\nOptimizing for time: %v\n", s.optimize, s.time)
	return nil
}

// ===================================================================
// Finalize
// ===================================================================

// Saves 's.log' to file in json format and prints it to terminal.
func (s *scheduler) finalize() error {
	log := s.log

	//save log to file
	if err := s.saveLogToFile(log); err != nil {
		return err
	}

	//print log to terminal
	s.printLog(log)
	return nil
}

// Saves 's.log' to file in json bytes.
func (s *scheduler) saveLogToFile(log []domain.Log) error {
	if s.filepath == "" {
		return fmt.Errorf("%w: invalid filepath", ErrSaveFile)
	}

	s.filepath += ".log"

	logFile, err := os.Create(s.filepath)
	if err != nil {
		return errors.Join(
			fmt.Errorf("%w: error saving log to file: %s, error", ErrSaveFile, s.filepath),
			err,
		)
	}
	defer logFile.Close()

	bytes, err := json.Marshal(log)
	if err != nil {
		return errors.Join(
			fmt.Errorf("%w: error marshaling s.Logs", ErrSaveFile),
			err,
		)
	}

	_, err = logFile.Write(bytes)
	if err != nil {
		return errors.Join(
			fmt.Errorf("%w: error writing %s to file %s", ErrSaveFile, string(bytes), s.filepath),
			err,
		)
	}
	return nil
}

// Pretty print 's.log' to terminal.
func (s *scheduler) printLog(logs []domain.Log) error {
	fmt.Println("Main Processes:")
	for i, entry := range logs {
		var proc, msg, cycle string

		cycle = fmt.Sprintf("%d", entry.Cycle)

		if entry.ProcName != "" {
			proc = entry.ProcName
		}

		if cycle != "" && proc != "" && i < 10000 { //cap messages printed in case of endless loop
			fmt.Printf("  %s: %s\n", cycle, proc)
		}

		if entry.Message != "" {
			msg = entry.Message
			fmt.Println(msg)
		}
	}

	fmt.Println("Stock:")
	for resource, quantity := range s.stock {
		fmt.Printf("  %s => %v\n", resource, quantity)
	}
	return nil
}

// ===================================================================
// Altering Stock
// ===================================================================

// Adds to 's' all p.Out values and removes from 's' all p.In.
// Returns error if any of p.In values are larger than s.[pIn.<key>]
func updateStock(p domain.Process, s domain.Resources) error {
	err := removeFromStock(p.In, s)
	if err != nil {
		return err
	}
	addToStock(p.Out, s)
	return nil
}

// The reverse of 'updateStock'. It is used when a process 'p' changes to 's' need to be undone.
// Reverse stock only logs error returned from removeFromStock as it is assumed that whatever is to be removed was added by this process before hand.
func restoreStock(p domain.Process, s domain.Resources) {
	err := removeFromStock(p.Out, s)
	if err != nil {
		fmt.Println("restoreStock", err.Error())
	}
	addToStock(p.In, s)
}

// Adds all 'r' values to 's' values with corresponding keys.
func addToStock(r domain.Resources, s domain.Resources) {
	for outResource, quantity := range r {
		s[outResource] += quantity
	}
}

// Removes all 'r' values from 's' values with corresponding keys. Returns error if any of r values are larger than the 's' value.
// Although error is returned the removal takes place allowing for the error to be ignored if must.
func removeFromStock(r domain.Resources, s domain.Resources) (err error) {
	for inResource, quantity := range r {
		if s[inResource] < quantity {
			err = fmt.Errorf("%w: for resource %s", ErrNotEnoughResources, inResource)
		}
		s[inResource] -= quantity
	}
	return err
}

// Returns true if there are enough resources on stock 's' to run process with 'pIn' input.
func enoughStock(pIn domain.Resources, s domain.Resources) bool {
	for k := range pIn {
		if pIn[k] > s[k] {
			return false
		}
	}
	return true
}

// ===================================================================
// Deep Copy Helpers
// ===================================================================

func deepCopyMap(source domain.Resources) domain.Resources {
	copy := make(map[string]int, len(source))
	maps.Copy(copy, source)
	return copy
}

// ===================================================================
// Logging and Errors
// ===================================================================

func noMoreProcess(cycle int) error {
	return fmt.Errorf("%w at cycle %v", ErrNoMoreProcess, cycle)
}

func LogNoMoreProc(log []domain.Log, cycle int) []domain.Log {
	msg := fmt.Sprintf("No more process doable at cycle %v", cycle)
	log = append(log, domain.Log{
		Message: msg,
	})
	return log
}
