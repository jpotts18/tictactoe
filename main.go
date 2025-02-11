package main

import (
	"flag"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/jpotts18/tictactoe/agent"
	"github.com/jpotts18/tictactoe/game"
)

// RewardScheme holds the reward values for different game outcomes
type RewardScheme struct {
	win   float64
	draw  float64
	loss  float64
	step  float64
}

func evaluateAgents(agent1 agent.Agent, agent2 agent.Agent, numGames int, rewards RewardScheme) (wins, draws, losses int) {
	for i := 0; i < numGames; i++ {
		board := game.NewBoard()
		
		// Randomly decide who goes first
		agent1GoesFirst := rand.Float64() < 0.5
		var firstAgent, secondAgent agent.Agent
		var firstPlayer, secondPlayer int
		
		if agent1GoesFirst {
			firstAgent, secondAgent = agent1, agent2
			firstPlayer, secondPlayer = 1, 2
		} else {
			firstAgent, secondAgent = agent2, agent1
			firstPlayer, secondPlayer = 2, 1
		}

		for {
			// First agent's turn
			oldState := firstAgent.GetStateKey(board)
			move := firstAgent.ChooseAction(board)
			board.MakeMove(move, firstPlayer)
			
			gameOver, winner := board.IsGameOver()
			if gameOver {
				var reward float64
				if agent1GoesFirst {
					if winner == 1 {
						wins++
						reward = rewards.win
					} else if winner == 0 {
						draws++
						reward = rewards.draw
					} else {
						losses++
						reward = rewards.loss
					}
				} else {
					if winner == 2 {
						losses++
						reward = rewards.loss
					} else if winner == 0 {
						draws++
						reward = rewards.draw
					} else {
						wins++
						reward = rewards.win
					}
				}
				firstAgent.Learn(oldState, move, reward, firstAgent.GetStateKey(board))
				break
			} else {
				firstAgent.Learn(oldState, move, rewards.step, firstAgent.GetStateKey(board))
			}
			
			// Second agent's turn
			oldState = secondAgent.GetStateKey(board)
			move = secondAgent.ChooseAction(board)
			board.MakeMove(move, secondPlayer)
			
			gameOver, winner = board.IsGameOver()
			if gameOver {
				var reward float64
				if agent1GoesFirst {
					if winner == 1 {
						wins++
						reward = -rewards.win
					} else if winner == 0 {
						draws++
						reward = rewards.draw
					} else {
						losses++
						reward = -rewards.loss
					}
				} else {
					if winner == 2 {
						losses++
						reward = -rewards.win
					} else if winner == 0 {
						draws++
						reward = rewards.draw
					} else {
						wins++
						reward = -rewards.loss
					}
				}
				secondAgent.Learn(oldState, move, reward, secondAgent.GetStateKey(board))
				break
			} else {
				secondAgent.Learn(oldState, move, rewards.step, secondAgent.GetStateKey(board))
			}
		}
	}
	return
}

// Add this helper function for consistent output format
func printResults(name string, wins, draws, losses, numGames int) {
	fmt.Printf("Results after %d games:\n", numGames)
	fmt.Printf("Agent1 - Wins: %d, Draws: %d, Losses: %d\n", wins, draws, losses)
	fmt.Printf("Win rate: %.2f%%, Draw rate: %.2f%%\n\n", 
		float64(wins)/float64(numGames)*100,
		float64(draws)/float64(numGames)*100)
}

// Add this helper function to evaluate Q-Learning progress
func evaluateProgress(qagent agent.Agent, opponent agent.Agent, numGames int) (winRate, drawRate float64) {
	wins, draws, _ := evaluateAgents(qagent, opponent, numGames, RewardScheme{})
	return float64(wins)/float64(numGames)*100, float64(draws)/float64(numGames)*100
}

func main() {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	// Define command-line flags
	trainCmd := flag.Bool("train", false, "Train all models")
	evalCmd := flag.Bool("evaluate", false, "Evaluate trained models")
	playCmd := flag.Bool("play", false, "Play against a trained model")
	flag.Parse()

	if !*trainCmd && !*evalCmd && !*playCmd {
		fmt.Println("Please specify one of the following commands:")
		fmt.Println("  -train     Train all models")
		fmt.Println("  -evaluate  Evaluate trained models")
		fmt.Println("  -play      Play against a trained model")
		return
	}

	if *trainCmd {
		trainModels()
	}
	if *evalCmd {
		evaluateModels()
	}
	if *playCmd {
		playGame()
	}
}

