# ♟️ Go Chess

A high-performance chess engine and GUI implementation written in Go, featuring bitboard-based representation and magic bitboards for fast move generation.

## ✨ Features

### 🎯 **Core Chess Engine**
- **Magic Bitboard Move Generation** - Fast sliding piece move calculation
- **Legal Move Validation** - Check/pin detection and filtering
- **Special Moves** - Castling, en passant, and pawn promotion
- **FEN Support** - Position parsing and generation
- **UCI Protocol** - Standard engine communication

### 🎮 **Interactive GUI**
- **Ebiten Graphics** - 2D rendering with SVG pieces
- **Multiple Game Modes** - Human vs AI, AI vs AI, debug mode
- **Audio Support** - Sound effects for game events

## 🚀 Quick Start

### Prerequisites
- Go 1.23+ installed
- Git for cloning the repository

### Installation & Running

```bash
# Clone the repository
git clone https://github.com/jgerontis/go-chess.git
cd go-chess

# Run the main GUI (starts with menu)
go run ./cmd/main

# Or skip directly to specific modes
go run ./cmd/main --debug           # Manual piece movement
go run ./cmd/main --human-vs-ai     # Play against AI
go run ./cmd/main --ai-vs-ai        # Watch AI vs AI

# UCI Client for engine communication  
go run ./cmd/engine-cli tcp localhost:8080    # Connect to TCP server
go run ./cmd/engine-cli exec stockfish        # Launch Stockfish
go run ./cmd/engine-cli exec ./my-engine      # Launch your custom engine

# GoChess Engine (standalone UCI engine)
go build -o my-engine ./cmd/engine            # Build your engine
go run ./cmd/engine-cli exec ./my-engine      # Test your engine

# Run with custom position
go run ./cmd/main --fen "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1"
# Legacy format still supported
go run ./cmd/main "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1"
```

### Development Commands

```bash
# Run all tests
go test ./...

# Run specific package tests
go test ./internal/chess -v

# Format code
go fmt ./...

# Build all binaries
go build ./cmd/...
```

## 🎯 Current Status

### ✅ **Completed**
- Move generation with magic bitboards
- Legal move filtering and check detection
- Special moves (en passant, castling, promotion)
- FEN parsing and board representation
- GUI with main menu and game mode selection
- Command line argument parsing and navigation
- UCI client library with TCP and process communication
- Shared UCI layer for engine integration

### 🚧 **In Progress**
- Checkmate and stalemate detection
- Draw condition detection
- Enhanced GUI features

### 📋 **Roadmap**

#### **Phase 1: Complete Game Rules** 
- [ ] Checkmate/stalemate detection
- [ ] Draw condition detection
- [ ] Game state management

#### **Phase 2: Engine Intelligence**
- [ ] Position evaluation function
- [ ] Minimax with alpha-beta pruning
- [ ] Iterative deepening search
- [ ] Transposition tables

#### **Phase 3: Advanced Features**
- [ ] Opening book integration
- [ ] Endgame tablebase support
- [ ] Time management
- [ ] Advanced search techniques

#### **Phase 4: Polish & Features**
- [ ] PGN import/export
- [ ] Multiple GUI themes
- [ ] Network multiplayer
- [ ] Analysis tools

## 🏛️ Architecture Overview

```
├── cmd/                    # Executable entry points
│   ├── main/              # Main GUI with menu system and all game modes
│   │   └── main.go        # GUI application entry point
│   ├── engine-cli/        # UCI client (works with any UCI engine: Stockfish, etc.)
│   │   └── main.go        # Interactive UCI client
│   └── engine/            # Standalone GoChess UCI engine
│       └── main.go        # Engine executable entry point
├── internal/
│   ├── chess/             # Core chess logic
│   │   ├── bitboard.go    # Bitboard operations & magic bitboards
│   │   ├── board.go       # Board state & move execution
│   │   ├── movegen.go     # Move generation & legal filtering
│   │   ├── move.go        # Move representation and execution
│   │   ├── move_consts.go # Move flags and constants
│   │   ├── fen.go         # FEN parsing/generation
│   │   ├── piece.go       # Piece representation
│   │   └── *_test.go      # Comprehensive test suite
│   ├── engine/            # Chess engine implementation
│   │   ├── engine.go      # Main engine with async search
│   │   └── search.go      # Search algorithms
│   └── uci/               # Complete UCI communication layer
│       ├── client.go      # UCI client (TCP and process)
│       ├── server.go      # UCI server infrastructure
│       ├── commands.go    # UCI command helpers
│       ├── responses.go   # Response parsing utilities
│       └── example_usage.go # Usage documentation
├── gui/                   # Ebiten-based GUI with menu system
│   ├── app.go            # Application state management
│   ├── menu.go           # Main menu interface
│   ├── game.go           # Chess game interaction
│   ├── ui.go             # Board rendering and graphics
│   └── audio.go          # Sound effects
└── assets/               # Graphics and audio resources
```

## 🤝 Why This Project?

This project combines several passions:
- **♟️ Chess** - I love chess, and while my rating isn't that high, I find the game fascinating.
- **🚀 Go Programming** - Exploring Go's performance and elegance  
- **🔬 Algorithms** - Implementing classic CS algorithms in a real application
- **🎨 Graphics Programming** - Creating an intuitive visual interface

It serves as both a learning experience in advanced programming concepts and a fully functional chess application that can keep up with commercial engines.

## 🛠️ Technical Highlights

- **Magic Bitboards**: Fast sliding piece move generation
- **Efficient Architecture**: Clean separation between engine, GUI, and chess logic
- **Comprehensive Testing**: Ensures correctness of chess rules

## 📄 License

This project is open source and available under the MIT License.

## 🙏 Acknowledgments

A ton of credit has to go out to [Sebastian Lague](https://www.youtube.com/@SebastianLague) for his excellent YouTube series *Coding Adventures*, which was the inspiration for this project.
His videos on chess programming and bitboards were invaluable in getting started.

The project also draws from various open source chess engines and libraries, which provided insights into move generation and game logic.
See:
- [Stockfish](https://stockfishchess.org/)
- [Chess Programming Wiki](https://chessprogramming.wikispaces.com/)
