package scheduler

import (
	"context"
	"errors"
	"fmt"
	"math"

	"platform.zone01.gr/git/vstefano/stock-exchange-sim/domain"
)

// Holds the process and its end cycle as calculated by start cycle and process.Cycles.
type processToComplete struct {
	process  *domain.Process
	endCycle int
}

type timeOptimizer struct {
	stock               domain.Resources
	processList         []domain.Process
	sourceLog           []domain.Log
	timeOptimizedLog    []domain.Log
	processesToComplete []processToComplete
}

// Runs all processes logged on 'sourceLog' maximizing the concurrent calls for each cycle respecting the ctx timeout.
// It depends on a sourceLog, stock and processList that results to a succesfull serial run of processes.
//
// On success returns a new time optimized log and ErrNoMoreProcess.
func (t *timeOptimizer) Run(ctx context.Context) (log []domain.Log, err error) {
	fmt.Println("Running time optimization")
	err = t.runLog(ctx)
	if err != nil {
		return t.timeOptimizedLog, errors.Join(ErrRunParallel, err)
	}
	return t.timeOptimizedLog, err
}

func (t *timeOptimizer) runLog(ctx context.Context) (err error) {
	var startCycle, endCycle int

	for {
		endCycle, err = t.runProcOnCycle(ctx, startCycle)
		if err != nil {
			return errors.Join(ErrRunLog, err)
		}

		nextCycle, found := t.nextEndCycleAfter(startCycle)
		if !found {
			t.timeOptimizedLog = LogNoMoreProc(t.timeOptimizedLog, endCycle+1)
			return noMoreProcess(endCycle + 1) // should return nil or err ??
		}
		startCycle = nextCycle
	}
}

// Adds to stock all outs from processes that are marked to end on 'cycle'.
// Runs all processes that the current state of stock allows.
// Adds each process to 'processesToComplete' for future cycle stock updates.
// Returns maxEndCycle
func (t *timeOptimizer) runProcOnCycle(ctx context.Context, cycle int) (maxEndCycle int, err error) {
	if ctx.Err() != nil {
		return cycle, errors.Join(ErrRunProcOnCycle, ctx.Err())
	}

	// Add to stock the the outs of processes that end on current cycle and remove them from 't.processesToComplete'
	for i := 0; i < len(t.processesToComplete); {
		p := t.processesToComplete[i]
		if p.endCycle == cycle {
			addToStock(p.process.Out, t.stock)
			t.processesToComplete = append(t.processesToComplete[:i], t.processesToComplete[i+1:]...)
			continue
		}
		i++
	}

	for i := 0; i < len(t.sourceLog); {
		currentLog := t.sourceLog[i]
		p, err := findProcessByName(currentLog.ProcName, t.processList)
		if err != nil {
			if errors.Is(err, ErrNoProcessFound) { // gave a process name and was not found
				// log.Fatal(err)
				return maxEndCycle, err
			}
			i++
			continue
		}

		if !enoughStock(p.In, t.stock) {
			i++
			continue
		}

		t.timeOptimizedLog = append(t.timeOptimizedLog, domain.Log{
			Cycle:    cycle,
			ProcName: p.Name,
		})

		maxEndCycle = max(p.Cycles, maxEndCycle)
		t.sourceLog = append(t.sourceLog[:i], t.sourceLog[i+1:]...)
		t.processesToComplete = append(t.processesToComplete,
			processToComplete{process: &p, endCycle: p.Cycles + cycle})

		err = removeFromStock(p.In, t.stock)
		if err != nil {
			return maxEndCycle, err
		}
	}

	return maxEndCycle + cycle, nil
}

// Returns the next bigger than 'cycle' end cycle in 'processesToComplete'
func (t *timeOptimizer) nextEndCycleAfter(cycle int) (nextCycle int, found bool) {
	nextCycle = math.MaxInt
	for _, p := range t.processesToComplete {
		if p.endCycle > cycle && p.endCycle < nextCycle {
			nextCycle = p.endCycle
			found = true
		}
	}
	return nextCycle, found
}

func NewTimeOptimizer(stock domain.Resources, sourceLog []domain.Log, processList []domain.Process) *timeOptimizer {
	return &timeOptimizer{
		sourceLog:   sourceLog,
		stock:       stock,
		processList: processList,
	}
}
