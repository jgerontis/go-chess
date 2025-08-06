package uci

// Example usage of the UCI package for reference

/*
Example: Basic UCI interaction

	// Connect to a local Stockfish engine
	client, err := uci.NewClientExec("stockfish")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// Initialize UCI communication
	err = client.InitializeEngine()
	if err != nil {
		log.Fatal(err)
	}

	// Wait for uciok response
	for {
		response, err := client.ReadResponse()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Engine:", response)
		if uci.IsUCIOK(response) {
			break
		}
	}

	// Check if engine is ready
	err = client.IsReady()
	if err != nil {
		log.Fatal(err)
	}

	// Wait for readyok
	for {
		response, err := client.ReadResponse()
		if err != nil {
			log.Fatal(err)
		}
		if uci.IsReadyOK(response) {
			break
		}
	}

	// Set starting position
	err = client.SetPosition("startpos")
	if err != nil {
		log.Fatal(err)
	}

	// Start search
	err = client.GoDepth(10)
	if err != nil {
		log.Fatal(err)
	}

	// Read responses until we get bestmove
	for {
		response, err := client.ReadResponse()
		if err != nil {
			log.Fatal(err)
		}

		// Parse different response types
		if info := uci.ParseInfoResponse(response); info != nil {
			fmt.Printf("Depth: %d, Score: %d %s, PV: %v\n", 
				info.Depth, info.Score, info.ScoreType, info.PV)
		} else if bestMove := uci.ParseBestMoveResponse(response); bestMove != nil {
			fmt.Printf("Best move: %s\n", bestMove.Move)
			break
		}
	}

Example: Async communication for GUI

	client, err := uci.NewClientExec("stockfish")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// Setup async communication
	responseChan := make(chan string, 100)
	errorChan := make(chan error, 10)
	client.StartResponseListener(responseChan, errorChan)

	// Handle responses in GUI
	go func() {
		for {
			select {
			case response := <-responseChan:
				// Update GUI with engine response
				if bestMove := uci.ParseBestMoveResponse(response); bestMove != nil {
					// Make the move in the game
					game.MakeMove(bestMove.Move)
				}
			case err := <-errorChan:
				// Handle error in GUI
				log.Printf("Engine error: %v", err)
				return
			}
		}
	}()

	// Send commands from GUI
	client.InitializeEngine()
	client.SetPosition("startpos")
	client.GoDepth(10)
*/