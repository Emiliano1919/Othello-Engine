# Othello-Engine
An Othello/Reversi AI and Engine in Go.

Currently features a Montecarlo tree search artificial intelligence, fully coded in golang from scratch.

Uses a Bitboard implementation to optimize for simulation speed and therefore number of simulations.
Uses Montecarlo Tree search to select the best move.

Next features:

    - AI will be able to play as not only black but also white
    - Will have a nicer UI (not just terminal)
    - Will have difficulty selection
    - Will have a harder AI

Notes:
    It seems like whoever is black has an advantage on the game, it seems like the AI loses more easily if it is white rather than black.
    
        Possible fixes:

        - Add more compute in case that AI is white
        - Modify the C parameter?
