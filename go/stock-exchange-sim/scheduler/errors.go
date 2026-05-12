package scheduler

import "errors"

var (
	ErrNoMoreProcess        = errors.New("no more process doable")
	ErrSaveFile             = errors.New("save fiel error")
	ErrNoProcessFound       = errors.New("no process found")
	ErrRunProcOnCycle       = errors.New("run processes on cycle error")
	ErrPrintLog             = errors.New("print log error")
	ErrRunLog               = errors.New("run log error")
	ErrRunParallel          = errors.New("run parallel error")
	ErrSerialAlt            = errors.New("run serial alternate error")
	ErrSerialPipe           = errors.New("run serial pipeline error")
	ErrNotEnoughResources   = errors.New("not enough resources in stock")
	ErrRemoveRedundantProcs = errors.New("remove redundant processes")
	ErrReverseChildren      = errors.New("reverse children error")
	ErrRunChildProcess      = errors.New("run child process error")
)
