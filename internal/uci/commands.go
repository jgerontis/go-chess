package uci

import "fmt"

// Common UCI commands and helpers

// InitializeEngine sends the standard UCI initialization sequence
func (c *Client) InitializeEngine() error {
	if err := c.SendCommand("uci"); err != nil {
		return err
	}
	return nil
}

// IsReady sends the isready command
func (c *Client) IsReady() error {
	return c.SendCommand("isready")
}

// SetPosition sets the current position using FEN or startpos
func (c *Client) SetPosition(fen string) error {
	if fen == "" || fen == "startpos" {
		return c.SendCommand("position startpos")
	}
	return c.SendCommand(fmt.Sprintf("position fen %s", fen))
}

// SetPositionWithMoves sets position and applies moves
func (c *Client) SetPositionWithMoves(fen string, moves []string) error {
	var cmd string
	if fen == "" || fen == "startpos" {
		cmd = "position startpos"
	} else {
		cmd = fmt.Sprintf("position fen %s", fen)
	}
	
	if len(moves) > 0 {
		cmd += " moves"
		for _, move := range moves {
			cmd += " " + move
		}
	}
	
	return c.SendCommand(cmd)
}

// GoDepth starts search with specified depth
func (c *Client) GoDepth(depth int) error {
	return c.SendCommand(fmt.Sprintf("go depth %d", depth))
}

// GoMovetime starts search for specified time in milliseconds
func (c *Client) GoMovetime(ms int) error {
	return c.SendCommand(fmt.Sprintf("go movetime %d", ms))
}

// GoInfinite starts infinite search
func (c *Client) GoInfinite() error {
	return c.SendCommand("go infinite")
}

// Stop stops the current search
func (c *Client) Stop() error {
	return c.SendCommand("stop")
}

// Quit sends quit command to engine
func (c *Client) Quit() error {
	return c.SendCommand("quit")
}