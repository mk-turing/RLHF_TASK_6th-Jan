package _12241

import (
	"math/rand"
	"testing"
)

type TestVariant string

func GetVariant() TestVariant {
	if rand.Intn(2) == 0 {
		return TestVariant("A")
	}
	return TestVariant("B")
}

func BenchmarkVariantA(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetVariant() // This will call the variant selection logic
	}
}

func BenchmarkVariantB(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetVariant() // This will call the variant selection logic
	}
}