func trainModels() {
	fmt.Println("=== Training Models ===")
	
	// Create learning agents
	qagent := agent.NewQAgent(1)
	sarsaAgent := agent.NewSarsaAgent(1)
	mcAgent := agent.NewMonteCarloAgent(1)
	
	// Training configurations
	iterations := 100000
	evalFrequency := 5000
	benchmarkRandom := agent.NewRandomAgent(2)

	// Train each agent
	trainAgent("Q-Learning", qagent, iterations, evalFrequency, benchmarkRandom)
	trainAgent("SARSA", sarsaAgent, iterations, evalFrequency, benchmarkRandom)
	trainAgent("Monte Carlo", mcAgent, iterations, evalFrequency, benchmarkRandom)

	// Save only learning agents
	if learningAgent, ok := qagent(agent.LearningAgent); ok {
		learningAgent.Save("models/qagent")
	}
	if learningAgent, ok := sarsaAgent.(agent.LearningAgent); ok {
		learningAgent.Save("models/sarsa")
	}
	if learningAgent, ok := mcAgent.(agent.LearningAgent); ok {
		learningAgent.Save("models/montecarlo")
	}
}

func trainAgent(name string, trainAgent agent.Agent, iterations, evalFrequency int, benchmark agent.Agent) {
	fmt.Printf("Training %s agent...\n", name)
	
	for i := 0; i < iterations; i++ {
		// Self-play training
		evaluateAgents(trainAgent, trainAgent, 1, RewardScheme{
			win: 1.0, draw: 0.5, loss: -2.0, step: -0.01,
		})

		// Periodic evaluation
		if (i+1) % evalFrequency == 0 {
			wins, draws, losses := evaluateAgents(trainAgent, benchmark, 100, RewardScheme{})
			
			// Create progress bar (30 chars wide)
			winChars := int(float64(wins) * 0.3)    // Scale to 30 chars total
			drawChars := int(float64(draws) * 0.3)
			lossChars := 30 - winChars - drawChars

			fmt.Printf("Iteration %d: [", i+1)
			// Print wins in green
			fmt.Printf("\033[32m%s", strings.Repeat("█", winChars))
			// Print draws in yellow
			fmt.Printf("\033[33m%s", strings.Repeat("█", drawChars))
			// Print losses in red
			fmt.Printf("\033[31m%s\033[0m", strings.Repeat("█", lossChars))
			fmt.Printf("] W:%d%% D:%d%% L:%d%%\n",
				wins, draws, losses)
		}
	}
	fmt.Println()
}

func evaluateModels() {
	fmt.Println("=== Evaluating Models ===")
	
	// Create agents
	qagent := agent.NewQAgent(1)
	sarsaAgent := agent.NewSarsaAgent(1)
	mcAgent := agent.NewMonteCarloAgent(1)

	// Load trained models
	if err := qagent.Load("models/qagent"); err != nil {
		fmt.Println("No trained Q-Learning model found")
	}
	if err := sarsaAgent.Load("models/sarsa"); err != nil {
		fmt.Println("No trained SARSA model found")
	}
	if err := mcAgent.Load("models/montecarlo"); err != nil {
		fmt.Println("No trained Monte Carlo model found")
	}

	// Number of games for evaluation
	numGames := 1000

	// Create benchmark agents
	random := agent.NewRandomAgent(2)
	minimax := agent.NewMinimaxAgent(2)

	// Evaluate each agent against random and minimax
	evaluateAgent("Q-Learning", qagent, random, minimax, numGames)
	evaluateAgent("SARSA", sarsaAgent, random, minimax, numGames)
	evaluateAgent("Monte Carlo", mcAgent, random, minimax, numGames)
}

func evaluateAgent(name string, testAgent, random, minimax agent.Agent, numGames int) {
	fmt.Printf("\nEvaluating %s agent:\n", name)
	
	// vs Random
	wins, draws, losses := evaluateAgents(testAgent, random, numGames, RewardScheme{})
	fmt.Printf("vs Random:  Win = %.1f%%, Draw = %.1f%%, Loss = %.1f%%\n",
		float64(wins)/float64(numGames)*100,
		float64(draws)/float64(numGames)*100,
		float64(losses)/float64(numGames)*100)

	// vs Minimax
	wins, draws, losses = evaluateAgents(testAgent, minimax, numGames, RewardScheme{})
	fmt.Printf("vs Minimax: Win = %.1f%%, Draw = %.1f%%, Loss = %.1f%%\n",
		float64(wins)/float64(numGames)*100,
		float64(draws)/float64(numGames)*100,
		float64(losses)/float64(numGames)*100)
}

