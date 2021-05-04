package main

import (
	"testing"
)

func BenchmarkGoogleSequential(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GoogleConcurrent("golang")
	}
}

func BenchmarkGoogleConcurrent(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GoogleConcurrent("golang")
	}
}

func BenchmarkGoogleReplicas(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GoogleReplicas("golang")
	}
}
