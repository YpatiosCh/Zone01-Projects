package scheduler

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"platform.zone01.gr/git/vstefano/stock-exchange-sim/domain"
)

type serial struct {
	log          []domain.Log     // final log shared among all serial runs
	parents      []domain.Process // the processes that produce the resource to optimize
	stock        domain.Resources
	allProcesses []domain.Process // all processes available
	depth        int              // recursion depth
	maxDepth     int              // recursion max depth
	drainStock   bool             // true if only optimizing time
}

type processState struct {
	current        *domain.Process // the process curently running
	children       []*processState // the dependencies of current
	startCycle     int
	targetResource string    // holds the resource that was needed for this process to run
	logId          uuid.UUID // the id of the current process log in s.logs
}

// Loops through 's.parrents' and calls serialRun. Keeps calling with the same parent until it returns 'ErrNoMoreProcess'.
// On each return it check the stock status and if the target in stock is not increasing it breaks the inner for loop
// When error is returned the loop continues by calling serialRun with the next in line until they all return 'ErrNoMoreProcess'.
func (s *serial) RunPipeline(ctx context.Context) (log []domain.Log, err error) {
	fmt.Println("Running pipeline")
	var currentCycle int
	var newCycle int

	previousStock := deepCopyMap(s.stock)

	var errOnLoop error
	for _, parent := range s.parents {
		for {
			ps := &processState{
				current:    &parent,
				startCycle: currentCycle,
			}

			s.depth = 0
			newCycle, errOnLoop = s.serialRun(ctx, ps)
			if errOnLoop != nil {
				err = errors.Join(errOnLoop, err)
				break
			}

			if !s.isStockIncreasing(&previousStock) && !s.drainStock {
				break
			}
			currentCycle = newCycle
		}
	}
	s.log = LogNoMoreProc(s.log, newCycle+1)
	if err != nil {
		return s.log, errors.Join(ErrSerialPipe, err)
	}
	return s.log, nil
}

// Loops through 's.parrents' processes for as long there is no process possible.
//
// After the current process returns one time the stock is passed to the next.
// If the process returns an error it is removed from future itterations.
//
// When all processes return error, it logs the next cycle of the end cycle of last process end cyle as 'no more process doable'.
func (s *serial) RunAlternate(ctx context.Context) (log []domain.Log, err error) {
	fmt.Println("Running alternate")
	var currentCycle int
	var i int
	var counter int
	var newCycle int
	previousStock := deepCopyMap(s.stock)

	var errOnLoop error
	for {
		if counter == len(s.parents) {
			s.log = LogNoMoreProc(s.log, newCycle+1)
			return s.log, errors.Join(ErrSerialAlt, err)
		}

		if i > len(s.parents)-1 {
			counter = 0
			i = 0
		}

		parent := s.parents[i]
		ps := &processState{current: &parent, startCycle: currentCycle}
		s.depth = 0
		newCycle, errOnLoop = s.serialRun(ctx, ps)
		if errOnLoop != nil {
			counter++
			err = errors.Join(errOnLoop, err)
			continue
		}

		if !s.isStockIncreasing(&previousStock) && !s.drainStock {
			continue
		}

		currentCycle = newCycle
		i++
	}
}

