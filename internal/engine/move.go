package engine

type Move uint16

// Move is a 16-bit integer with the following format:
// ffff tttttt ssssss
// s bits 0-5 = source square (0-63)
// t bits 6-11 = target square (0-63)
// f bits 12-15 = flags
