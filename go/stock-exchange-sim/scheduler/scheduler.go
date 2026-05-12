package scheduler

import (
	"context"
	"errors"
	"fmt"
	"time"

	"platform.zone01.gr/git/vstefano/stock-exchange-sim/domain"
)

type scheduler struct {
	optimize           []string         // the products that are targets for optimization.
	parents            []domain.Process // the processes that create the target products.
	time               bool             // determines if the time optimization is run.
	filepath           string           // log file path.
	log                []domain.Log     // holds each process run.
	Pipeline           bool             // Pipeline runs first proc given until it returns error, then moves to the next. Else loop through procs and execute each once until they return error.
	onlyOptimizingTime bool             // if true simply consume stock.
	stock              domain.Resources
	allProcesses       []domain.Process // All processes including parents.
}

// Runs all processes necessary to produce as many as possible of the products listed in s.optimize
// and generates a log of the executed processes.
//
// Each log entry contains the process that was run and the cycle in which the process started.
// If resources are insufficient to run more processes, the message "no more process possible doable" is logged.
//
// Logs are saved to a JSON file specified by s.filepath.
//
// If s.time is true, a second pass is executed in which, for each cycle, all processes that can run
// are executed.
//
// If multiple optimization targets are provided, the pipeline option determines whether to prioritize
// targets in the order they appear in s.optimize, or to give equal weight to all processes.
//
// Always returns error.
//
// Unless input is malformed or falls on infinite loop the expected error is 'ErrNoMoreProcess'
// -similar to EOF- as the stock is guarnateed to be finite.
//
// Usage:
//
//	 sc := scheduler.NewScheduler(s) // construct s using package domain type Scheduller
//	 err = sc.Optimize(context.Background(), timeout)
//	 if !errors.Is(err, scheduler.ErrNoMoreProcess) && err != nil {
//		 fmt.Println(err)
//	 }
func (s *scheduler) Optimize(ctx context.Context, timeout time.Duration) (err error) {
	fmt.Println("Running", s.filepath)

	childCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	defer func() {
		if ferr := s.finalize(); ferr != nil {
			if err != nil {
				err = errors.Join(err, ferr)
			} else {
				err = ferr
			}
		}
	}()

	err = s.setEntryProcs()
	if err != nil {
		return err
	}

	// Make a copy of original stock to use for paralell
	initialStock := deepCopyMap(s.stock)

	serial := s.newSerial()
	if s.Pipeline {
		s.log, err = serial.RunPipeline(childCtx)
	} else {
		s.log, err = serial.RunAlternate(childCtx)
	}

	if !s.time || errors.Is(err, context.DeadlineExceeded) {
		return err
	}

	fmt.Println("Serial complete")

	// s.printLog(s.log) // print serial log
	t := NewTimeOptimizer(initialStock, s.log, s.allProcesses)
	s.log, err = t.Run(childCtx)
	return err
}

func (s *scheduler) newSerial() *serial {
	return &serial{
		parents:      s.parents,
		stock:        s.stock,
		allProcesses: s.allProcesses,
		log:          s.log,
		depth:        0,
		maxDepth:     1000 * 1000,
		drainStock:   s.onlyOptimizingTime,
	}
}

func NewScheduler(s *domain.Scheduler) *scheduler {
	return &scheduler{
		stock:        s.Stock,
		allProcesses: s.Processes,
		time:         s.Time,
		filepath:     s.Filepath,
		Pipeline:     s.Pipeline,
		optimize:     s.Optimize,
	}
}
