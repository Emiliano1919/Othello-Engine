# Othello-Engine
An Othello/Reversi AI and Engine in Go.

Currently features a Montecarlo tree search artificial intelligence, fully coded in golang from scratch.

Uses a Bitboard implementation to optimize for simulation speed and therefore number of simulations.
Uses Montecarlo Tree search to select the best move.

## Next features:

    1. AI will be able to play as not only black but also white (DONE)
    2. Will have a nicer UI (not just terminal) (DONE, but can be improved)
    3. Will have difficulty selection
    4. Will have a harder AI
    5. WIll be available on itchio or something (Available on Itchio https://nanuklovesfish3.itch.io/simple-othello)

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
    - ~~Implement PUCT~~ (DONE)
    - Ask for questions to the respective creators of the evaluation functions (?)
- Get EDAX or Egaroucid running to test the game
- Maybe try to introduce the evaluation pattern used by Logistello (in some sort of way)
- At endgame, run another Algorithm instead of MCTS maybe Minimax (The depth should be small enough to get the actual best move)
- ~~When calling NextNodeFromInput we create a new node, but maybe we can take a node that already exists, if it is kept in the tree. This way we are saving the information gained from the backpropagation that has reached that node. Additionally we can cut a subtree starting from that node, that way the backpropagation algorithm does not have to run until the initial root node (the one that started the game). This would improve the amount of information we have at any time and the speed of the program~~ (DONE)
- Add a way to simulate based on time rather than simulation count
- Add a way to simulate while the opponent makes its move
- Implement parent Q initialization
- Virtual loss for the parallelization (?)
- Implement NegaScout Algorithm (?)
- Improve speed and memory allocation
    - ~~Change to smaller types where possible~~ (DONE for uint8)
    - ~~Start using profiler and benchmarking to test~~ (DONE) 
    - ~~Optimize the random number generator~~ (DONE)
    - ~~Optimize the memory expensive ArrayOfPositionalMoves, and ArrayOfMoves~~ (DONE)
- Improve code readability
    - CurrentScore can be a gamestate function instead of a node function
    - Try to avoid repetition
    - Improve comments

## Current benchmarking results:

The 3 MCTS algorithms are running 500 simulations in total, we can see that single run parallelization performs great.
Versus was running 2 models, MCTS original vs single run parallelization, both with 500 simulations in total.

Current Benchmark results:

    goos: darwin
    goarch: arm64
    pkg: othello
    cpu: Apple M1
    BenchmarkInnacurateMonteCarloTreeSearch-8           	      36	  30143149 ns/op	  174613 B/op	    2867 allocs/op
    BenchmarkOriginalMonteCarloTreeSearch-8             	      38	  30215717 ns/op	   88595 B/op	    1395 allocs/op
    BenchmarkMonteCarloTreeSearchPUCT-8                 	      38	  30551758 ns/op	  301920 B/op	    3874 allocs/op
    BenchmarkSingleRunParallelizationMCTS-8             	     169	   7022824 ns/op	  141911 B/op	    1518 allocs/op
    BenchmarkRollout-8                                  	   18679	     64245 ns/op	       0 B/op	       0 allocs/op
    BenchmarkRolloutParallel-8                          	   92790	     11574 ns/op	       0 B/op	       0 allocs/op
    BenchmarkInnacurateMonteCarloTreeSearchParallel-8   	     182	   5911194 ns/op	  179843 B/op	    2909 allocs/op
    BenchmarkVersus-8                                   	       2	 593293375 ns/op	 6632804 B/op	   79010 allocs/op
    PASS
    ok  	othello	9.883s



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


The results of the profiler indicate that the rng generator is a performance bottle neck in both cpu and memory reworking the rng in all the mcts improved performance all around.
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

### Inneficient way of representing moves

My previous implementation used an array of []uint8 for each node, and i also used to convert this into a positional [2]uint8 representation  (row and column). Chainging this to a simple representation and calculating the row and column parts on the fly using bit manipulation tricks (modulo using & and division using shift) reduced massively the use of memory and the speed of execution.

Before:
    BenchmarkSingleRunParallelizationMCTS-8   	     136	   8759955 ns/op	 1757296 B/op	  115100 allocs/op
    PASS
    ok  	othello	1.934s
                                                                        
    File: othello.test
    Type: alloc_space
    Time: 2025-11-18 17:07:16 EST
    Entering interactive mode (type "help" for commands, "o" for options)
    (pprof) top
    Showing nodes accounting for 236.10MB, 98.04% of 240.83MB total
    Dropped 24 nodes (cum <= 1.20MB)
    Showing top 10 nodes out of 35
        flat  flat%   sum%        cum   cum%
        126MB 52.32% 52.32%      126MB 52.32%  othello.ArrayOfPositionalMoves (inline)
    89.50MB 37.16% 89.49%    89.50MB 37.16%  othello.ArrayOfMoves (inline)
        8.04MB  3.34% 92.82%     8.04MB  3.34%  math/rand.newSource
        6.50MB  2.70% 95.52%       10MB  4.15%  othello.NewNode
        2MB  0.83% 96.36%        2MB  0.83%  runtime.allocm
        1.50MB  0.62% 96.98%     2.50MB  1.04%  github.com/hajimehoshi/ebiten/v2/internal/gamepaddb.parseLine
        1.17MB  0.48% 97.46%     1.74MB  0.72%  compress/flate.(*compressor).init
        0.88MB  0.37% 97.83%     2.62MB  1.09%  compress/flate.NewWriter (inline)
        0.50MB  0.21% 98.04%       10MB  4.15%  othello.(*Node).Expand
            0     0% 98.04%     2.62MB  1.09%  compress/gzip.(*Writer).Write
    go tool pprof cpu.out                                                                                 
    File: othello.test
    Type: cpu
    Time: 2025-11-17 19:48:56 EST
    Duration: 1.32s, Total samples = 5.51s (418.89%)
    Entering interactive mode (type "help" for commands, "o" for options)
    (pprof) top
    Showing nodes accounting for 5230ms, 94.92% of 5510ms total
    Dropped 38 nodes (cum <= 27.55ms)
    Showing top 10 nodes out of 76
        flat  flat%   sum%        cum   cum%
        1570ms 28.49% 28.49%     1610ms 29.22%  othello.shift (inline)
        930ms 16.88% 45.37%      930ms 16.88%  runtime.pthread_cond_signal
        650ms 11.80% 57.17%     1900ms 34.48%  othello.generateMoves
        570ms 10.34% 67.51%      570ms 10.34%  runtime.madvise
        460ms  8.35% 75.86%      460ms  8.35%  runtime.usleep
        360ms  6.53% 82.40%      360ms  6.53%  runtime.pthread_cond_wait
        230ms  4.17% 86.57%      310ms  5.63%  othello.ArrayOfMoves
        190ms  3.45% 90.02%      620ms 11.25%  othello.resolveMove
        140ms  2.54% 92.56%      140ms  2.54%  runtime.asyncPreempt
        130ms  2.36% 94.92%      130ms  2.36%  runtime.pthread_kill

After:

    BenchmarkSingleRunParallelizationMCTS-8   	     160	   7394913 ns/op	  130839 B/op	    2324 allocs/op
    PASS
    ok  	othello	1.817s
    File: othello.test
    Type: alloc_space
    Time: 2025-11-18 23:22:53 EST
    Entering interactive mode (type "help" for commands, "o" for options)
    (pprof) unit MB
    did you mean: unit=MB
    (pprof) unit=MB
    (pprof) top    
    Showing nodes accounting for 24.82MB, 86.12% of 28.82MB total
    Showing top 10 nodes out of 76
        flat  flat%   sum%        cum   cum%
        8.54MB 29.64% 29.64%     8.54MB 29.64%  math/rand.newSource
        8.50MB 29.49% 59.13%    10.50MB 36.43%  othello.NewNode
        1.50MB  5.21% 64.34%     1.50MB  5.21%  github.com/hajimehoshi/ebiten/v2/internal/gamepaddb.parseLine
        1.50MB  5.20% 69.55%     1.50MB  5.20%  othello.ArrayOfPositionalMoves (inline)
        1.16MB  4.01% 73.56%     1.16MB  4.01%  runtime/pprof.StartCPUProfile
        1MB  3.48% 77.05%        1MB  3.48%  github.com/hajimehoshi/ebiten/v2/internal/graphics.shaderSuffix
        1MB  3.48% 80.52%        1MB  3.48%  runtime.allocm
        0.55MB  1.90% 82.42%     0.55MB  1.90%  github.com/hajimehoshi/ebiten/v2.imageToBytesSlow
        0.55MB  1.90% 84.33%     0.55MB  1.90%  image.NewNRGBA
        0.52MB  1.79% 86.12%     0.52MB  1.79%  regexp.(*bitState).reset
    Showing nodes accounting for 5030ms, 96.92% of 5190ms total
    Dropped 31 nodes (cum <= 25.95ms)
    Showing top 10 nodes out of 65
        flat  flat%   sum%        cum   cum%
        1630ms 31.41% 31.41%     1630ms 31.41%  othello.shift (inline)
        1480ms 28.52% 59.92%     1480ms 28.52%  runtime.pthread_cond_signal
        860ms 16.57% 76.49%     2090ms 40.27%  othello.generateMoves
        330ms  6.36% 82.85%      330ms  6.36%  runtime.madvise
        300ms  5.78% 88.63%      710ms 13.68%  othello.resolveMove
        200ms  3.85% 92.49%      200ms  3.85%  runtime.pthread_cond_wait
        90ms  1.73% 94.22%       90ms  1.73%  runtime.memclrNoHeapPointers
        70ms  1.35% 95.57%       70ms  1.35%  runtime.usleep
        40ms  0.77% 96.34%     2860ms 55.11%  othello.SimulateRollout
        30ms  0.58% 96.92%       30ms  0.58%  internal/runtime/atomic.(*UnsafePointer).Load

Additionally there were massive improvements in terms of allocations and bit operations:
Before:

        BenchmarkInnacurateMonteCarloTreeSearch-8           	      30	  34751862 ns/op	 1706158 B/op	  110736 allocs/op
        BenchmarkOriginalMonteCarloTreeSearch-8             	      33	  34726694 ns/op	 1655930 B/op	  110337 allocs/op
        BenchmarkSingleRunParallelizationMCTS-8             	     136	   8713589 ns/op	 1757373 B/op	  115098 allocs/op
        BenchmarkRollout-8                                  	   16195	     73988 ns/op	    3311 B/op	     232 allocs/op
        BenchmarkRolloutParallel-8                          	   77290	     15196 ns/op	    3314 B/op	     232 allocs/op
        BenchmarkInnacurateMonteCarloTreeSearchParallel-8   	     158	   7617586 ns/op	 1744979 B/op	  113841 allocs/op
        BenchmarkInitialNodeCreationParallel-8              	10059080	       101.3 ns/op	     128 B/op	       3 allocs/op
        BenchmarkVersus-8                                   	       2	 746328500 ns/op	49786616 B/op	 3259000 allocs/op

After (Notice the 0s):

        BenchmarkInnacurateMonteCarloTreeSearch-8           	      34	  29879812 ns/op	  167925 B/op	    5234 allocs/op
        BenchmarkOriginalMonteCarloTreeSearch-8             	      38	  30451938 ns/op	   83539 B/op	    2506 allocs/op
        BenchmarkSingleRunParallelizationMCTS-8             	     163	   7167516 ns/op	  130636 B/op	    2326 allocs/op
        BenchmarkRollout-8                                  	   18514	     64359 ns/op	       0 B/op	       0 allocs/op
        BenchmarkRolloutParallel-8                          	   92695	     12469 ns/op	       0 B/op	       0 allocs/op
        BenchmarkInnacurateMonteCarloTreeSearchParallel-8   	     200	   6247379 ns/op	  164263 B/op	    4874 allocs/op
        BenchmarkInitialNodeCreationParallel-8              	11385991	        97.09 ns/op	     128 B/op	       3 allocs/op
        BenchmarkVersus-8                                   	       2	 599196854 ns/op	 6858336 B/op	  162794 allocs/op

#### Chainging to just uint8 to represent moves

Before:

    BenchmarkRollout-8                                  	   18514	     64359 ns/op	       0 B/op	       0 allocs/op
    BenchmarkRolloutParallel-8                          	   92695	     12469 ns/op	       0 B/op	       0 allocs/op
    BenchmarkInnacurateMonteCarloTreeSearchParallel-8   	     200	   6247379 ns/op	  164263 B/op	    4874 allocs/op
    BenchmarkInitialNodeCreationParallel-8              	11385991	        97.09 ns/op	     128 B/op	       3 allocs/op
    BenchmarkVersus-8                                   	       2	 599196854 ns/op	 6858336 B/op	  162794 allocs/op

After:

    BenchmarkInnacurateMonteCarloTreeSearch-8           	      37	  30288100 ns/op	  174411 B/op	    2866 allocs/op
    BenchmarkOriginalMonteCarloTreeSearch-8             	      38	  30469754 ns/op	   88376 B/op	    1397 allocs/op
    BenchmarkMonteCarloTreeSearchPUCT-8                 	      38	  31097039 ns/op	  302717 B/op	    3881 allocs/op
    BenchmarkSingleRunParallelizationMCTS-8             	     156	   7634688 ns/op	  142079 B/op	    1518 allocs/op
    BenchmarkRollout-8                                  	   18026	     65334 ns/op	       0 B/op	       0 allocs/op
    BenchmarkRolloutParallel-8                          	   88774	     13735 ns/op	       0 B/op	       0 allocs/op
    BenchmarkInnacurateMonteCarloTreeSearchParallel-8   	     134	   8810127 ns/op	  179820 B/op	    2909 allocs/op
    BenchmarkVersus-8                                   	       2	 728645666 ns/op	 6560808 B/op	   78171 allocs/op

There is still work to be done, optimizing the allocation of the size of the slices used. Paritcularly in the Node creation as we know how many children nodes there will be as a maximum.

### Implementation of MCTS with PUCT
I had to implement a whole new MCTS fit to use PUCT, but it seems it is better than the common UCT.
It is not currently optimized so it can still be improved in terms of speed, but the results against montecarlo tree search with UCT are promising.

Benchmark: 

    BenchmarkMonteCarloTreeSearchPUCT-8                 	      37	  31263037 ns/op	  296617 B/op	    4993 allocs/op
    BenchmarkOriginalMonteCarloTreeSearch-8             	      38	  30599989 ns/op	   83594 B/op	    2508 allocs/op

As you can see above the writes are triple those of the original MCTS and the allocations are almos double. I believe this is to my way of representing the moves as [2]uint8, maybe a switch to []uint8 would improve this.

Results as the white opponent of the MCTS UCT (both having 500 simulations per turn):

    Opponent Wins: 57
    Draws: 5
    Total Games ran: 100
    Total run time for all the games: 1m40.238196083s

    Opponent Wins: 67
    Draws: 7
    Total Games ran: 100
    Total run time for all the games: 1m50.454568167s

I have ran it as white multiple times and it appears that we get variable results but all above 57 and below 69 wins. This is quite a lot, but I also wonder if the fact that they both share the rng affects the code. I might try it later with 2 independent rng sources (Theoretically it should not affect, but i could do it just to be sure).

Results as the black opponent of the MCTS UCT (both having 500 simulations per turn):

    Opponent Wins: 64
    Draws: 6
    Total Games ran: 100
    Total run time for all the games: 1m40.353152166s

### Align fields to avoid cache line problems

Before: 

    BenchmarkInnacurateMonteCarloTreeSearch-8           	      37	  30288100 ns/op	  174411 B/op	    2866 allocs/op
    BenchmarkOriginalMonteCarloTreeSearch-8             	      38	  30469754 ns/op	   88376 B/op	    1397 allocs/op
    BenchmarkMonteCarloTreeSearchPUCT-8                 	      38	  31097039 ns/op	  302717 B/op	    3881 allocs/op
    BenchmarkSingleRunParallelizationMCTS-8             	     156	   7634688 ns/op	  142079 B/op	    1518 allocs/op
    BenchmarkRollout-8                                  	   18026	     65334 ns/op	       0 B/op	       0 allocs/op
    BenchmarkRolloutParallel-8                          	   88774	     13735 ns/op	       0 B/op	       0 allocs/op
    BenchmarkInnacurateMonteCarloTreeSearchParallel-8   	     134	   8810127 ns/op	  179820 B/op	    2909 allocs/op
    BenchmarkVersus-8                                   	       2	 728645666 ns/op	 6560808 B/op	   78171 allocs/op

After:

    BenchmarkInnacurateMonteCarloTreeSearch-8           	      36	  30143149 ns/op	  174613 B/op	    2867 allocs/op
    BenchmarkOriginalMonteCarloTreeSearch-8             	      38	  30215717 ns/op	   88595 B/op	    1395 allocs/op
    BenchmarkMonteCarloTreeSearchPUCT-8                 	      38	  30551758 ns/op	  301920 B/op	    3874 allocs/op
    BenchmarkSingleRunParallelizationMCTS-8             	     169	   7022824 ns/op	  141911 B/op	    1518 allocs/op
    BenchmarkRollout-8                                  	   18679	     64245 ns/op	       0 B/op	       0 allocs/op
    BenchmarkRolloutParallel-8                          	   92790	     11574 ns/op	       0 B/op	       0 allocs/op
    BenchmarkInnacurateMonteCarloTreeSearchParallel-8   	     182	   5911194 ns/op	  179843 B/op	    2909 allocs/op
    BenchmarkVersus-8                                   	       2	 593293375 ns/op	 6632804 B/op	   79010 allocs/op
    PASS
    ok  	othello	9.883s