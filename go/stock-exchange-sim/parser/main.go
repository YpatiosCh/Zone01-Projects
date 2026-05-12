package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"platform.zone01.gr/git/vstefano/stock-exchange-sim/domain"
)

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Run(filepathString string) (*domain.Scheduler, error) {
	// read the file
	systemPath := filepath.FromSlash(filepathString)
	lines, err := p.readFile(systemPath)
	if err != nil {
		return nil, err
	}

	// handle comments
	lines = p.handleComments(lines)

	// check file format
	stock, process, optimize, pipeline, err := p.checkFormat(lines)
	if err != nil {
		return nil, err
	}

	// create scheduler
	scheduler, err := p.createScheduler(stock, process, optimize, systemPath, pipeline)
	if err != nil {
		return nil, err
	}

	return scheduler, nil
}

// readFile reads the file and returns a []string with the containing lines
func (p *Parser) readFile(filepath string) ([]string, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	return strings.Split(string(data), "\n"), nil
}

// handleComments removes lines that start with '#' and trims inline comments.
func (p *Parser) handleComments(lines []string) []string {
	var result []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Skip full-line comments
		if strings.HasPrefix(trimmed, "#") {
			continue
		}
		// Remove inline comments
		if idx := strings.Index(trimmed, "#"); idx != -1 {
			trimmed = strings.TrimSpace(trimmed[:idx])
		}
		// Skip if the line becomes empty after removing comments
		if trimmed == "" {
			continue
		}
		result = append(result, trimmed)
	}

	return result
}

// checkFormat verifies that the configuration file follows this format:
//
// Stock: at least one (<stock_name>:<quantity>)
//
// Process: at least one (<name>:(<need>:<quantity>;...):(<result>:<quantity>;...):<nb_cycle>)
//
// Optimize: exactly one  (optimize:(<stock_name>|time))
//
// Any deviation (wrong order, missing section, or invalid format) returns an error.
// Returns stock, process and optimize lines to create scheduler
func (p *Parser) checkFormat(file []string) ([]string, []string, string, bool, error) {
	// formats
	stock := regexp.MustCompile(`^[A-Za-z_]+\s*:\s*[0-9]+$`)
	process := regexp.MustCompile(`^\s*[A-Za-z_]+\s*:\s*\(\s*(?:[A-Za-z_]+\s*:\s*[0-9]+\s*)(?:\s*;\s*[A-Za-z_]+\s*:\s*[0-9]+\s*)*\)\s*:\s*\(\s*(?:[A-Za-z_]+\s*:\s*[0-9]+\s*)(?:\s*;\s*[A-Za-z_]+\s*:\s*[0-9]+\s*)*\)\s*:\s*[0-9]+\s*$`)
	optimize := regexp.MustCompile(`^\s*optimize\s*:\s*\(\s*[A-Za-z_]+(?:\s*;\s*[A-Za-z_]+)*\s*\)\s*$`)

	hasStock := false
	hasProcess := false
	hasOptimize := false

	// current expected section
	mode := "stock"

	var stockLines []string
	var processLines []string
	var optimizeLines string
	pipeline := true

	for i, line := range file {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		switch mode {
		case "stock":
			if stock.MatchString(line) {
				hasStock = true
				stockLines = append(stockLines, line)
				continue
			}
			if process.MatchString(line) {
				if !hasStock {
					return stockLines, processLines, "", pipeline, fmt.Errorf("missing stock section before process line: %v", line)
				}
				mode = "process"
				hasProcess = true
				processLines = append(processLines, line)
				continue
			}
			if optimize.MatchString(line) {
				return stockLines, processLines, "", pipeline, fmt.Errorf("optimize section appeared before process section at line %d", i+1)
			}
			return stockLines, processLines, "", pipeline, fmt.Errorf("invalid stock format: %v", line)

		case "process":
			if process.MatchString(line) {
				// hasProcess = true    <- no need maybe
				processLines = append(processLines, line)
				continue
			}
			if optimize.MatchString(line) {
				if !hasProcess {
					return stockLines, processLines, "", pipeline, fmt.Errorf("missing process section before optimize line: %v", line)
				}
				mode = "optimize"
				hasOptimize = true
				optimizeLines = line
				continue
			}
			if stock.MatchString(line) {
				return stockLines, processLines, "", pipeline, fmt.Errorf("stock line found after process section: %v", line)
			}
			// not pass process, not pass stock, means expecting optimize
			return stockLines, processLines, "", pipeline, fmt.Errorf("invalid optimize format: %v", line)

		case "optimize":
			if strings.ToLower(line) == "alternate" {
				pipeline = false
				mode = "done"
			} else if strings.ToLower(line) == "pipeline" {
				mode = "done"
				continue
			} else {
				return stockLines, processLines, "", pipeline, fmt.Errorf("unexpected line after optimize section: %v", line)
			}

		case "done":
			return stockLines, processLines, "", pipeline, fmt.Errorf("unexpected line after optimize section: %v", line)
		}
	}

	// final validation: all required sections must exist, in order
	if !hasStock {
		return stockLines, processLines, "", pipeline, fmt.Errorf("missing stock section")
	}
	if !hasProcess {
		return stockLines, processLines, "", pipeline, fmt.Errorf("missing process section")
	}
	if !hasOptimize {
		return stockLines, processLines, "", pipeline, fmt.Errorf("missing optimize section")
	}

	return stockLines, processLines, optimizeLines, pipeline, nil
}

