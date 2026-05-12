package domain

import "github.com/google/uuid"

type Resources map[string]int // resources keys are types ie. "euros, fame" and values are quantity for each.

type Scheduler struct {
	Optimize  []string
	Stock     Resources
	Processes []Process
	Time      bool
	Filepath  string
	Pipeline  bool
}

type Process struct {
	Name   string
	In     Resources // the types and ammount of resources consumed
	Cycles int
	Out    Resources // the types and amount of resources produced
	// NeededResourcesSorted []string
}

type Log struct {
	Id       uuid.UUID
	Cycle    int    `json:"cycle"`
	ProcName string `json:"proc_name"`
	Message  string `json:"message"`
}
