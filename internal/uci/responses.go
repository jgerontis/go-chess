package uci

import (
	"strconv"
	"strings"
)

// Response types for parsing UCI engine responses

// InfoResponse represents an "info" response from the engine
type InfoResponse struct {
	Depth    int
	Score    int
	ScoreType string // "cp" (centipawns), "mate" (mate in X)
	PV       []string // Principal variation
	Time     int      // Search time in ms
	Nodes    int64    // Nodes searched
}

// BestMoveResponse represents a "bestmove" response
type BestMoveResponse struct {
	Move   string
	Ponder string // Ponder move (optional)
}

// ParseInfoResponse parses an "info" line from the engine
func ParseInfoResponse(line string) *InfoResponse {
	if !strings.HasPrefix(line, "info ") {
		return nil
	}
	
	parts := strings.Fields(line[5:]) // Skip "info "
	info := &InfoResponse{}
	
	for i := 0; i < len(parts); i++ {
		switch parts[i] {
		case "depth":
			if i+1 < len(parts) {
				if depth, err := strconv.Atoi(parts[i+1]); err == nil {
					info.Depth = depth
				}
				i++
			}
		case "score":
			if i+2 < len(parts) {
				info.ScoreType = parts[i+1]
				if score, err := strconv.Atoi(parts[i+2]); err == nil {
					info.Score = score
				}
				i += 2
			}
		case "time":
			if i+1 < len(parts) {
				if time, err := strconv.Atoi(parts[i+1]); err == nil {
					info.Time = time
				}
				i++
			}
		case "nodes":
			if i+1 < len(parts) {
				if nodes, err := strconv.ParseInt(parts[i+1], 10, 64); err == nil {
					info.Nodes = nodes
				}
				i++
			}
		case "pv":
			// Collect all remaining parts as the principal variation
			info.PV = parts[i+1:]
			break
		}
	}
	
	return info
}

// ParseBestMoveResponse parses a "bestmove" line from the engine
func ParseBestMoveResponse(line string) *BestMoveResponse {
	if !strings.HasPrefix(line, "bestmove ") {
		return nil
	}
	
	parts := strings.Fields(line)
	if len(parts) < 2 {
		return nil
	}
	
	response := &BestMoveResponse{
		Move: parts[1],
	}
	
	// Check for ponder move
	for i := 2; i < len(parts)-1; i++ {
		if parts[i] == "ponder" && i+1 < len(parts) {
			response.Ponder = parts[i+1]
			break
		}
	}
	
	return response
}

// IsUCIOK checks if the response is "uciok"
func IsUCIOK(line string) bool {
	return strings.TrimSpace(line) == "uciok"
}

// IsReadyOK checks if the response is "readyok"
func IsReadyOK(line string) bool {
	return strings.TrimSpace(line) == "readyok"
}