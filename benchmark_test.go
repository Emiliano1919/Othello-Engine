package main

import (
	"math/rand"
	"testing"
	"time"
)

func BenchmarkInnacurateMonteCarloTreeSearch(b *testing.B) {
	node := InitialRootNode()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for b.Loop() {
		InnacurateMonteCarloTreeSearch(node, 500, OPTIMIZE_FOR_BLACK, rng)
	}
}

func BenchmarkOriginalMonteCarloTreeSearch(b *testing.B) {
	node := InitialRootNode()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for b.Loop() {
		OriginalMonteCarloTreeSearch(node, 500, rng)
	}
}
func BenchmarkSingleRunParallelizationMCTS(b *testing.B) {
	node := InitialRootNode()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for b.Loop() {
		SingleRunParallelizationMCTS(node, 50, rng)
	}
}

func BenchmarkRollout(b *testing.B) {
	node := InitialRootNode()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for b.Loop() {
		SimulateRollout(node.GameState, rng)
	}
}

func BenchmarkRolloutParallel(b *testing.B) {
	nodeP := InitialRootNode()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			SimulateRollout(nodeP.GameState, rng)
		}
	})
}

func BenchmarkInnacurateMonteCarloTreeSearchParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			nodePM := InitialRootNode()
			rng := rand.New(rand.NewSource(time.Now().UnixNano()))
			InnacurateMonteCarloTreeSearch(nodePM, 500, OPTIMIZE_FOR_BLACK, rng)
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
