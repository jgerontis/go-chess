# Go Chess

## Why?

I love the game of chess, and while I don't play very often, I still usually do my daily puzzles.
I also have wanted to do more programming for fun, and I thought this would be a good project to work on.
Also, I haven't worked as much in Go lately as I would like, so I thought this would be a good opportunity to get back into it.

## How?

For now, you can run it with `go run ./cmd/debug {optional FEN string goes here}`

## Structure:
This project is more or less split into three parts:
1. The GUI that that will be used to play the game
2. The internal chess logic to be used by the engine and the GUI
3. The chess engine that will be used to play against the player


## GUI:
### Done:
- Basic board rendering
- Basic piece rendering
- Basic piece movement
- Basic piece selection
- Basic piece deselection
- Basic piece capture

### Todo:
- Piece promotion
- Castling
- En passant
- Check and checkmate detection
- Stalemate detection
- Draw by repetition detection
- Draw by insufficient material detection
- Draw by 50-move rule detection
- Sound effects
- UCI protocol
- Animations (maybe)
- Themes (maybe)
- Multiplayer (maybe)

## Chess logic:
### Done:
- Board representation
- Piece representation
- Move representation
- Bitboards
- Basic FEN parsing to load a board state
- Full FEN parsing
- Full FEN generation
- Basic piece moving
- Turn order
- Pseudo-legal move calculation with Bitboards 🙌🙌
- Psuedo-legal move calculation with magic bitboards 🧙‍♂️


#### To do:
- Checks and checkmates
- Pins
- Castling
- En passant
- Pawn promotion
- Stalemate
- Draw by repetition
- Draw by insufficient material
- Draw by 50-move rule
- Draw by 3-fold repetition


## Making the chess AI:

Todo:
- Board evaluation
- Minimax algorithm
- Alpha-beta pruning
- Iterative deepening
- Transposition tables
- Quiescence search
- Move ordering
- Null move pruning
- Late move reduction
- Killer move heuristic
- History heuristic
- Futility pruning
- Razoring
- Static exchange evaluation
- Principal variation search
- Aspiration windows
- Time management
- Multi-threading
- Opening book
- Endgame tablebases

## Stretch goals:

- PGN parsing
- PGN generation