func (p *Parser) createScheduler(stockLines, process []string, optimizeLine, systemPath string, pipeline bool) (*domain.Scheduler, error) {
	// retrieve stock from file
	stock, err := p.stock(stockLines)
	if err != nil {
		return nil, err
	}

	// retrieve processes from file
	processes, err := p.processes(process)
	if err != nil {
		return nil, err
	}

	// retrieve entryProc and time
	optimize, time, err := p.optimize(*processes, optimizeLine)
	if err != nil {
		return nil, err
	}

	return &domain.Scheduler{
		Optimize:  optimize,
		Stock:     *stock,
		Processes: *processes,
		Time:      time,
		Filepath:  systemPath,
		Pipeline:  pipeline,
	}, nil
}

// stock retrieves the stock resources from the file given
func (p *Parser) stock(stockLines []string) (*domain.Resources, error) {
	mStock := make(domain.Resources)

	for _, stock := range stockLines {
		split := strings.Split(stock, ":")
		quantity, err := strconv.Atoi(string(strings.TrimSpace(split[1])))
		if err != nil {
			return nil, fmt.Errorf("invalid stock quantity: %v", stock)
		}
		mStock[string(strings.TrimSpace(split[0]))] = quantity
	}

	return &mStock, nil
}

// processes retrieves the processes from the file given
func (p *Parser) processes(processLines []string) (*[]domain.Process, error) {
	var processes []domain.Process

	// yolo:(euro:8):(material:1):10
	// do_cabinet:(doorknobs:2;background:1;shelf:3):(cabinet:1):30

	// 0 name
	// 1 in
	// 2 out
	// 3 cycles

	// sell_phone:(smartphone:1):(euro:100):2
	// sell_phone: smartphone.1 : euro.100 :2
	for i, proc := range processLines {
		inParenth := false
		runes := []rune(proc)
		for j, ch := range runes {
			if inParenth && ch == ':' {
				runes[j] = '.'
			}
			if ch == '(' {
				runes[j] = ' '
				inParenth = true
			}
			if ch == ')' {
				runes[j] = ' '
				inParenth = false
			}
		}
		processLines[i] = string(runes)
	}

	for _, proc := range processLines {
		var process domain.Process
		splitProc := strings.Split(proc, ":")

		// name of each process
		process.Name = strings.TrimSpace(splitProc[0])

		// ins of each process
		ins := strings.Split(splitProc[1], ";")
		inResource := make(domain.Resources)
		for _, in := range ins {
			split := strings.Split(in, ".")
			quantity, err := strconv.Atoi(strings.TrimSpace(split[1]))
			if err != nil {
				return nil, fmt.Errorf("invalid resource quantity: %v", in)
			}
			inResource[string(strings.TrimSpace(split[0]))] = quantity
		}
		process.In = inResource

		// outs of each process
		outs := strings.Split(splitProc[2], ";")
		outResource := make(domain.Resources)
		for _, out := range outs {
			split := strings.Split(out, ".")
			quantity, err := strconv.Atoi(strings.TrimSpace(split[1]))
			if err != nil {
				return nil, fmt.Errorf("invalid resource quantity: %v", out)
			}
			outResource[string(strings.TrimSpace(split[0]))] = quantity
		}
		process.Out = outResource

		// cycles of each process
		cycles, err := strconv.Atoi(strings.TrimSpace(splitProc[3]))
		if err != nil {
			return nil, fmt.Errorf("invalid process cyle: %v : %v ", proc, splitProc[3])
		}
		process.Cycles = cycles

		processes = append(processes, process)
	}

	return &processes, nil
}

// optimize retrieves the entryProc process and time.
// if only time is provided or no entryProc (or multiple???) found, we return error
func (p *Parser) optimize(processes []domain.Process, optimize string) ([]string, bool, error) {
	hasTime := false

	// optimize:(time;cabinet)
	// time;cabinet;yolo
	inParenth := false
	var opt string
	for i, ch := range optimize {
		if i == len(optimize)-1 {
			break
		}
		if inParenth {
			opt += string(ch)
		}
		if ch == '(' {
			inParenth = true
		}
	}

	// trim spaces
	opts := strings.Split(opt, ";")
	for i, spl := range opts {
		opts[i] = strings.TrimSpace(spl)
	}

	var targets []string
	// find out what to optimize
	for _, target := range opts {
		if target == "time" {
			hasTime = true
		} else {
			if !canOptimize(target, processes) {
				return nil, false, fmt.Errorf("No process found resulting to %v. Make sure to optimize with a doable result.", target)
			}
			targets = append(targets, target)
		}
	}

	// fmt.Println("TARGETS:")
	// fmt.Println(targets)

	return targets, hasTime, nil
}

func canOptimize(target string, processes []domain.Process) bool {
	for _, proc := range processes {
		for k := range proc.Out {
			if k == target {
				return true
			}
		}
	}
	return false
}
