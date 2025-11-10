# Othello-Engine
An Othello/Reversi AI and Engine in Go.

Currently features a Montecarlo tree search artificial intelligence, fully coded in golang from scratch.

Uses a Bitboard implementation to optimize for simulation speed and therefore number of simulations.
Uses Montecarlo Tree search to select the best move.

Next features:

    1. ~~AI will be able to play as not only black but also white~~ (DONE)
    2. ~~Will have a nicer UI (not just terminal)~~ (DONE, but can be improved)
    3. Will have difficulty selection
    4. Will have a harder AI
    5. ~~WIll be available on itchio or something~~ (Available on Itchio https://nanuklovesfish3.itch.io/simple-othello)

Notes:
    It seems like whoever is black has an advantage on the game, it seems like the AI loses more easily if it is white rather than black. 

        Possible fixes:

        - Add more compute in case that AI is white
        - Modify the C parameter?

    Actual root cause: The backpropagation algorithm is not taking into account if the machine is black or white, defaults to black. Therefore it is easier because technically the AI currently is choosing the moves that are more likely to make black win, even when the machine is white. The machine is acting against its own interests.

    The problem has now been fixed. But it seems like the AI plays a kind of different strategy (I know it is just simulation and probability) when it plays white. It plays a long game, in my opinion.

Current Ideas:
- ~~Implement a way to test 2 AIs against each other, so that they can be benchmarked~~ (DONE)
    - Then once it is confirmed that the erroneous implementation is better, try to think why is it better
- Implement leaf or root parallelization (Hopefully this will be straightforward)
- Create a Neural Network that analyzes the current leaf to see how it will play out (maybe through self play maybe through a dataset)
- Replace UCT with other methods seen in previous research paper
    - Ask for questions to the respective creators of the evaluation functions (?)
- Get EDAX or Egaroucid running to test the game
- Maybe try to introduce the evaluation pattern used by Logistello (in some sort of way)
- At endgame, run another Algorithm instead of MCTS maybe Minimax (The depth should be small enough to get the actual best move)
- When calling NextNodeFromInput we create a new node, but maybe we can take a node that already exists, if it is kept in the tree. This way we are saving the information gained from the backpropagation that has reached that node. Additionally we can cut a subtree starting from that node, that way the backpropagation algorithm does not have to run until the initial root node (the one that started the game). This would improve the amount of information we have at any time and the speed of the program
- Implement NegaScout Algorithm (?)

Current Results:

With my double expansion greedy MCTS as the opponent (black) out of 100 games this was the result against the Original MCTS with UCT
 
    Opponent Wins: 61
    Draws: 4
    Total Games ran: 100

So it seems it is actually better than the normal MCTS, but the output for my algorithm as the opponent (white) was

    Opponent Wins: 31
    Draws: 4
    Total Games ran: 100
    Total run time for all the games: 19m0.365755959s% 

So it is worse in this case. I actually suspect this is a bug the ratio looks flipped. I will investigate. 

It was a bug, I was not awarding the wins to the nodes correctly. Only nodes that have a parent with white turn should receive wins when we have wins for white, same story for black. (I will write more notes about this bug later). This bug seems to benefit Black, which is kind of interesting without the bug i get 41 wins and 14 draws out of a 100 games.
Also it seems that the AI benefits from a more slow approach when it is white (probably because it is the second one to move) removing the double expansion and fixing the bug in both my version and the original one. 
    Opponent Wins: 7
    Draws: 0
    Total Games ran: 10 
    Total run time for all the games: 2m2.1476085s% 

This simulations take some time to run. Parallelization now seems like a necessary improvement to run the code at a faster speed. The code also needs some improvements, userIsBlack is confusing to use as a variable when running 2 algorithms.
Note: Both algorithms were doing 5000 rollouts at each leaf Node.

Current Benchmark results:
        goos: darwin
        goarch: arm64
        pkg: othello
        cpu: Apple M1
        BenchmarkMonteCarloTreeSearc-8            	      32	  34703569 ns/op	 1709309 B/op	  110738 allocs/op
        BenchmarkRollout-8                        	   16146	     74191 ns/op	    3315 B/op	     232 allocs/op
        BenchmarkRolloutParallel-8                	   81868	     15022 ns/op	    3313 B/op	     232 allocs/op
        BenchmarkMonteCarloTreeSearchParallel-8   	     145	   7361030 ns/op	 1740413 B/op	  113845 allocs/op
        BenchmarkInitialNodeCreationParallel-8    	11684953	       101.2 ns/op	     128 B/op	       3 allocs/op
        BenchmarkVersus-8                         	       1	10628987833 ns/op	465059776 B/op	31064478 allocs/op
        PASS
        ok  	othello	17.534s