func playGame() {
	fmt.Println("=== Play Against AI ===")
	
	// Create agents
	agents := map[string]agent.Agent{
		"Random":      agent.NewRandomAgent(1),
		"Minimax":     agent.NewMinimaxAgent(1),
		"Q-Learning":  agent.NewQAgent(1),
		"SARSA":       agent.NewSarsaAgent(1),
		"Monte Carlo": agent.NewMonteCarloAgent(1),
	}

	for {
		fmt.Println("\nChoose your opponent:")
		options := make([]string, 0, len(agents))
		for name := range agents {
			options = append(options, name)
		}
		for i, name := range options {
			fmt.Printf("%d. %s\n", i+1, name)
		}
		fmt.Printf("%d. Exit\n", len(options)+1)

		var choice int
		fmt.Print("Enter your choice: ")
		fmt.Scan(&choice)

		if choice == len(options)+1 {
			fmt.Println("Thanks for playing!")
			break
		}

		if choice < 1 || choice > len(options) {
			fmt.Println("Invalid choice, please try again")
			continue
		}

		agentName := options[choice-1]
		fmt.Printf("\nPlaying against %s agent\n", agentName)
		playAgainstAgent(agents[agentName])

		fmt.Print("\nPlay another game? (y/n): ")
		var response string
		fmt.Scan(&response)
		if response != "y" {
			fmt.Println("Thanks for playing!")
			break
		}
	}
}

func playAgainstAgent(agent agent.Agent) {
	board := game.NewBoard()
	for {
		// Agent's turn
		move := agent.ChooseAction(board)
		board.MakeMove(move, 1)
		fmt.Printf("\nAgent plays position %d:\n%s\n", move+1, board.String())
		
		gameOver, winner := board.IsGameOver()
		if gameOver {
			if winner == 1 {
				fmt.Println("Agent wins!")
			} else if winner == 0 {
				fmt.Println("It's a draw!")
			} else {
				fmt.Println("You win!")
			}
			break
		}
		
		// Human player's turn
		var humanMove int
		for {
			fmt.Print("Enter your move (1-9): ")
			fmt.Scan(&humanMove)
			if humanMove >= 1 && humanMove <= 9 && board.MakeMove(humanMove-1, 2) {
				break
			}
			fmt.Println("Invalid move, try again")
		}
		
		fmt.Printf("\nYour move:\n%s\n", board.String())
		
		gameOver, winner = board.IsGameOver()
		if gameOver {
			if winner == 1 {
				fmt.Println("Agent wins!")
			} else if winner == 0 {
				fmt.Println("It's a draw!")
			} else {
				fmt.Println("You win!")
			}
			break
		}
	}
}

// Add a function to compare different reward schemes
func compareRewardSchemes() {
	schemes := []RewardScheme{
		{win: 1.0, draw: 0.0, loss: -1.0, step: 0.0},    // Standard
		{win: 1.0, draw: 0.5, loss: -1.0, step: 0.0},    // Reward draws
		{win: 2.0, draw: 0.0, loss: -1.0, step: 0.0},    // Emphasize winning
		{win: 1.0, draw: 0.0, loss: -2.0, step: 0.0},    // Emphasize avoiding losses
		{win: 1.0, draw: 0.0, loss: -1.0, step: -0.1},   // Penalize long games
	}

	for i, scheme := range schemes {
		fmt.Printf("\n=== Testing Reward Scheme %d ===\n", i+1)
		fmt.Printf("Win: %.1f, Draw: %.1f, Loss: %.1f, Step: %.1f\n", 
			scheme.win, scheme.draw, scheme.loss, scheme.step)

		// Create fresh agents for each scheme
		qagent := agent.NewQAgent(1)
		sarsaAgent := agent.NewSarsaAgent(1)
		mcAgent := agent.NewMonteCarloAgent(1)
		randomAgent := agent.NewRandomAgent(2)

		// Train agents
		for j := 0; j < 5; j++ {
			// Training phase
			evaluateAgents(qagent, qagent, 1000, scheme)
			evaluateAgents(sarsaAgent, sarsaAgent, 1000, scheme)
			evaluateAgents(mcAgent, mcAgent, 1000, scheme)

			// Evaluation phase against random agent
			fmt.Printf("\nAfter %d iterations:\n", (j+1)*1000)
			
			qagent.Evaluating = true
			sarsaAgent.Evaluating = true
			mcAgent.Evaluating = true
			
			wins, draws, _ := evaluateAgents(qagent, randomAgent, 200, scheme)
			fmt.Printf("Q-Learning - Win: %.1f%%, Draw: %.1f%%\n", 
				float64(wins)/2.0, float64(draws)/2.0)
			
			wins, draws, _ = evaluateAgents(sarsaAgent, randomAgent, 200, scheme)
			fmt.Printf("SARSA      - Win: %.1f%%, Draw: %.1f%%\n", 
				float64(wins)/2.0, float64(draws)/2.0)
			
			wins, draws, _ = evaluateAgents(mcAgent, randomAgent, 200, scheme)
			fmt.Printf("Monte Carlo - Win: %.1f%%, Draw: %.1f%%\n", 
				float64(wins)/2.0, float64(draws)/2.0)
			
			qagent.Evaluating = false
			sarsaAgent.Evaluating = false
			mcAgent.Evaluating = false
		}
	}
} 