// Calls all needed processes to produce the needed in resources for 'p',
// 'endCycle' is the cycle just before an error occured down the process tree.
func (s *serial) serialRun(ctx context.Context, ps *processState) (endCycle int, err error) {
	if s.depth > s.maxDepth {
		return endCycle, fmt.Errorf("recursion depth limit exceeded")
	}
	s.depth++

	if ctx.Err() != nil {
		return ps.startCycle, errors.Join(noMoreProcess(ps.startCycle+1), ctx.Err())
	}

	// Call all workers that produce resources not in stock
	for resource, quantity := range ps.current.In {
		ps.targetResource = resource
		if s.stock[ps.targetResource] < quantity {
			ps.startCycle, err = s.runChildProcess(ctx, ps)
			if err != nil {
				if errR := s.reverseChildren(ctx, ps); errR != nil {
					err = errors.Join(err, errR)
				}
				return ps.startCycle, errors.Join(noMoreProcess(ps.startCycle+1), err)
			}
		}
	}

	for k := range ps.current.In {
		// Check if a child has consumed needed resources and if possible aquire them from a worker.
		if ps.current.In[k] > s.stock[k] {
			//don't try to acquire resource if it's a "worker" (same in and out)
			//workers are not produced by processes, just used and released
			//so running a process to get a worker makes no sense
			if ps.current.In[k] != ps.current.Out[k] {
				ps.targetResource = k
				ps.startCycle, err = s.runChildProcess(ctx, ps)
			}
			if err != nil {
				if errR := s.reverseChildren(ctx, ps); errR != nil {
					err = errors.Join(err, errR)
				}
				return ps.startCycle, errors.Join(noMoreProcess(ps.startCycle+1), err)
			}
		}
	}

	err = updateStock(*ps.current, s.stock)
	if err != nil {
		return endCycle, errors.Join(noMoreProcess(endCycle), err)
	}

	ps.logId = uuid.New()

	s.log = append(s.log, domain.Log{
		Id:       ps.logId,
		Cycle:    ps.startCycle,
		ProcName: ps.current.Name,
	})

	endCycle = ps.current.Cycles + ps.startCycle
	return endCycle, nil
}

// Runs all processes that provide the given resource in quantity. Returns error if this or a descendant process is unsuccesfull.
func (s *serial) runChildProcess(ctx context.Context, ps *processState) (endCycle int, err error) {
	child, err := findProcessByOutput(ps.targetResource, s.allProcesses)
	if err != nil {
		return ps.startCycle, errors.Join(fmt.Errorf("%w: parent: %s, target resource: %s ", ErrRunChildProcess, ps.current.Name, ps.targetResource), err)
	}

	neededInstances := max(1, ps.current.In[ps.targetResource]/child.Out[ps.targetResource])

	for range neededInstances {
		childPs := &processState{
			current:    child,
			startCycle: ps.startCycle,
		}

		c, err := s.serialRun(ctx, childPs)
		if err != nil {
			return ps.startCycle, err
		}
		ps.startCycle = c
		ps.children = append(ps.children, childPs)
	}
	return ps.startCycle, nil
}

// Restores 's.stock' for all children of 'ps' ins and outs. Deletes the from 's.logs' all children loggings.
// Restores the 'ps.startCycle' at the state it was before all children where run.
func (s *serial) reverseChildren(ctx context.Context, ps *processState) error {
	for _, child := range ps.children {
		if err := ctx.Err(); err != nil {
			return errors.Join(ErrReverseChildren, err)
		}
		s.reverseProcess(child)
		ps.startCycle -= child.current.Cycles
	}
	return nil
}

// Adds to 's.stock' all of 'ps.current' ins and removes all outs.
// Deletes the from 's.logs' the log with id 'ps.logId'.
func (s *serial) reverseProcess(ps *processState) {
	restoreStock(*ps.current, s.stock)
	for i, log := range s.log {
		if log.Id.ID() == ps.logId.ID() {
			s.log = append(s.log[:i], s.log[i+1:]...)
		}
	}
}

// Returns true if any of the target resources of optimize has increased on 's.stock'.
func (s *serial) isStockIncreasing(previousStock *domain.Resources) bool {
	atLeastOneIncreased := false
	prev := *previousStock
	for _, p := range s.parents {
		r := p.Name
		if prev[r] <= s.stock[r] {
			//fmt.Printf("%s previous: %v, current: %v\n", r, ps[r], s.stock[r])
			atLeastOneIncreased = true
		}
	}
	current := deepCopyMap(s.stock)
	*previousStock = current
	return atLeastOneIncreased
}
