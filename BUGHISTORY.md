# Original Notes done during development related to bugs
## Incorrect implementation was better than original MCTS**

**TLDR: The origianl MCTS was implemented incorrectly so my agressive MCTS was performing better**
Notes before discovering the bugs:
    It seems like whoever is black has an advantage on the game, it seems like the AI loses more easily if it is white rather than black. 

        Possible fixes:

        - Add more compute in case that AI is white
        - Modify the C parameter?

    Actual root cause: The backpropagation algorithm is not taking into account if the machine is black or white, defaults to black. Therefore it is easier because technically the AI currently is choosing the moves that are more likely to make black win, even when the machine is white. The machine is acting against its own interests.

    The problem has now been fixed. But it seems like the AI plays a kind of different strategy (I know it is just simulation and probability) when it plays white. It plays a long game, in my opinion.


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

After running both my version and the Montecarlo tree search original version, now with a correct implmentation. The MCTS original version is better
My algo as Black:
Opponent Wins: 30
Draws: 14
Total Games ran: 100
Total run time for all the games: 18m7.259308792s%  
My algo as white:
Opponent Wins: 39
Draws: 7
Total Games ran: 100
Total run time for all the games: 18m12.5226s%  