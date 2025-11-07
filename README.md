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
