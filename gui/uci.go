package gui

import (
	"bufio"
	"os/exec"
	"strings"
)

type UCIEngine struct {
	cmd    *exec.Cmd
	stdin  *bufio.Writer
	stdout *bufio.Scanner
}

// creates a new UCI Engine
func NewUCIEngine(enginePath string) (*UCIEngine, error) {
	cmd := exec.Command(enginePath)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(stdout)
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	engine := &UCIEngine{
		cmd:    cmd,
		stdin:  bufio.NewWriter(stdin),
		stdout: scanner,
	}

	// Initialize UCI
	// engine.sendCommand("uci")
	// engine.waitFor("uciok")

	return engine, nil
}

func (e *UCIEngine) SetPosition(moves []string) {
	cmd := "position startpos"
	if len(moves) > 0 {
		cmd += " moves " + strings.Join(moves, " ")
	}
	// e.sendCommand(cmd)
}
