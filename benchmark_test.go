package main

import (
	"testing"
)

func BenchmarkInnacurateMonteCarloTreeSearch(b *testing.B) {
	node := InitialRootNode()
	for b.Loop() {
		InnacurateMonteCarloTreeSearch(node, 500, OPTIMIZE_FOR_BLACK)
	}
}

func BenchmarkOriginalMonteCarloTreeSearch(b *testing.B) {
	node := InitialRootNode()
	for b.Loop() {
		OriginalMonteCarloTreeSearch(node, 500)
	}
}
func BenchmarkSingleRunParallelizationMCTS(b *testing.B) {
	node := InitialRootNode()
	for b.Loop() {
		SingleRunParallelizationMCTS(node, 200)
	}
}

func BenchmarkRollout(b *testing.B) {
	node := InitialRootNode()
	for b.Loop() {
		SimulateRollout(node.GameState, 0)
	}
}

func BenchmarkRolloutParallel(b *testing.B) {
	nodeP := InitialRootNode()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			SimulateRollout(nodeP.GameState, 0)
		}
	})
}

func BenchmarkInnacurateMonteCarloTreeSearchParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			nodePM := InitialRootNode()
			InnacurateMonteCarloTreeSearch(nodePM, 500, OPTIMIZE_FOR_BLACK)
		}
	})
}

func BenchmarkInitialNodeCreationParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			InitialRootNode()
		}
	})
}

func BenchmarkVersus(b *testing.B) {
	for b.Loop() {
		Versus()
	}
}
