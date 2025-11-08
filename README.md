# Othello-Engine
An Othello/Reversi AI and Engine in Go.

Currently features a Montecarlo tree search artificial intelligence, fully coded in golang from scratch.

Uses a Bitboard implementation to optimize for simulation speed and therefore number of simulations.
Uses Montecarlo Tree search to select the best move.

Next features:

    - AI will be able to play as not only black but also white (DONE)
    - Will have a nicer UI (not just terminal) (DONE, but can be improved)
    - Will have difficulty selection
    - Will have a harder AI
    - WIll be available on itchio or something (Available on Itchio https://nanuklovesfish3.itch.io/simple-othello)

Notes:
    It seems like whoever is black has an advantage on the game, it seems like the AI loses more easily if it is white rather than black. 

        Possible fixes:

        - Add more compute in case that AI is white
        - Modify the C parameter?

    Actual root cause: The backpropagation algorithm is not taking into account if the machine is black or white, defaults to black. Therefore it is easier because technically the AI currently is choosing the moves that are more likely to make black win, even when the machine is white. The machine is acting against its own interests.

    The problem has now been fixed. But it seems like the AI plays a kind of different strategy (I know it is just simulation and probability) when it plays white. It plays a long game, in my opinion.

Current Ideas:
- Implement a way to test 2 AIs against each other, so that they can be benchmarked
    - Then once it is confirmed that the erroneous implementation is better, try to think why is it better
- Implement leaf or root parallelization (Hopefully this will be straightforward)
- Create a Neural Network that analyzes the current leaf to see how it will play out (maybe through self play maybe through a dataset)
- Replace UCT with other methods seen in previous research paper
    - Ask for questions to the respective creators of the evaluation functions (?)
- Get EDAX or Egaroucid running to test the game
- Maybe try to introduce the evaluation pattern used by Logistello (in some sort of way)
- At endgame, run another Algorithm instead of MCTS maybe Minimax (The depth should be small enough to get the actual best move)

- Implement NegaScout Algorithm (?)