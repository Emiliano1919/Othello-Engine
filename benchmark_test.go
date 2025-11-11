package main

import (
	"testing"
)

func BenchmarkMonteCarloTreeSearc(b *testing.B) {
	node := InitialRootNode()
	for b.Loop() {
		MonteCarloTreeSearch(node, 500, OPTIMIZE_FOR_BLACK)
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

func BenchmarkMonteCarloTreeSearchParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			nodePM := InitialRootNode()
			MonteCarloTreeSearch(nodePM, 500, OPTIMIZE_FOR_BLACK)
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
