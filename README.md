# Bioinformatics algorithms

This repository holds Go implementations of the bioinformatics algorithms covered during my undergraduate degree, see the [course page](https://www.cl.cam.ac.uk/teaching/2526/Bioinfo/) for more details about the course.
Most of the algorithms here were coded as an optional exercise, rather than for completion of the course.
All completed algorithms are in their own package, the `main.go` file executes an algorithm according to command line input.

Here are some notes about the implementation and details of each algorithm.

## Needleman-Wunsch: global alignment

The Needleman-Wunsch algorithm finds a maximal alignment of two strings.
The algorithm is a slight variation on the problem of finding the longest path in a directed acyclic graph (DAG), since the shape of the DAG is known to be a Manhattan-style grid.
This also lets the algorithm know that there are only three possible moves to consider, rather than having to look at all predecessors.
