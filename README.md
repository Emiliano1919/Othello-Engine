# Othello-Engine
An Othello/Reversi AI and Engine in Go.

Currently features a Montecarlo tree search artificial intelligence, fully coded in golang from scratch.

Uses a Bitboard implementation to optimize for simulation speed and therefore number of simulations.
Uses Montecarlo Tree search to select the best move.

## Next features:

    1. ~~AI will be able to play as not only black but also white~~ (DONE)
    2. ~~Will have a nicer UI (not just terminal)~~ (DONE, but can be improved)
    3. Will have difficulty selection
    4. Will have a harder AI
    5. ~~WIll be available on itchio or something~~ (Available on Itchio https://nanuklovesfish3.itch.io/simple-othello)

## Current Ideas:
- ~~Implement a way to test 2 AIs against each other, so that they can be benchmarked~~ (DONE)
    - ~~Then once it is confirmed that the erroneous implementation is better, try to think why is it better~~ (It was not better)
- ~~Implement leaf or root parallelization~~ (DONE for root)
    - ~~Root otherwise known as Single run parallelization~~ (DONE)
        - ~~Test this version against the unparallelized~~ (DONE)
    - Implement leaf parallelization
        - Test this version against the original version and the single run parallelization
- Create a Neural Network that analyzes the current leaf to see how it will play out (maybe through self play maybe through a dataset)
- Replace UCT with other methods seen in previous research paper
    - Ask for questions to the respective creators of the evaluation functions (?)
- Get EDAX or Egaroucid running to test the game
- Maybe try to introduce the evaluation pattern used by Logistello (in some sort of way)
- At endgame, run another Algorithm instead of MCTS maybe Minimax (The depth should be small enough to get the actual best move)
- ~~When calling NextNodeFromInput we create a new node, but maybe we can take a node that already exists, if it is kept in the tree. This way we are saving the information gained from the backpropagation that has reached that node. Additionally we can cut a subtree starting from that node, that way the backpropagation algorithm does not have to run until the initial root node (the one that started the game). This would improve the amount of information we have at any time and the speed of the program~~ (DONE)
- Implement NegaScout Algorithm (?)
- Improve speed and memory allocation
    - ~~Change to smaller types where possible~~ (DONE for uint8)

## Current benchmarking results:

The 3 MCTS algorithms are running 500 simulations in total, we can see that single run parallelization performs great.
Versus was running 2 models, MCTS original vs single run parallelization, both with 500 simulations in total.

Current Benchmark results:

    goos: darwin
    goarch: arm64
    pkg: othello
    cpu: Apple M1
    BenchmarkInnacurateMonteCarloTreeSearch-8           	      30	  34751862 ns/op	 1706158 B/op	  110736 allocs/op
    BenchmarkOriginalMonteCarloTreeSearch-8             	      33	  34726694 ns/op	 1655930 B/op	  110337 allocs/op
    BenchmarkSingleRunParallelizationMCTS-8             	     136	   8713589 ns/op	 1757373 B/op	  115098 allocs/op
    BenchmarkRollout-8                                  	   16195	     73988 ns/op	    3311 B/op	     232 allocs/op
    BenchmarkRolloutParallel-8                          	   77290	     15196 ns/op	    3314 B/op	     232 allocs/op
    BenchmarkInnacurateMonteCarloTreeSearchParallel-8   	     158	   7617586 ns/op	 1744979 B/op	  113841 allocs/op
    BenchmarkInitialNodeCreationParallel-8              	10059080	       101.3 ns/op	     128 B/op	       3 allocs/op
    BenchmarkVersus-8                                   	       2	 746328500 ns/op	49786616 B/op	 3259000 allocs/op
    PASS
    ok  	othello	11.083s



## Optimizations:

### Results on memory of changing the Move type to [2]uint8 to reduce strain

Before:

        goos: darwin
        goarch: arm64
        pkg: othello
        cpu: Apple M1
        BenchmarkInnacurateMCTS-8                 	      27	  39392144 ns/op	15171193 B/op	  238426 allocs/op
        BenchmarkRollout-8                        	   12151	     96240 ns/op	   31677 B/op	     511 allocs/op
        BenchmarkRolloutParallel-8                	   48339	     27689 ns/op	   31654 B/op	     511 allocs/op
        BenchmarkInnacurateMCTS-8                  	      79	  15042758 ns/op	15511110 B/op	  245948 allocs/op
        BenchmarkInitialNodeCreationParallel-8    	 4720350	       246.2 ns/op	     280 B/op	       7 allocs/op
        BenchmarkVersus-8                         	       1	13137906167 ns/op	4164307368 B/op	67774181 allocs/op
        PASS
        ok  	othello	22.145s

After:

        goos: darwin
        goarch: arm64
        pkg: othello
        cpu: Apple M1
        BenchmarkInnacurateMCTS-8                	      32	  34703569 ns/op	 1709309 B/op	  110738 allocs/op
        BenchmarkRollout-8                        	   16146	     74191 ns/op	    3315 B/op	     232 allocs/op
        BenchmarkRolloutParallel-8                	   81868	     15022 ns/op	    3313 B/op	     232 allocs/op
        BenchmarkInnacurateMCTS-8                	     145	   7361030 ns/op	 1740413 B/op	  113845 allocs/op
        BenchmarkInitialNodeCreationParallel-8    	11684953	       101.2 ns/op	     128 B/op	       3 allocs/op
        BenchmarkVersus-8                         	       1	10628987833 ns/op	465059776 B/op	31064478 allocs/op
        PASS
        ok  	othello	17.534s

### Single Run Parallelization

Here you can see that I can fit in single run parallelization 2000 thousand rollouts in nearly the same amount of time as doing just 500 in the original MCTS.

    goos: darwin
    goarch: arm64
    pkg: othello
    cpu: Apple M1
    BenchmarkInnacurateMonteCarloTreeSearch-8           	      25	  40492423 ns/op	 4400324 B/op	  111577 allocs/op
    BenchmarkOriginalMonteCarloTreeSearch-8             	      28	  40401527 ns/op	 4346984 B/op	  111051 allocs/op
    BenchmarkSingleRunParallelizationMCTS-8             	      28	  39768070 ns/op	17519120 B/op	  457123 allocs/op
    BenchmarkRollout-8                                  	   14156	     84577 ns/op	    8694 B/op	     233 allocs/op
    BenchmarkRolloutParallel-8                          	   66181	     18163 ns/op	    8690 B/op	     233 allocs/op
    BenchmarkInnacurateMonteCarloTreeSearchParallel-8   	     126	   8891650 ns/op	 4428303 B/op	  114367 allocs/op
    BenchmarkInitialNodeCreationParallel-8              	10660340	       115.8 ns/op	     128 B/op	       3 allocs/op
    BenchmarkVersus-8                                   	       1	1018168333 ns/op	205979192 B/op	 3139643 allocs/op
    PASS
    ok  	othello	9.983s

Running a 100 simulated games, where white (the opponent) is the single run parallelization running 2000 rollouts vs 500 from original MCTS. As we can see above this is logical because root parallelization can do those 2000 simulations in the same window of time that lets the original MCTS to do just 500.

The results are:

    Opponent Wins: 53
    Draws: 21
    Total Games ran: 100
    Total run time for all the games: 2m25.277620666s%   

Now with the parallelization as black:

    Opponent Wins: 59
    Draws: 16
    Total Games ran: 100
    Total run time for all the games: 2m24.837461458s%   

### Inneficient use of random number generator


    The results of the profiler indicate that the rng generator is a performance bottle neck in both cpu and memory reworking the rng in all the mcts improved performance all around
    Extract from profiler for OriginalMonteCarloTreeSearch:

    go tool pprof mem.out
        flat  flat%   sum%        cum   cum%
    67.85MB 57.48% 57.48%    67.85MB 57.48%  math/rand.newSource
    go tool pprof cpu.out
    File: othello.test
    Type: cpu
    Time: 2025-11-17 11:36:33 EST
    Duration: 1.43s, Total samples = 1.01s (70.58%)
    Showing nodes accounting for 860ms, 85.15% of 1010ms total
    Showing top 10 nodes out of 85
        flat  flat%   sum%        cum   cum%
        510ms 50.50% 50.50%      510ms 50.50%  othello.shift (inline)
        80ms  7.92% 58.42%      410ms 40.59%  othello.generateMoves
        50ms  4.95% 63.37%       60ms  5.94%  math/rand.(*rngSource).Seed
    
    And the benchmark results:

        goos: darwin
        goarch: arm64
        pkg: othello
        cpu: Apple M1
        BenchmarkInnacurateMonteCarloTreeSearch-8           	      27	  43043940 ns/op	 4401202 B/op	  111478 allocs/op
        BenchmarkOriginalMonteCarloTreeSearch-8             	      28	  40368275 ns/op	 4346360 B/op	  111039 allocs/op
        BenchmarkSingleRunParallelizationMCTS-8             	     100	  11413553 ns/op	 4386934 B/op	  115242 allocs/op
        BenchmarkRollout-8                                  	   14127	     86011 ns/op	    8688 B/op	     233 allocs/op
        BenchmarkRolloutParallel-8                          	   61320	     18258 ns/op	    8688 B/op	     233 allocs/op
        BenchmarkInnacurateMonteCarloTreeSearchParallel-8   	     130	   8809887 ns/op	 4427154 B/op	  114317 allocs/op
        BenchmarkInitialNodeCreationParallel-8              	11673699	        99.94 ns/op	     128 B/op	       3 allocs/op
        BenchmarkVersus-8                                   	       2	 936545000 ns/op	209554416 B/op	 3285514 allocs/op
        PASS
        ok  	othello	11.180s
    
    The bytes per operation (third column of results) went down dramatically and there is some improvement in nano second per operation (second column of results). After fixing the performance issue by sharing one random number generator for sequential execution, and one rng per goroutine inside the single run parallelization.:

        goos: darwin
        goarch: arm64
        pkg: othello
        cpu: Apple M1
        BenchmarkInnacurateMonteCarloTreeSearch-8           	      30	  34751862 ns/op	 1706158 B/op	  110736 allocs/op
        BenchmarkOriginalMonteCarloTreeSearch-8             	      33	  34726694 ns/op	 1655930 B/op	  110337 allocs/op
        BenchmarkSingleRunParallelizationMCTS-8             	     136	   8713589 ns/op	 1757373 B/op	  115098 allocs/op
        BenchmarkRollout-8                                  	   16195	     73988 ns/op	    3311 B/op	     232 allocs/op
        BenchmarkRolloutParallel-8                          	   77290	     15196 ns/op	    3314 B/op	     232 allocs/op
        BenchmarkInnacurateMonteCarloTreeSearchParallel-8   	     158	   7617586 ns/op	 1744979 B/op	  113841 allocs/op
        BenchmarkInitialNodeCreationParallel-8              	10059080	       101.3 ns/op	     128 B/op	       3 allocs/op
        BenchmarkVersus-8                                   	       2	 746328500 ns/op	49786616 B/op	 3259000 allocs/op
        PASS
        ok  	othello	11.083s

    Extract from profiler after improvemnt for OriginalMonteCarloTreeSearch:
    go tool pprof mem.out
    Showing nodes accounting for 53.28MB, 96.37% of 55.28MB total
    Showing top 10 nodes out of 66
        flat  flat%   sum%        cum   cum%
    28.50MB 51.56% 51.56%    28.50MB 51.56%  othello.ArrayOfPositionalMoves (inline)
        17MB 30.75% 82.31%       17MB 30.75%  othello.ArrayOfMoves (inline)
        1.50MB  2.72% 85.03%     1.50MB  2.72%  runtime.allocm
    go tool pprof cpu.out                                                                        
    File: othello.test
    Type: cpu
    Time: 2025-11-17 19:36:42 EST
    Duration: 1.31s, Total samples = 1.07s (81.54%)
    Entering interactive mode (type "help" for commands, "o" for options)
    (pprof) top
    Showing nodes accounting for 1.01s, 94.39% of 1.07s total
    Showing top 10 nodes out of 74
        flat  flat%   sum%        cum   cum%
        0.52s 48.60% 48.60%      0.52s 48.60%  othello.shift (inline)
        0.25s 23.36% 71.96%      0.25s 23.36%  runtime.madvise
        0.08s  7.48% 79.44%      0.09s  8.41%  othello.ArrayOfMoves
        0.04s  3.74% 83.18%      0.41s 38.32%  othello.generateMoves
        0.04s  3.74% 86.92%      0.04s  3.74%  runtime.memclrNoHeapPointers
        0.03s  2.80% 89.72%      0.05s  4.67%  othello.(*Node).Expand
        0.02s  1.87% 91.59%      0.02s  1.87%  runtime.scanobject
        0.01s  0.93% 92.52%      0.01s  0.93%  internal/runtime/atomic.(*Uint32).Add
        0.01s  0.93% 93.46%      0.01s  0.93%  math/rand.(*rngSource).Uint64
        0.01s  0.93% 94.39%      0.06s  5.61%  othello.ArrayOfPositionalMoves
