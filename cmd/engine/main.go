package main

import (
	"github.com/jgerontis/go-chess/internal/engine"
)

func main() {
	// Start the UCI engine - this will handle stdin/stdout UCI protocol
	engine.StartEngine()
}