# Gosta
detecting Goroutine bugs with Static analysis methodologies

## Current Features
+ Detecting goroutine-deadlock by abstract execution on GFSM(Goroutine Finite State Machine)


## TODO LIST(Small)
+ Kill the redundant SFSM(Small Finite State Machine for a Basic Block) when transforming SFSMs to GFSMs.
+ Add heuristic algorithms, which can lead the tool to the paths that are most likely to have goroutine deadlocks. 


## TODO LIST(Big)
+ Support goroutine-leak
+ After a BUG detected, we still need to check the path sensitivity using SMT-Solver.
