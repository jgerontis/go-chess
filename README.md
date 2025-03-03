# GO Chess

## Why?

I love the game of chess, and while I don't play very often, I still usually do my daily puzzles.
I also have wanted to do more programming for fun, and I thought this would be a good project to work on.

## How?

For now, you can run it with `go run ./cmd/main.go {optional FEN string goes here}`
I'll look into adding more features later.

## Progress report:

Done:

- Board representation
- FEN parsing (haven't done castling and turn yet)
- Basic piece moving
  To do:
- Turn order
- Legal move calculation
- Restricting moves based on legal moves
- Checks and checkmates
- Castling
- En passant
- Pawn promotion
- Stalemate
- Draw by repetition
- Draw by insufficient material
- Draw by 50-move rule
- Draw by 3-fold repetition
- Full FEN parsing
- Full FEN generation

## Making the chess AI:

(I let Co-Pilot fill this section out for me)

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
- UCI protocol
